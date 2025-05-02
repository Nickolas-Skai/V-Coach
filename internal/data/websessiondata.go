//session for website cookies and allthat

package data

import (
	"database/sql"

	"net/http"

	"github.com/cohune-cabbage/di/internal/validator"
	"github.com/gorilla/sessions"
)

// SessionManager manages user sessions.
type SessionManager struct {
	DB        *sql.DB
	Store     *sessions.CookieStore
	CSRFKey   []byte
	CSRFToken string
	Validator *validator.Validator
}

// NewSessionManager creates a new SessionManager.
func NewSessionManager(db *sql.DB, secretKey string) *SessionManager {
	store := sessions.NewCookieStore([]byte(secretKey))
	store.Options = &sessions.Options{}

	// Set the session expiration time to 1 hour

	store.Options.MaxAge = 3600
	store.Options.HttpOnly = true
	store.Options.Secure = true
	store.Options.SameSite = http.SameSiteStrictMode
	return &SessionManager{
		DB:        db,
		Store:     store,
		CSRFKey:   []byte(secretKey),
		Validator: validator.NewValidator(),
	}
}

// put session data in the session
func (m *SessionManager) Put(r *http.Request, w http.ResponseWriter, key string, value interface{}) error {
	session, err := m.Store.Get(r, "session")
	if err != nil {
		return err
	}
	session.Values[key] = value
	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

// renew token
func (m *SessionManager) RenewToken(r *http.Request, w http.ResponseWriter) error {
	session, err := m.Store.Get(r, "session")
	if err != nil {
		return err
	}
	session.Values["csrf_token"] = m.CSRFToken
	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}
