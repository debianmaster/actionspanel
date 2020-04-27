package gh_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/alexedwards/scs"
	"github.com/alexedwards/scs/stores/memstore"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v30/github"
	"github.com/palantir/go-githubapp/oauth2"
	"github.com/phunki/actionspanel/mock"
	"github.com/phunki/actionspanel/pkg/constants"
	"github.com/phunki/actionspanel/pkg/gh"
	"github.com/phunki/actionspanel/pkg/testutil"
	"github.com/stretchr/testify/assert"
	goauth2 "golang.org/x/oauth2"
)

func Test_Login_MapInstallationIDs_SingleInstallation(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAppsService := mock.NewMockAppsService(ctrl)

	mockInstallations := []*github.Installation{
		{
			ID: github.Int64(12345),
			Account: &github.User{
				Login: github.String("abatilo"),
			},
		},
	}
	mockResponse := &github.Response{
		Response: &http.Response{
			StatusCode: 200,
		},
	}
	var mockError error

	expectedMap := map[string]int64{
		"abatilo": 12345,
	}

	mockAppsService.EXPECT().ListUserInstallations(gomock.Any(), gomock.Any()).Return(mockInstallations, mockResponse, mockError)

	actual, err := gh.MapInstallationIDs(mockAppsService)
	assert.Nil(err)
	assert.Equal(expectedMap, actual)
}

func Test_Login_MapInstallationIDs_Pagination(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAppsService := mock.NewMockAppsService(ctrl)

	// Page 1
	mockInstallationsPage1 := []*github.Installation{
		{
			ID: github.Int64(12345),
			Account: &github.User{
				Login: github.String("abatilo"),
			},
		},
	}
	mockResponsePage1 := &github.Response{
		NextPage: 1,
		Response: &http.Response{
			StatusCode: 200,
		},
	}

	// Page 2
	mockInstallationsPage2 := []*github.Installation{
		{
			ID: github.Int64(54321),
			Account: &github.User{
				Login: github.String("phunki"),
			},
		},
	}
	mockResponsePage2 := &github.Response{
		NextPage: 0,
		Response: &http.Response{
			StatusCode: 200,
		},
	}

	var mockError error

	expectedMap := map[string]int64{
		"abatilo": 12345,
		"phunki":  54321,
	}

	gomock.InOrder(
		// First
		mockAppsService.EXPECT().ListUserInstallations(gomock.Any(), gomock.Any()).
			Return(mockInstallationsPage1, mockResponsePage1, mockError).
			Times(1),

		// Second
		mockAppsService.EXPECT().ListUserInstallations(gomock.Any(), gomock.Any()).
			Return(mockInstallationsPage2, mockResponsePage2, mockError).
			Times(1),
	)

	actual, err := gh.MapInstallationIDs(mockAppsService)
	assert.Nil(err)
	assert.Equal(expectedMap, actual)
}

func Test_Login_MapInstallationIDs_AppsServiceErrNotNil(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAppsService := mock.NewMockAppsService(ctrl)

	mockResponse := &github.Response{
		Response: &http.Response{
			StatusCode: 200,
		}}

	mockAppsService.EXPECT().ListUserInstallations(gomock.Any(), gomock.Any()).
		Return(nil, mockResponse, errors.New("failure")).
		Times(1)

	actual, err := gh.MapInstallationIDs(mockAppsService)
	assert.NotNil(err)
	assert.Nil(actual)
}

func Test_Login_MapInstallationIDs_AppsServiceResponseNot200(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAppsService := mock.NewMockAppsService(ctrl)

	mockResponse := &github.Response{
		Response: &http.Response{
			StatusCode: 500,
		}}

	mockAppsService.EXPECT().ListUserInstallations(gomock.Any(), gomock.Any()).
		Return(nil, mockResponse, errors.New("GitHub API returned non 200")).
		Times(1)

	actual, err := gh.MapInstallationIDs(mockAppsService)
	assert.NotNil(err)
	assert.Nil(actual)
}

func Test_Login_LoginCallbackSetsSessionToken(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	installationsURL, _ := url.Parse("https://api.github.com/user/installations")

	mockClientCreator := mock.NewMockClientCreator(ctrl)
	mockInstallations := struct {
		Installations []*github.Installation `json:"installations"`
	}{
		Installations: []*github.Installation{
			{
				ID: github.Int64(12345),
				Account: &github.User{
					Login: github.String("abatilo"),
				},
			},
		},
	}
	mockBytes, _ := json.Marshal(mockInstallations)

	// Creates a new github.Client where we inject a custom RoundTripper to returning the mock data we want
	mockClient := github.NewClient(testutil.NewTestClient(func(req *http.Request) *http.Response {
		if req.Method == "GET" && req.URL.String() == installationsURL.String() {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBuffer(mockBytes)),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		t.Fatal("Installations request was never matched")
		return nil
	}))
	mockClientCreator.EXPECT().NewTokenClient(gomock.Any()).Return(mockClient, nil)

	// Create an in memory store just for the test environment
	sessionManager := scs.NewManager(memstore.New(1 * time.Second))

	// The scs library doesn't give us any way to create our own Session. All the
	// fields are private, so the only thing we can do is use the
	// scs.Manager.Load function which lazily creates and populates a Session. So
	// we create a session to use based on an empty http request which we never
	// use
	throwaway := &http.Request{}
	session := sessionManager.Load(throwaway)

	// Inject our new session into our mock request, so that we can assert on it later
	ctx := sessionManager.AddToContext(context.Background(), session)
	req, err := http.NewRequestWithContext(ctx, "GET", "/api/github/auth", nil)
	assert.Nil(err)

	// We need some AccessToken to inject into the mocked clientCreator
	login := &oauth2.Login{
		Token: &goauth2.Token{
			AccessToken: "fdsa",
		},
	}

	resp := httptest.NewRecorder()
	// Create the callback to test
	loginCallback := gh.NewOnLoginCallback(sessionManager, mockClientCreator)

	// Verify state before executing the callback
	assert.False(session.Exists(constants.AccessToken))
	loginCallback(resp, req, login)

	_, err = ioutil.ReadAll(resp.Result().Body)
	assert.Nil(err)
	defer resp.Result().Body.Close()

	// Verify state after executing the callback
	assert.Equal(http.StatusOK, resp.Result().StatusCode)
	assert.True(session.Exists(constants.AccessToken))
	accessToken, _ := session.GetString(constants.AccessToken)
	assert.Equal("fdsa", accessToken)
}
