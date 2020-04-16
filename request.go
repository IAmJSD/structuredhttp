package structuredhttp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Request defines the request that will be ran.
type Request struct {
	URL            string            `json:"url"`
	Method         string            `json:"method"`
	Headers        map[string]string `json:"headers"`
	CurrentTimeout *time.Duration    `json:"timeout"`
	CurrentReader  io.Reader         `json:"-"`
	Error          *error            `json:"-"`
}

// Header sets a header.
func (r *Request) Header(key string, value string) *Request {
	if r.Error != nil {
		return r
	}
	r.Headers[key] = value
	return r
}

// Timeout sets a timeout. 0 is infinite.
func (r *Request) Timeout(Value time.Duration) *Request {
	if r.Error != nil {
		return r
	}
	r.CurrentTimeout = &Value
	return r
}

// Bytes sets the data to the bytes specified.
func (r *Request) Bytes(Data []byte) *Request {
	if r.Error != nil {
		return r
	}
	Buffer := bytes.NewReader(Data)
	r.Headers["Content-Length"] = strconv.Itoa(len(Data))
	r.CurrentReader = Buffer
	return r
}

// JSON sets the data to the JSON specified.
func (r *Request) JSON(Data interface{}) *Request {
	if r.Error != nil {
		return r
	}
	JSONData, err := json.Marshal(Data)
	if err != nil {
		r.Error = &err
		return r
	}
	r.Headers["Content-Length"] = strconv.Itoa(len(JSONData))
	r.Headers["Content-Type"] = "application/json"
	r.Bytes(JSONData)
	return r
}

// Reader sets the IO reader to the one specified.
func (r *Request) Reader(Data io.Reader) *Request {
	if r.Error != nil {
		return r
	}
	r.CurrentReader = Data
	return r
}

// URLEncodedForm sets the data to a URL encoded form.
func (r *Request) URLEncodedForm(Data url.Values) *Request {
	if r.Error != nil {
		return r
	}
	Encoded := Data.Encode()
	r.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	r.Headers["Content-Length"] = strconv.Itoa(len(Encoded))
	r.CurrentReader = strings.NewReader(Encoded)
	return r
}

// MultipartForm sets the data to a multipart form.
func (r *Request) MultipartForm(Buffer *bytes.Buffer, ContentType string) *Request {
	if r.Error != nil {
		return r
	}
	r.Headers["Content-Type"] = ContentType
	r.CurrentReader = Buffer
	return r
}

// Plugin allows for third party functions to be chained into the request.
func (r *Request) Plugin(Function func(r *Request)) *Request {
	if r.Error != nil {
		return r
	}
	Function(r)
	return r
}

// Query adds a URL query argument to the URL.
func (r *Request) Query(Key string, Value string) *Request {
	if r.Error != nil {
		return r
	}
	u, err := url.Parse(r.URL)
	if err != nil {
		r.Error = &err
		return r
	}
	q := u.Query()
	q.Set(Key, Value)
	u.RawQuery = q.Encode()
	r.URL = u.String()
	return r
}
