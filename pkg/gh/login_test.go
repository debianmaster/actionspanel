package gh_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v30/github"
	"github.com/phunki/actionspanel/mock"
	"github.com/phunki/actionspanel/pkg/gh"
	"github.com/stretchr/testify/assert"
)

func Test_Login_MapInstallationIDs_SingleInstallation(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAppsService := mock.NewMockAppsService(ctrl)

	mockInstallations := []*github.Installation{
		&github.Installation{
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
		&github.Installation{
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
		&github.Installation{
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
