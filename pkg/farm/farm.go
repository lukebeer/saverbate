package farm

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

// State is a state of the farm
type State string

const (
	// StateInactive is inactive state of the farm
	StateInactive State = "inactive"
	// StateActive is active state
	StateActive State = "active"
	// StateFailed is setted when farm in failing state
	StateFailed State = "failed"
)

// Instance is instance of the farm
type Instance struct {
	ID              int64  `db:"id"`
	Name            string `db:"name"`
	Password        string `db:"password"`
	MailerState     State  `db:"mailer_state"`
	SubscriberState State  `db:"subscriber_state"`
	PortalUsername  string `db:"portal_username"`
	PortalPassword  string `db:"portal_password"`
}

// FindFree finds free (inactive) farm for occupation
func FindFree(db *sqlx.DB, serviceName string) (*Instance, error) {
	farm := &Instance{}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	row := tx.QueryRow(`
		SELECT
			id,
			name,
			mailer_state,
			subscriber_state,
			password,
			portal_username,
			portal_password
		FROM farms WHERE ` + serviceName + `_state = 'inactive' LIMIT 1 FOR UPDATE
	`)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = row.Scan(
		&(farm.ID),
		&(farm.Name),
		&(farm.MailerState),
		&(farm.SubscriberState),
		&(farm.Password),
		&(farm.PortalUsername),
		&(farm.PortalPassword),
	)
	if err == sql.ErrNoRows {
		tx.Rollback()
		return nil, errors.New("No free farms")
	} else {
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	_, err = tx.Exec(`UPDATE farms SET `+serviceName+`_state = 'active' WHERE id = $1`, farm.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	if serviceName == "mailer" {
		farm.MailerState = StateActive
	} else {
		farm.SubscriberState = StateActive
	}

	return farm, nil
}

// Release move farm to inactive
func Release(db *sqlx.DB, serviceName, name string) error {
	_, err := db.Exec(`UPDATE farms SET `+serviceName+`_state = 'inactive' WHERE name = $1`, name)
	if err != nil {
		return err
	}

	return nil
}
