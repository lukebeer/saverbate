package downloader

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"saverbate/pkg/broadcast"
	"sync"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/jmoiron/sqlx"
)

const (
	userAgent = `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36`
)

type performer struct {
	Name string `json:"performer_name"`
}

type Downloads struct {
	rs         *redsync.Redsync
	mutexes    map[string]*redsync.Mutex
	guardMutex *sync.Mutex
	db         *sqlx.DB
}

func New(rs *redsync.Redsync, db *sqlx.DB) *Downloads {
	return &Downloads{
		rs:         rs,
		mutexes:    make(map[string]*redsync.Mutex),
		guardMutex: &sync.Mutex{},
		db:         db,
	}
}

func (d *Downloads) Start(name string) {
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

	go copyOutput(stdout)
	go copyOutput(stderr)
	cmd.Wait()

	if err := r.Finish(d.db); err != nil {
		log.Printf("ERROR: %v", err)
		return
	}
}

func copyOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
