package downloader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
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

	// Handle shutdown
	quit chan struct{}
	done chan struct{}
}

// New creates new instance of Downloads
func New(rs *redsync.Redsync, db *sqlx.DB, nc *nats.Conn) *Downloads {
	return &Downloads{
		rs:              rs,
		db:              db,
		nc:              nc,
		mutexes:         make(map[string]*redsync.Mutex),
		guardMutex:      &sync.Mutex{},
		activeCmds:      make(map[string]*exec.Cmd),
		guardActiveCmds: &sync.Mutex{},
		performers:      make(chan string, 1),
		quit:            make(chan struct{}),
		done:            make(chan struct{}),
	}
}

// Run runs main loop of downloads
func (d *Downloads) Run() {
	for {
		select {
		case performer := <-d.performers:
			go d.start(performer)
		case _ = <-d.quit:
			d.close()
			return
		}
	}
}

// Start runs new download by name
func (d *Downloads) Start(name string) {
	d.performers <- name
}

func (d *Downloads) start(name string) {
	d.guardMutex.Lock()
	if _, ok := d.mutexes[name]; !ok {
		d.mutexes[name] = d.rs.NewMutex("downloads:locks:"+name, redsync.WithExpiry(36*time.Hour))
	}
	d.guardMutex.Unlock()

	if err := d.mutexes[name].Lock(); err != nil {
		log.Printf("Download of "+name+" already run: %v", err)
		return
	}

	r, err := broadcast.NewRecord(d.db, name)
	if err != nil {
		log.Printf("ERROR: failed to save record %v", err)
		return
	}

	cmd := exec.Command(
		"youtube-dl",
		"--quiet",
		"--no-warnings",
		"--no-color",
		"--no-call-home",
		"--no-progress",
		"--user-agent", userAgent,
		"-f", "bestvideo[filesize<3G][height<=?720]+bestaudio/best",
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
	if err := cmd.Wait(); err != nil {
		log.Printf("ERROR: error execute command: %v, name: %s", err, name)
		return
	}
	d.finishDownload(r)
}

func copyOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func (d *Downloads) finishDownload(r *broadcast.Record) {
	log.Printf("DEBUG: finish download for %s", r.BroadcasterName)

	if ok, err := d.mutexes[r.BroadcasterName].Unlock(); !ok || err != nil {
		log.Printf("Could not release crawler lock: %v", err)
		return
	}
	d.guardMutex.Lock()
	delete(d.mutexes, r.BroadcasterName)
	d.guardMutex.Unlock()

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

	d.guardActiveCmds.Lock()
	delete(d.activeCmds, r.BroadcasterName)
	d.guardActiveCmds.Unlock()

	log.Printf("DEBUG: downlad finished for %s", r.BroadcasterName)
}

// Close closes all current downloads
func (d *Downloads) Close() chan struct{} {
	d.quit <- struct{}{}

	return d.done
}

func (d *Downloads) close() {
	close(d.performers)
	close(d.quit)

	// Send SIGINT to all commands
	for name, cmd := range d.activeCmds {
		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			log.Printf("ERROR: sending SIGINT failed: %v, for %s", err, name)
		}
	}

	// Wait for all
	for l := len(d.activeCmds); l > 0; {
		log.Printf("INFO: wait for finish %d active commands", l)
		time.Sleep(1 * time.Second)
	}
	close(d.done)
}
