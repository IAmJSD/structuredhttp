package tests

import (
	"structuredhttp"
	"testing"
	"time"
)

func TestJSONPost(t *testing.T) {
	response, err := structuredhttp.POST(
		"https://postman-echo.com/post").Timeout(10 * time.Second).JSON(&map[string]string{
			"hello": "world",
	}).Run()

	if err != nil {
		t.Error(err.Error())
		return
	}

	err = response.RaiseForStatus()
	if err != nil {
		t.Error(err.Error())
		return
	}

	i, err := response.JSON()
	if err != nil {
		t.Error(err.Error())
		return
	}

	World := i.(map[string]interface{})["data"].(map[string]interface{})["hello"].(string)
	if World != "world" {
		t.Error("Invalid string returned (" + World + ").")
		return
	}
	t.Log("Success!")
}
