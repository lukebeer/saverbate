package downloader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"saverbate/pkg/broadcast"
	"sync"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
)

const (
	userAgent = `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36`
)

// Downloads handles downloading
type Downloads struct {
	rs *redsync.Redsync
	db *sqlx.DB
	nc *nats.Conn

	mutexes    map[string]*redsync.Mutex
	guardMutex *sync.Mutex

	activeCmds      map[string]*exec.Cmd
	guardActiveCmds *sync.Mutex

	performers chan string
}

// New creates new instance of Downloads
func New(rs *redsync.Redsync, db *sqlx.DB, nc *nats.Conn) *Downloads {
	return &Downloads{
		rs:              rs,
		mutexes:         make(map[string]*redsync.Mutex),
		guardMutex:      &sync.Mutex{},
		db:              db,
		nc:              nc,
		activeCmds:      make(map[string]*exec.Cmd),
		guardActiveCmds: &sync.Mutex{},
		performers:      make(chan string, 1),
	}
}

// Start runs new download by name
func (d *Downloads) Start(name string) {
	for {
		select {
		case performer := <-d.performers:
			go d.start(performer)
		}
	}
}

func (d *Downloads) start(name string) {
	d.guardMutex.Lock()
	if _, ok := d.mutexes[name]; !ok {
		d.mutexes[name] = d.rs.NewMutex(name, redsync.WithExpiry(36*time.Hour))
	}
	d.guardMutex.Unlock()

	if err := d.mutexes[name].Lock(); err != nil {
		log.Println("Download of " + name + " already run")
		return
	}

	defer func() {
		if ok, err := d.mutexes[name].Unlock(); !ok || err != nil {
			log.Println("Could not release crawler lock")
			return
		}
	}()

	r, err := broadcast.NewRecord(d.db, name)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return
	}

	cmd := exec.Command(
		"youtube-dl",
		"--no-color",
		"--no-call-home",
		"--no-progress",
		"--user-agent", userAgent,
		"-f", "best[height<=560]",
		"--output", "/app/downloads/"+name+"/"+r.UUID+".%(ext)s",
		"https://chaturbate.com/"+name+"/",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	err = cmd.Start()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	d.guardActiveCmds.Lock()
	d.activeCmds[name] = cmd
	d.guardActiveCmds.Unlock()

	go copyOutput(stdout)
	go copyOutput(stderr)
	cmd.Wait()

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
	}
}

func copyOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
