package structuredhttp

import (
	"testing"
	"time"
)

func TestDefaultTimeout(t *testing.T) {
	SetDefaultTimeout(5 * time.Second)
	t.Log("Default timeout successfully set.")
}
