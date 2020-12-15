package user

import (
	"context"
	"database/sql"
	"log"
	"time"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"

	"github.com/volatiletech/authboss/v3"
)

// User is type of user
type User struct {
	ID              int64  `json:"id" db:"id"`
	Name            string `json:"name" db:"name"`
	Email           string `json:"email" db:"email"`
	Password        string
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
	Confirmed       bool           `db:"confirmed"`
	ConfirmSelector string         `db:"confirm_selector"`
	ConfirmVerifier string         `db:"confirm_verifier"`
	RecoverSelector sql.NullString `db:"recover_selector"`
	RecoverVerifier sql.NullString `db:"recover_verifier"`
	RecoverExpiry   sql.NullTime   `db:"recover_expiry"`
}

// GetEmail returns user email
func (u *User) GetEmail() string {
	return u.Email
}

// GetConfirmed returns confirmed flag
func (u *User) GetConfirmed() bool {
	return u.Confirmed
}

// GetConfirmSelector returns confirm selector - hash for confirmation
func (u *User) GetConfirmSelector() string {
	return u.ConfirmSelector
}

// GetConfirmVerifier returns string for verify email
func (u *User) GetConfirmVerifier() string {
	return u.ConfirmVerifier
}

// GetPassword return password string
func (u *User) GetPassword() string {
	return u.Password
}

// PutPassword changes password field
func (u *User) PutPassword(password string) {
	u.Password = password
}

// PutEmail changes email field
func (u *User) PutEmail(email string) {
	u.Email = email
}

// PutConfirmed changes confirmed flag
func (u *User) PutConfirmed(confirmed bool) {
	u.Confirmed = confirmed
}

// PutConfirmSelector changes confirmation selector
func (u *User) PutConfirmSelector(confirmSelector string) {
	u.ConfirmSelector = confirmSelector
}

// PutConfirmVerifier changes confirmation verifier
func (u *User) PutConfirmVerifier(confirmVerifier string) {
	u.ConfirmVerifier = confirmVerifier
}

// GetPID returns unique property of user for identify in system
func (u *User) GetPID() string {
	return u.Email
}

// PutPID changes unique identify (email for us)
func (u *User) PutPID(email string) {
	u.Email = email
}

// GetArbitrary returns map of additional user fields such as Name
func (u *User) GetArbitrary() map[string]string {
	return map[string]string{
		"name": u.Name,
	}
}

// GetRecoverSelector returns recover token
func (u *User) GetRecoverSelector() string {
	if u.RecoverSelector.Valid {
		return u.RecoverSelector.String
	}
	return ""
}

// GetRecoverVerifier returns recover token for verification
func (u *User) GetRecoverVerifier() string {
	if u.RecoverVerifier.Valid {
		return u.RecoverVerifier.String
	}

	return ""
}

// GetRecoverExpiry returns recover expiration
func (u *User) GetRecoverExpiry() (expiry time.Time) {
	if u.RecoverExpiry.Valid {
		return u.RecoverExpiry.Time
	}
	return time.Now()
}

// PutRecoverSelector saves recover token to User
func (u *User) PutRecoverSelector(selector string) {
	u.RecoverSelector = sql.NullString{String: selector, Valid: true}
}

// PutRecoverVerifier saves recover verification token to User
func (u *User) PutRecoverVerifier(verifier string) {
	u.RecoverVerifier = sql.NullString{String: verifier, Valid: true}
}

// PutRecoverExpiry saves expiration time of recover token
func (u *User) PutRecoverExpiry(expiry time.Time) {
	u.RecoverExpiry = sql.NullTime{Time: expiry, Valid: true}
}

// PutArbitrary changes additional user fields
func (u *User) PutArbitrary(arbitrary map[string]string) {
	if n, ok := arbitrary["name"]; ok {
		u.Name = n
	}
}

// Storer represent logic of user storage
type Storer struct {
	Db    *sqlx.DB
	redis *goredislib.Client
}

// NewStorer creates storer object with given db connection
func NewStorer(db *sqlx.DB, redis *goredislib.Client) *Storer {
	return &Storer{Db: db, redis: redis}
}

// New returns empty User object
func (s Storer) New(ctx context.Context) authboss.User {
	return &User{}
}

// Save Updates user
func (s Storer) Save(ctx context.Context, user authboss.User) error {
	usr := user.(*User)
	u := &User{}

	findStatement := `SELECT * FROM users WHERE lower(email) = lower($1)`
	err := s.Db.Get(u, findStatement, usr.Email)
	if err == sql.ErrNoRows {
		return authboss.ErrUserNotFound
	}
	if err != nil {
		return err
	}

	updateStatement := `UPDATE users
	  SET name = :name,
	  	confirmed = :confirmed,
	  	password = :password,
	  	confirm_selector = :confirm_selector,
	  	confirm_verifier = :confirm_verifier,
			updated_at = NOW(),
			recover_selector = :recover_selector,
			recover_verifier = :recover_verifier,
			recover_expiry = :recover_expiry
	  WHERE lower(email) = lower(:email)`

	_, err = s.Db.NamedExec(updateStatement,
		map[string]interface{}{
			"name":             usr.Name,
			"email":            usr.Email,
			"confirmed":        usr.Confirmed,
			"password":         usr.Password,
			"confirm_selector": usr.ConfirmSelector,
			"confirm_verifier": usr.ConfirmVerifier,
			"recover_selector": usr.RecoverSelector,
			"recover_verifier": usr.RecoverVerifier,
			"recover_expiry":   usr.RecoverExpiry,
		})

	return err
}

// Load returns User for given identity (email)
func (s Storer) Load(ctx context.Context, key string) (authboss.User, error) {
	u := &User{}

	findStatement := `SELECT * FROM users WHERE lower(email) = lower($1) LIMIT 1`
	err := s.Db.Get(u, findStatement, key)

	if err == sql.ErrNoRows {
		return nil, authboss.ErrUserNotFound
	}
	if err != nil {
		log.Printf("ERROR: %v", err)
		return nil, err
	}

	return u, nil
}

// Create saves user into database
func (s Storer) Create(ctx context.Context, user authboss.User) error {
	usr := user.(*User)
	u := &User{}

	findStatement := `SELECT * FROM users WHERE lower(email) = lower($1) OR name = $2`
	err := s.Db.Get(u, findStatement, usr.Email, usr.Name)
	if err != sql.ErrNoRows && err != nil {
		return err
	}
	if u.ID != 0 {
		return authboss.ErrUserFound
	}

	// Create user if OK
	insertStatement := `INSERT INTO users (
		name,
		email,
		password,
		confirmed,
		confirm_selector,
		confirm_verifier,
		updated_at,
		created_at) VALUES
		  (:name, lower(:email), :password, :confirmed, :confirm_selector, :confirm_verifier, NOW(), NOW())`

	_, err = s.Db.NamedExec(insertStatement,
		map[string]interface{}{
			"name":             usr.Name,
			"email":            usr.Email,
			"confirmed":        usr.Confirmed,
			"password":         usr.Password,
			"confirm_selector": usr.ConfirmSelector,
			"confirm_verifier": usr.ConfirmVerifier,
		})

	return err
}

// LoadByConfirmSelector implements logic of confirmation: loads user by confirm hash
func (s Storer) LoadByConfirmSelector(ctx context.Context, selector string) (authboss.ConfirmableUser, error) {
	u := &User{}

	findStatement := `SELECT * FROM users WHERE confirm_selector = $1`
	err := s.Db.Get(u, findStatement, selector)
	if err == sql.ErrNoRows {
		return nil, authboss.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return u, nil
}

// AddRememberToken adds remeber token to pid
func (s Storer) AddRememberToken(ctx context.Context, pid, token string) error {
	return s.redis.SAdd(ctx, "remember_tokens:"+pid, token).Err()
}

// UseRememberToken checks given token
func (s Storer) UseRememberToken(ctx context.Context, pid, token string) error {
	ok, err := s.redis.SIsMember(ctx, "remember_tokens:"+pid, token).Result()
	if err != nil {
		return err
	}
	if !ok {
		return authboss.ErrTokenNotFound
	}
	return nil
}

// DelRememberTokens removes all tokens for the given pid
func (s Storer) DelRememberTokens(ctx context.Context, pid string) error {
	return s.redis.Del(ctx, "remember_tokens:"+pid).Err()
}

// LoadByRecoverSelector loads user from db by recover token
func (s Storer) LoadByRecoverSelector(ctx context.Context, selector string) (authboss.RecoverableUser, error) {
	u := &User{}

	findStatement := `SELECT * FROM users WHERE recover_selector = $1`
	err := s.Db.Get(u, findStatement, selector)
	if err == sql.ErrNoRows {
		return nil, authboss.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return u, nil
}
