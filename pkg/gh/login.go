package gh

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs"
	"github.com/google/go-github/v30/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/palantir/go-githubapp/oauth2"
	"github.com/phunki/actionspanel/pkg/constants"
	"github.com/phunki/actionspanel/pkg/log"
)

// LoginHandler creates a GitHub oauth2 handler. It's convenience more than
// anything.
type LoginHandler struct {
	sessionManager *scs.Manager
	tls            bool
	githubConfig   githubapp.Config
}

// NewLoginHandler creates a new GitHub oauth2 handler
func NewLoginHandler(sessionManager *scs.Manager, tls bool, githubConfig githubapp.Config) LoginHandler {
	return LoginHandler{sessionManager: sessionManager, tls: tls, githubConfig: githubConfig}
}

// NewOAuthHandler returns a new handler that we can register to handle an OAuth flow
func (h LoginHandler) NewOAuthHandler() http.Handler {
	return oauth2.NewHandler(
		oauth2.GetConfig(h.githubConfig, nil),
		// force generated URLs to use HTTPS; useful if the app is behind a reverse proxy
		oauth2.ForceTLS(h.tls),
		// set the callback for successful logins
		oauth2.OnLogin(h.onLogin),
		oauth2.OnError(h.onError),
		oauth2.WithStore(&oauth2.SessionStateStore{Sessions: h.sessionManager}),
	)
}

// MapInstallationIDs retrieves installation IDs and maps names to installationIDs
func MapInstallationIDs(appsService AppsService) (map[string]int64, error) {
	allAvailableInstallations := make([]*github.Installation, 0)
	opt := &github.ListOptions{}

	for {
		availableInstallations, resp, err := appsService.ListUserInstallations(context.Background(), opt)
		if resp.StatusCode != 200 || err != nil {
			log.Err(err, "couldn't list user's installations")
			return nil, err
		}
		allAvailableInstallations = append(allAvailableInstallations, availableInstallations...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	installationMap := make(map[string]int64)
	for _, installation := range allAvailableInstallations {
		installationMap[installation.GetAccount().GetLogin()] = installation.GetID()
	}
	return installationMap, nil
}

func (h LoginHandler) onLogin(w http.ResponseWriter, r *http.Request, login *oauth2.Login) {
	log.Info("Handling a new oauth login")

	session := h.sessionManager.Load(r)
	client := github.NewClient(login.Client)
	installationMap, err := MapInstallationIDs(client.Apps)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Reset installation map
	err = session.Remove(w, constants.InstallationMap)
	if err != nil {
		log.Err(err, "failed to clear installation map from session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = session.PutObject(w, constants.InstallationMap, installationMap)
	if err != nil {
		log.Err(err, "failed to put installation map into session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = session.PutString(w, constants.AccessToken, login.Token.AccessToken)
	if err != nil {
		log.Err(err, "failed to put access token into session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h LoginHandler) onError(w http.ResponseWriter, r *http.Request, err error) {
	log.Infof("Attempted url: %v", r.URL)
	log.Err(err, "couldn't login")
}
