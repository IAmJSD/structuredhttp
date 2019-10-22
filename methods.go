package structuredhttp

// GET adds support for a GET request.
func GET(URL string) *Request {
	return &Request{
		URL:            URL,
		Method:         "GET",
		Headers: map[string]string{},
	}
}

// POST adds support for a POST request.
func POST(URL string) *Request {
	return &Request{
		URL:            URL,
		Method:         "POST",
		Headers: map[string]string{},
	}
}

// PUT adds support for a PUT request.
func PUT(URL string) *Request {
	return &Request{
		URL:            URL,
		Method:         "PUT",
		Headers: map[string]string{},
	}
}

// PATCH adds support for a PATCH request.
func PATCH(URL string) *Request {
	return &Request{
		URL:            URL,
		Method:         "PATCH",
		Headers: map[string]string{},
	}
}

// DELETE adds support for a DELETE request.
func DELETE(URL string) *Request {
	return &Request{
		URL:            URL,
		Method:         "DELETE",
		Headers: map[string]string{},
	}
}

// OPTIONS adds support for a OPTIONS request.
func OPTIONS(URL string) *Request {
	return &Request{
		URL:            URL,
		Method:         "OPTIONS",
		Headers: map[string]string{},
	}
}

// HEAD adds support for a HEAD request.
func HEAD(URL string) *Request {
	return &Request{
		URL:            URL,
		Method:         "HEAD",
		Headers: map[string]string{},
	}
}
