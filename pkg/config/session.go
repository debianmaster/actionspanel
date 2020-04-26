package config

import "github.com/alexedwards/scs"

// SessionManagerFactory creates a new session manager
type SessionManagerFactory struct {
	sessionManagerType  string
	cookieSessionSecret string
}

// CreateSessionManager will create a session manager based on the type of session manager it is configured to make
func (f *SessionManagerFactory) CreateSessionManager() *scs.Manager {
	switch f.sessionManagerType {
	case "cookie":
		return scs.NewCookieManager(f.cookieSessionSecret)
	}
	return nil
}

// NewSessionManagerFactory creates a session factory
//
// This function panics if an invalid session type is passed in. We want this
// behavior since having a valid session manager is crucial to the application operation.
func NewSessionManagerFactory(cfg Config) *SessionManagerFactory {
	switch cfg.SessionManagerType {
	case "cookie":
		if cfg.CookieSessionSecret == "" {
			return &SessionManagerFactory{sessionManagerType: "cookie", cookieSessionSecret: DefaultCookieSessionSecret}
		}
		return &SessionManagerFactory{sessionManagerType: "cookie", cookieSessionSecret: cfg.CookieSessionSecret}
	default:
		panic("session manager type is unrecognized")
	}
}
