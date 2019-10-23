# structuredhttp
A lightweight structured HTTP client for Go with no third party dependencies. The client works by chaining functions:
```go
package main

import (
	"github.com/jakemakesstuff/structuredhttp"
	"time"
)

func main() {
	response, err := structuredhttp.GET(
		"https://httpstat.us/200").Timeout(10 * time.Second).Run()

	if err != nil {
		panic(err)
	}

	err = response.RaiseForStatus()
	if err != nil {
		panic(err)
	}
 
	println("Ayyyy! A 200!")
}
```

## Making a request
To make a request, you need to call the function representing the HTTP method. For example, if you want to make a HTTP GET request, you will call the `GET` function. There are several functions you can call in the request chain:
- `Timeout` - Check the "Handling timeouts" documentation below.
- `Header` - This adds a header into your request. This function takes a key and a value.
- `Bytes` - Puts the bytes specified into the body.
- `JSON` - Serializes the item specified into JSON (**make sure you provide a pointer**) and puts it into the body.
- `Reader` - Allows you to provide your own I/O reader.
- `URLEncodedForm` - This will take the values specified and turn it into a URL encoded form.
- `MultipartForm` - This will take the buffer and content type after the creation of a multipart form and handle it.
- `Plugin` - This will pass through to a third party function specified. The plugin will need to take `*structuredhttp.Request` as an argument.

After you have made the request chain, you should call `Run`. This function will then return a pointer to the Response structure (described below) and an error which will not be null if something went wrong.

## Handling timeouts
There are 2 ways to handle timeouts:
1. **Call the `Timeout` function in the request chain:** Calling the timeout function in the request chain will override the default timeout.
2. **Set a default timeout:** To set a default timeout, you can call `SetDefaultTimeout`:
    ```go
    structuredhttp.SetDefaultTimeout(5 * time.Second)
    ```

## The Response structure
The response structure has several useful functions:
- `Bytes` - This returns the response as bytes.
- `JSON` - This returns the response as an interface which can be casted to different types.
- `RaiseForStatus` - This just returns an error. The error will not be null if it's a HTTP error.
- `Text` - This returns the response as text.

If you need the raw response, the `RawResponse` attribute contains a pointer to the `http.Response` from the request.

## Request error handling
The Request structure has an `Error` attribute. If there is an error, the error should be attached to this attribute. Any other functions in the chain will be skipped, and in the `Run` function the error will be thrown.
