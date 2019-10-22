package tests

import (
	"structuredhttp"
	"testing"
	"time"
)

func TestDefaultTimeout(t *testing.T) {
	structuredhttp.SetDefaultTimeout(5 * time.Second)
	t.Log("Default timeout successfully set.")
}
