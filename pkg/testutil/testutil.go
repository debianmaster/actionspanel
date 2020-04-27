// Package testutil has util functions for writing tests
package testutil

import "net/http"

// RoundTripFunc wraps a func definition with a type, so that we can attach a
// function receiver on it.
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip is a function we use to satisfy the http.RoundTripper interface so
// that we can intercept an http client's request and return some mock data.
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient creates an *http.Client where we can pass in a RoundTripper
// and intercept the request so that we can return mock data.
//
// The GitHub API Client provided by Google does not implement an interface, so
// mocking the client became a huge challenge. However, one of the ways we can
// create a github.Client is to create one from an existing http client. That's
// where this utility comes in. We create a new http.Client where we pass in an
// intercepting RoundTripper function where we force the RoundTripper to return
// mock data.
//
// Taken from:
// http://hassansin.github.io/Unit-Testing-http-client-in-Go
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}
