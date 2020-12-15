package downloader

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
)

type Handler struct {
	db              *sqlx.DB
	rs              *redsync.Redsync
	nc              *nats.Conn
	activeCmds      map[string]*Download
	guardActiveCmds *sync.Mutex
	sub             *nats.Subscription

	performers chan string

	// Handle shutdown
	quit chan struct{}
	done chan struct{}
}

func New(rs *redsync.Redsync, db *sqlx.DB, nc *nats.Conn) *Handler {
	return &Handler{
		db:              db,
		rs:              rs,
		nc:              nc,
		activeCmds:      make(map[string]*Download),
		guardActiveCmds: &sync.Mutex{},
		quit:            make(chan struct{}),
		done:            make(chan struct{}),
		performers:      make(chan string),
	}
}

func (h *Handler) Run() {
	log.Println("DEBUG: run main loop")
	// Subscribe
	subscribtion, err := h.nc.QueueSubscribe("downloading", "download", func(m *nats.Msg) {
		log.Printf("DEBUG: got message: %s", string(m.Data[:]))

		h.performers <- string(m.Data[:])
	})
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	h.sub = subscribtion

	defer log.Println("DEBUG: stopped main loop")

	for {
		select {
		case record := <-h.performers:
			go h.Start(record)
		case _ = <-h.quit:
			h.close()
			return
		}
	}
}

func (h *Handler) Start(performer string) {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Minute)
	defer cancel()

	d := &Download{
		db:        h.db,
		nc:        h.nc,
		performer: performer,
		mutex:     h.rs.NewMutex("downloads:locks:"+performer, redsync.WithExpiry(8*time.Hour)),
		done:      make(chan struct{}),
	}
	h.guardActiveCmds.Lock()
	h.activeCmds[performer] = d
	h.guardActiveCmds.Unlock()

	done, err := d.Run()
	if err != nil {
		log.Printf("ERROR: %v", err)
		return
	}

	select {
	case <-done:
		log.Println("INFO: Download complete...")
	case <-ctx.Done():
		log.Println("INFO: Timed out...")
		d.Close()
	}

	h.guardActiveCmds.Lock()
	delete(h.activeCmds, performer)
	h.guardActiveCmds.Unlock()
}

// Close gracefully stops all run tasks
func (h *Handler) Close() chan struct{} {
	close(h.quit)

	return h.done
}

func (h *Handler) close() {
	close(h.performers)

	if err := h.sub.Unsubscribe(); err != nil {
		log.Printf("ERROR: %v", err)
	}

	// Send SIGINT to all commands
	log.Println("DEBUG: wait for finish all active downloads...")
	for _, cmd := range h.activeCmds {
		cmd.Close()
	}

	close(h.done)
}
