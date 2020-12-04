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
	Followers       int64     `json:"-" db:"followers"`
}

// NewRecord creates Record
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

// Finish sets finish_at
func (r *Record) Finish(db *sqlx.DB) error {
	if _, err := db.Exec(`UPDATE records SET finish_at = NOW() WHERE id = $1`, r.ID); err != nil {
		return err
	}
	return nil
}

func TotalFeaturedRecords(db *sqlx.DB) (int64, error) {
	var total int64

	if err := db.Get(
		&total,
		`SELECT
			COUNT(*)
		FROM records t1
		JOIN (
			SELECT
				r.broadcaster_id AS broadcaster_id,
				MAX(r.finish_at) AS max_finish_at
			FROM records r
			WHERE r.finish_at IS NOT NULL
			GROUP BY r.broadcaster_id
		) t2 ON t1.broadcaster_id = t2.broadcaster_id AND t1.finish_at IS NOT NULL AND t1.finish_at = t2.max_finish_at`,
	); err != nil {
		return 0, err
	}

	return total, nil
}

// FeaturedRecords forms list of feaured performers
func FeaturedRecords(db *sqlx.DB, page int64, limit int) ([]*Record, error) {
	r := []*Record{}

	sql := `
		SELECT
			t1.id AS id,
			t1.broadcaster_id,
			t1.created_at,
			t1.finish_at,
			t1.uuid,
			b.name AS broadcaster_name,
			COALESCE(b.followers, 0) AS followers
		FROM records t1
		JOIN (
			SELECT
				r.broadcaster_id AS broadcaster_id,
				MAX(r.finish_at) AS max_finish_at
			FROM records r
			WHERE r.finish_at IS NOT NULL
			GROUP BY r.broadcaster_id
		) t2 ON t1.broadcaster_id = t2.broadcaster_id AND t1.finish_at = t2.max_finish_at
		INNER JOIN broadcasters b ON b.id = t1.broadcaster_id
		ORDER BY date_trunc('day', t1.finish_at) DESC, COALESCE(b.followers, 0) DESC, t1.id DESC LIMIT $1 OFFSET $2`

	err := db.Select(&r, sql, limit, (page-1)*int64(limit))
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
