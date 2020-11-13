package broadcast

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Record is record of broadcast
type Record struct {
	ID              int64
	BroadcasterName string
	BroadcasterID   int64
	UUID            string
	CreatedAt       time.Time
	FinishAt        time.Time
}

func NewRecord(db *sqlx.DB, broadcasterName string) (*Record, error) {
	r := &Record{
		UUID:            uuid.New().String(),
		BroadcasterName: broadcasterName,
		CreatedAt:       time.Now(),
	}

	b := &Broadcaster{}

	var insertedID int64

	err := db.Get(b, `SELECT id FROM broadcasters WHERE name = $1`, broadcasterName)
	if err == sql.ErrNoRows {
		if err := db.Get(
			&insertedID,
			`INSERT INTO broadcasters (name, created_at) VALUES ($1, NOW()) RETURNING id`,
			broadcasterName,
		); err != nil {
			return nil, err
		}
		r.BroadcasterID = insertedID

	} else if err != nil {
		return nil, err
	} else {
		r.BroadcasterID = b.ID
	}

	if err := db.Get(
		&insertedID,
		`INSERT INTO records
			(uuid, broadcaster_id, created_at)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		r.UUID,
		r.BroadcasterID,
		r.CreatedAt,
	); err != nil {
		return nil, err
	}

	r.ID = insertedID

	return r, nil
}

func (r *Record) Finish(db *sqlx.DB) error {
	if _, err := db.Exec(`UPDATE records SET finish_at = NOW() WHERE id = $1`, r.ID); err != nil {
		return err
	}
	return nil
}
