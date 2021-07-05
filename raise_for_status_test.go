package structuredhttp

import (
	"testing"
	"time"
)

func TestRaiseForStatus(t *testing.T) {
	response, err := GET(
		"https://httpstat.us/403").Timeout(10 * time.Second).Run()

	if err != nil {
		t.Error(err.Error())
		return
	}

	err = response.RaiseForStatus()
	if err == nil {
		t.Error("Failed to raise for a status 403.")
	} else {
		t.Log("Successfully raised for a status 403.")
	}
}
