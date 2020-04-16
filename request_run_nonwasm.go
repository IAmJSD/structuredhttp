// +build !js

package structuredhttp

import (
	"net/http"
	"strings"
	"time"
)

// Run executes the request.
func (r *Request) Run() (*Response, error) {
	if r.Error != nil {
		return nil, *r.Error
	}
	var CurrentTimeout time.Duration
	if r.CurrentTimeout == nil {
		CurrentTimeout = DefaultTimeout
	} else {
		CurrentTimeout = *r.CurrentTimeout
	}
	Client := http.Client{
		Timeout: CurrentTimeout,
	}
	Reader := r.CurrentReader
	if Reader == nil {
		Reader = strings.NewReader("")
	}
	RawRequest, err := http.NewRequest(r.Method, r.URL, Reader)
	if err != nil {
		return nil, err
	}
	for k, v := range r.Headers {
		RawRequest.Header.Set(k, v)
	}
	RawResponse, err := Client.Do(RawRequest)
	if err != nil {
		return nil, err
	}
	return &Response{
		RawResponse: RawResponse,
	}, nil
}
