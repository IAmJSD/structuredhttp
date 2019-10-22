package structuredhttp

import "time"

// DefaultTimeout defines the default timeout.
var DefaultTimeout = time.Duration(0)

// SetDefaultTimeout allows the user to set the default timeout in this library.
func SetDefaultTimeout(Timeout time.Duration) {
	DefaultTimeout = Timeout
}
