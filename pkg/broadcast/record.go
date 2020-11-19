package broadcast

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Record is record of broadcast
type Record struct {
	ID              int64     `json:"id,omitempty" db:"id"`
	BroadcasterName string    `json:"name" db:"broadcaster_name"`
	BroadcasterID   int64     `json:"-" db:"broadcaster_id"`
	UUID            string    `json:"uuid" db:"uuid"`
	CreatedAt       time.Time `json:"-" db:"created_at"`
	FinishAt        time.Time `json:"-" db:"finish_at"`
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

func FeaturedRecords(db *sqlx.DB) ([]*Record, error) {
	r := []*Record{}

	err := db.Select(
		&r, `
		SELECT
			x.id,
			x.broadcaster_id,
			x.broadcaster_name,
			x.created_at,
			x.finish_at,
			x.uuid
		FROM (
			SELECT
					r.id AS id,
					b.id AS broadcaster_id,
					b.name AS broadcaster_name,
					b.followers,
					r.created_at,
					r.finish_at,
					r.uuid,
					row_number() OVER (PARTITION BY b.id ORDER BY r.finish_at DESC) AS n
				FROM records r INNER JOIN broadcasters b ON b.id = r.broadcaster_id
				WHERE r.finish_at IS NOT NULL) x
		WHERE x.n = 1 ORDER BY date_trunc('day', x.finish_at) DESC, x.followers DESC NULLS LAST`,
	)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// FindByUUID finds record by UUID
func FindByUUID(db *sqlx.DB, uuid string) (*Record, error) {
	r := &Record{}

	err := db.Get(r,
		`SELECT
			r.id AS id,
			b.id AS broadcaster_id,
			b.name AS broadcaster_name,
			r.created_at,
			r.finish_at,
			r.uuid
		FROM records r INNER JOIN broadcasters b ON b.id = r.broadcaster_id
		WHERE r.uuid = $1 AND r.finish_at IS NOT NULL LIMIT 1`,
		uuid,
	)
	if err != nil {
		return nil, err
	}
	return r, nil
}
