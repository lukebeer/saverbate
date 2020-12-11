package thumbnails

import (
	"encoding/json"
	"log"
	"os/exec"
	"saverbate/pkg/broadcast"
	"saverbate/pkg/utils"

	"github.com/go-redsync/redsync/v4"
	"github.com/nats-io/nats.go"
)

// Thumbnail is thumbnail of video
type Thumbnail struct {
	record *broadcast.Record
	mutex  *redsync.Mutex
}

// Make makes thumbnails from uuid
func (t *Thumbnail) Make(nc *nats.Conn) {
	if err := t.mutex.Lock(); err != nil {
		log.Printf("ERROR: Thumbnail of "+t.record.BroadcasterName+" already run: %v", err)
		return
	}

	defer func() {
		if ok, err := t.mutex.Unlock(); !ok || err != nil {
			log.Printf("ERROR: Could not release lock: %v", err)
			return
		}
	}()

	cmd := exec.Command(
		"/usr/local/bin/thumbnail.sh",
		t.record.UUID,
		t.record.BroadcasterName,
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

	go utils.CopyOutput(stdout)
	go utils.CopyOutput(stderr)
	if err := cmd.Wait(); err != nil {
		log.Printf("ERROR: error execute command: %v, name: %s", err, t.record.BroadcasterName)
		return
	}

	message, err := json.Marshal(t.record)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return
	}
	if err := nc.Publish("thumbnails_complete", message); err != nil {
		log.Printf("ERROR: %v", err)
		return
	}

	if err := nc.Flush(); err != nil {
		log.Printf("ERROR: %v", err)
		return
	}

	log.Printf("DEBUG: making thumbnails finished for %s", t.record.BroadcasterName)
}
