package downloader

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"saverbate/pkg/broadcast"
	"saverbate/pkg/utils"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
)

const (
	userAgent = `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36`
)

type Download struct {
	db *sqlx.DB
	nc *nats.Conn

	performer string
	mutex     *redsync.Mutex
	cmd       *exec.Cmd
	done      chan struct{}
}

func (d *Download) Run() (chan struct{}, error) {
	if err := d.mutex.Lock(); err != nil {
		return nil, fmt.Errorf("ERROR: Download of "+d.performer+" already run: %v", err)
	}

	defer func() {
		if ok, err := d.mutex.Unlock(); !ok || err != nil {
			log.Printf("ERROR: Could not release lock: %v", err)
		}
	}()

	r, err := broadcast.NewRecord(d.db, d.performer)
	if err != nil {
		return nil, fmt.Errorf("ERROR: failed to save record %v", err)
	}

	done := make(chan struct{})
	cmd := exec.Command(
		"youtube-dl",
		"--socket-timeout", "10",
		"--retries", "3",
		"--quiet",
		"--no-warnings",
		"--no-color",
		"--no-call-home",
		"--no-progress",
		"--user-agent", userAgent,
		"-f", "best[height<720]",
		"--output", "/app/downloads/"+d.performer+"/"+r.UUID+".%(ext)s",
		"https://chaturbate.com/"+d.performer+"/",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("ERROR: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("ERROR: %v", err)
	}
	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("ERROR: %v", err)
	}
	d.cmd = cmd

	go utils.CopyOutput(stdout)
	go utils.CopyOutput(stderr)

	go func() {
		if err := cmd.Wait(); err != nil {
			log.Printf("ERROR: error execute command: %v, name: %s", err, d.performer)
			return
		}

		d.cmd = nil

		d.finish(r)

		close(d.done)
		close(done)
	}()

	return done, nil
}

func (d *Download) Close() {
	if d.cmd == nil {
		return
	}

	for i := 0; i < 10; i++ {
		select {
		case <-d.done:
			return
		case <-time.After(2 * time.Second):
			if i > 0 {
				log.Printf("INFO: send SIGINT after %d try", i)
			}

			if err := d.cmd.Process.Signal(os.Interrupt); err != nil {
				log.Printf("ERROR: sending SIGINT failed: %v, for %s", err, d.performer)
			}
		}
	}
}

func (d *Download) finish(r *broadcast.Record) {
	if err := r.Finish(d.db); err != nil {
		log.Printf("ERROR: %v", err)
		return
	}

	message, err := json.Marshal(r)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return
	}
	if err := d.nc.Publish("download_complete", message); err != nil {
		log.Printf("ERROR: %v", err)
		return
	}

	if err := d.nc.Flush(); err != nil {
		log.Printf("ERROR: %v", err)
		return
	}
}
