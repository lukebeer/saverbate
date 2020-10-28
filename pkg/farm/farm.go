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
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	State    State  `db:"state"`
	Password string `db:"password"`
}

// FindFree finds free (inactive) farm for occupation
func FindFree(db *sqlx.DB) (*Instance, error) {
	farm := &Instance{}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	row := tx.QueryRow(`SELECT id, name, state, password FROM farms WHERE state = 'inactive' LIMIT 1 FOR UPDATE`)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = row.Scan(&(farm.ID), &(farm.Name), &(farm.State), &(farm.Password))
	if err == sql.ErrNoRows {
		tx.Rollback()
		return nil, errors.New("No free farms")
	} else {
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	_, err = tx.Exec(`UPDATE farms SET state = 'active' WHERE id = $1`, farm.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	farm.State = StateActive

	return farm, nil
}

// Release move farm to inactive
func Release(db *sqlx.DB, name string) error {
	_, err := db.Exec(`UPDATE farms SET state = 'inactive' WHERE name = $1`, name)
	if err != nil {
		return err
	}

	return nil
}
