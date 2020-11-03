package downloader

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/go-redsync/redsync/v4"
)

type performer struct {
	Name string `json:"performer_name"`
}

type Downloads struct {
	rs         *redsync.Redsync
	mutexes    map[string]*redsync.Mutex
	guardMutex *sync.Mutex
}

func New(rs *redsync.Redsync) *Downloads {
	return &Downloads{
		rs:         rs,
		mutexes:    make(map[string]*redsync.Mutex),
		guardMutex: &sync.Mutex{},
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

	cmd := exec.Command(
		"streamlink",
		"https://chaturbate.com/"+name+"/",
		"best",
		"-o", "/app/downloads/"+name+".mp4", // TODO: add timestamp
		"--loglevel", "debug",
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
}

func copyOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
