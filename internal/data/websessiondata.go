//session for website cookies and allthat

package data

import (
	"crypto/sha256"
	"database/sql"

	"net/http"

	"github.com/cohune-cabbage/di/internal/validator"
	"github.com/gorilla/sessions"
)

// SessionManager manages user sessions.
type SessionManager struct {
	DB        *sql.DB
	Store     *sessions.Cooki8eStore
	CSRFKey   []byte
	CSRFToken string
	Validator *validator.Validator
}

func NewSessionManager(db *sql.DB, secretKey string) *SessionManager {
	// Hash the secret key using SHA-256
	hashedKey := sha256.Sum256([]byte(secretKey))

	store := sessions.NewCookieStore(hashedKey[:])
	store.Options = &sessions.Options{
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
		Secure:   true, // Set to true in production
		SameSite: http.SameSiteStrictMode,
	}

	return &SessionManager{
		DB:        db,
		Store:     store,
		CSRFKey:   hashedKey[:],
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
	session.Values[""] = m.CSRFToken
	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

// exists check if the session exists
func (m *SessionManager) Exists(r *http.Request, key string) bool {
	session, err := m.Store.Get(r, "session")
	if err != nil {
		return false
	}
	_, ok := session.Values[key]
	return ok
}

// get string from session
func (m *SessionManager) GetString(r *http.Request, key string) string {
	session, err := m.Store.Get(r, "session")
	if err != nil {
		return ""
	}
	value, ok := session.Values[key].(string)
	if !ok {
		return ""
	}
	return value
}

// get int from session
func (m *SessionManager) GetInt(r *http.Request, key string) int {
	session, err := m.Store.Get(r, "session")
	if err != nil {
		return 0
	}
	value, ok := session.Values[key].(int)
	if !ok {
		return 0
	}
	return value
}

// get string from session
func (m *SessionManager) Get(r *http.Request, key string) interface{} {
	session, err := m.Store.Get(r, "session")
	if err != nil {
		return nil
	}
	value, ok := session.Values[key]
	if !ok {
		return nil
	}
	return value
}

// get school id from session
func (m *SessionManager) GetSchools(r *http.Request) []int {
	session, err := m.Store.Get(r, "session")
	if err != nil {
		return nil
	}
	value, ok := session.Values["school_id"].([]int)
	if !ok {
		return nil
	}
	return value
}
