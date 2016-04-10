package main

import (
	"errors"
	"net/http"
	"time"
)

// SendRequests sends multiple http requests asynchronously and then returns the result
func SendRequests(requests []*http.Request, resultCh chan error, timeoutSec int, handleResponse func(*http.Response) error) {
	client := &http.Client{}
	timeoutChan := time.After(time.Duration(timeoutSec) * time.Second)

	responseChan := make(chan error)
	for _, req := range requests {
		go func(req *http.Request) {
			resp, err := client.Do(req)
			if err != nil {
				responseChan <- err
			} else if err = handleResponse(resp); err != nil {
				responseChan <- err
			} else {
				responseChan <- nil
			}
		}(req)
	}

	for _, _ = range requests {
		select {
		case err := <-responseChan:
			if err != nil {
				resultCh <- err
				return
			}
		case <-timeoutChan:
			resultCh <- errors.New("SendRequests timed out!")
			return
		}
	}
	resultCh <- nil
}
