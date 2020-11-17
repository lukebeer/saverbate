package thumbnails

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"saverbate/pkg/broadcast"

	"github.com/go-redsync/redsync/v4"
	"github.com/nats-io/nats.go"
)

// Handler makes thumbnails from videos
type Handler struct {
	rs              *redsync.Redsync
	nc              *nats.Conn
	activeCmds      map[string]*Thumbnail
	guardActiveCmds *sync.Mutex
	sub             *nats.Subscription

	records chan *broadcast.Record

	// Handle shutdown
	quit chan struct{}
	done chan struct{}
}

// New creates new instance of Thumbnails
func New(rs *redsync.Redsync, nc *nats.Conn) *Handler {
	h := &Handler{
		rs:              rs,
		nc:              nc,
		activeCmds:      make(map[string]*Thumbnail),
		guardActiveCmds: &sync.Mutex{},
		records:         make(chan *broadcast.Record),
		quit:            make(chan struct{}),
		done:            make(chan struct{}),
	}

	return h
}

// Run runs main loop
func (t *Handler) Run() {
	// Subscribe
	subscribtion, err := t.nc.QueueSubscribe("download_complete", "download", func(m *nats.Msg) {
		record := &broadcast.Record{}
		if err := json.Unmarshal(m.Data, record); err != nil {
			log.Printf("ERROR: Unmarshal error: %v", err)
			return
		}

		t.records <- record
	})
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	t.sub = subscribtion

	for {
		select {
		case record := <-t.records:
			go t.Start(record)
		case _ = <-t.quit:
			t.close()
			return
		}
	}

}

// Start starts new task for get thumbnails by uuid
func (t *Handler) Start(record *broadcast.Record) {
	thumbnail := &Thumbnail{
		record: record,
		mutex:  t.rs.NewMutex("thumbnails:locks:"+record.UUID, redsync.WithExpiry(36*time.Hour)),
	}
	t.guardActiveCmds.Lock()
	t.activeCmds[record.UUID] = thumbnail
	t.guardActiveCmds.Unlock()

	thumbnail.Make()

	t.guardActiveCmds.Lock()
	delete(t.activeCmds, record.UUID)
	t.guardActiveCmds.Unlock()
}

// Close gracefully stops all run tasks
func (t *Handler) Close() chan struct{} {
	t.quit <- struct{}{}

	return t.done
}

func (t *Handler) close() {
	close(t.records)
	close(t.quit)
	if err := t.sub.Unsubscribe(); err != nil {
		log.Printf("ERROR: %v", err)
	}

	// Send SIGINT to all commands
	//for name, cmd := range t.activeCmds {
	//if err := cmd.Process.Signal(os.Interrupt); err != nil {
	//	log.Printf("ERROR: sending SIGINT failed: %v, for %s", err, name)
	//}
	//}

	// Wait for all
	log.Println("DEBUG: wait for finish all active thumbnails...")
	for len(t.activeCmds) > 0 {
		log.Printf("INFO: wait for finish %d active commands", len(t.activeCmds))
		time.Sleep(1 * time.Second)
	}
	log.Println("DEBUG: closed...")

	close(t.done)
}
