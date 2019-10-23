package structuredhttp

import (
	"net/url"
	"strings"
	"time"
)

// RouteHandler defines the base HTTP URL/timeout which is used for routes.
type RouteHandler struct {
	BaseURL string             `json:"base_url"`
	Timeout *time.Duration     `json:"-"`
	Headers *map[string]string `json:"headers"`
}

// GenerateURL takes a path and returns the URL with the path added.
func (r *RouteHandler) GenerateURL(Path string) (string, error) {
	u, err := url.Parse(r.BaseURL)
	if err != nil {
		return "", err
	}
	u.Path = strings.TrimRight(u.Path, "/") + Path
	return u.String(), nil
}

// GET does a GET request based on this base.
func (r *RouteHandler) GET(Path string) *Request {
	url, err := r.GenerateURL(Path)
	if err != nil {
		url = ""
	}
	req := GET(url)
	if err != nil {
		req.Error = &err
		return req
	}
	if r.Timeout != nil {
		req = req.Timeout(*r.Timeout)
	}
	if r.Headers != nil {
		for k, v := range *r.Headers {
			req = req.Header(k, v)
		}
	}
	return req
}

// POST does a POST request based on this base.
func (r *RouteHandler) POST(Path string) *Request {
	url, err := r.GenerateURL(Path)
	if err != nil {
		url = ""
	}
	req := POST(url)
	if err != nil {
		req.Error = &err
		return req
	}
	if r.Timeout != nil {
		req = req.Timeout(*r.Timeout)
	}
	if r.Headers != nil {
		for k, v := range *r.Headers {
			req = req.Header(k, v)
		}
	}
	return req
}

// PUT does a PUT request based on this base.
func (r *RouteHandler) PUT(Path string) *Request {
	url, err := r.GenerateURL(Path)
	if err != nil {
		url = ""
	}
	req := PUT(url)
	if err != nil {
		req.Error = &err
		return req
	}
	if r.Timeout != nil {
		req = req.Timeout(*r.Timeout)
	}
	if r.Headers != nil {
		for k, v := range *r.Headers {
			req = req.Header(k, v)
		}
	}
	return req
}

// PATCH does a PATCH request based on this base.
func (r *RouteHandler) PATCH(Path string) *Request {
	url, err := r.GenerateURL(Path)
	if err != nil {
		url = ""
	}
	req := PATCH(url)
	if err != nil {
		req.Error = &err
		return req
	}
	if r.Timeout != nil {
		req = req.Timeout(*r.Timeout)
	}
	if r.Headers != nil {
		for k, v := range *r.Headers {
			req = req.Header(k, v)
		}
	}
	return req
}

// DELETE does a DELETE request based on this base.
func (r *RouteHandler) DELETE(Path string) *Request {
	url, err := r.GenerateURL(Path)
	if err != nil {
		url = ""
	}
	req := DELETE(url)
	if err != nil {
		req.Error = &err
		return req
	}
	if r.Timeout != nil {
		req = req.Timeout(*r.Timeout)
	}
	if r.Headers != nil {
		for k, v := range *r.Headers {
			req = req.Header(k, v)
		}
	}
	return req
}

// OPTIONS does a OPTIONS request based on this base.
func (r *RouteHandler) OPTIONS(Path string) *Request {
	url, err := r.GenerateURL(Path)
	if err != nil {
		url = ""
	}
	req := OPTIONS(url)
	if err != nil {
		req.Error = &err
		return req
	}
	if r.Timeout != nil {
		req = req.Timeout(*r.Timeout)
	}
	if r.Headers != nil {
		for k, v := range *r.Headers {
			req = req.Header(k, v)
		}
	}
	return req
}

// HEAD does a HEAD request based on this base.
func (r *RouteHandler) HEAD(Path string) *Request {
	url, err := r.GenerateURL(Path)
	if err != nil {
		url = ""
	}
	req := HEAD(url)
	if err != nil {
		req.Error = &err
		return req
	}
	if r.Timeout != nil {
		req = req.Timeout(*r.Timeout)
	}
	if r.Headers != nil {
		for k, v := range *r.Headers {
			req = req.Header(k, v)
		}
	}
	return req
}
