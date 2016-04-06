package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendRequest(t *testing.T) {
	testCreateHttp := func(method string, url string) *http.Request {
		req, _ := http.NewRequest(method, url, nil)
		return req
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	var sendRequestTests = []struct {
		requests []*http.Request
		noError  bool
	}{
		{
			[]*http.Request{
				testCreateHttp("GET", "http://www.google.com"),
				testCreateHttp("GET", "http://www.facebook.com"),
			},
			true,
		},
		{
			[]*http.Request{
				testCreateHttp("GET", "http://www.google.com"),
				testCreateHttp("GET", server.URL),
			},
			false,
		},
	}

	for _, test := range sendRequestTests {
		resultCh := make(chan error)
		go SendRequests(test.requests, resultCh, 3, func(resp *http.Response) error {
			if resp.StatusCode != 200 {
				return errors.New("Response status code is not 200!")
			}

			return nil
		})
		err := <-resultCh
		if err != nil && test.noError {
			t.Errorf("TestSendRequest failed: Error when sending %v", err)
		} else if err == nil && !test.noError {
			t.Errorf("TestSendRequest failed: Expected to return an error")
		}
	}
}
