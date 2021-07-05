package structuredhttp

import (
	"testing"
	"time"
)

func TestRouteHandler(t *testing.T) {
	Timeout := time.Second * 10
	handler := RouteHandler{
		BaseURL: "https://httpstat.us",
		Timeout: &Timeout,
	}

	response, err := handler.GET("/200").Run()

	if err != nil {
		t.Error(err.Error())
		return
	}

	err = response.RaiseForStatus()
	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Log("Handler works!")
}
