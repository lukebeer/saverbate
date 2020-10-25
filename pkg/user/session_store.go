package user

import (
	"encoding/base64"
	"time"

	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	abclientstate "github.com/volatiletech/authboss-clientstate"
	"github.com/volatiletech/authboss/v3"
)

const (
	sessionCookieName = "saverbate_session"
	sessionTTL        = int((14 * 24 * time.Hour) / time.Second)
)

// SessionStore represents session store
type SessionStore struct {
	CookieStorer  authboss.ClientStateReadWriter
	SessionStorer authboss.ClientStateReadWriter
}

// NewSessionStore builds SessionStore
func NewSessionStore() *SessionStore {
	cookieStoreKey, _ := base64.StdEncoding.DecodeString(viper.GetString("cookieStoreKey"))
	sessionStoreKey, _ := base64.StdEncoding.DecodeString(viper.GetString("sessionStoreKey"))

	cookieStore := abclientstate.NewCookieStorer(cookieStoreKey, nil)
	cookieStore.HTTPOnly = false
	cookieStore.Secure = true

	sessionStore := abclientstate.NewSessionStorer(sessionCookieName, sessionStoreKey, nil)
	cstore := sessionStore.Store.(*sessions.CookieStore)
	cstore.Options.HttpOnly = false
	cstore.Options.Secure = true
	cstore.MaxAge(sessionTTL)

	return &SessionStore{CookieStorer: cookieStore, SessionStorer: sessionStore}
}
