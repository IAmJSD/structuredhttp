// +build js

package structuredhttp

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"syscall/js"
)

// HACK: Promises are not explicitly supported in Go right now, so we have to improvise.
func promiseHack(CalledPromise js.Value) (js.Value, error) {
	// Defines the result channel.
	resultChan := make(chan interface{})

	CalledPromise.Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Set the catch to set the error.
		if len(args) == 0 {
			resultChan <- errors.New("")
		} else {
			var msg string
			msgAttr := args[0].Get("message")
			if msgAttr.IsUndefined() {
				msg = js.Global().Call("String", args[0]).String()
			} else {
				msg = msgAttr.String()
			}
			resultChan <- errors.New(msg)
		}
		return nil
	})).Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// If there was no error, set the result.
		resultChan <- args[0]
		return nil
	}))

	// Get the promise result.
	res := <-resultChan
	switch x := res.(type) {
	case error:
		return js.Value{}, x
	case js.Value:
		return x, nil
	default:
		panic("internal error - unknown type")
	}
}

func strmap2obj(m map[string]string) js.Value {
	obj := js.Global().Call("Object")
	for k, v := range m {
		obj.Set(k, v)
	}
	return obj
}

func fetch(URL string, Args map[string]interface{}) (js.Value, error) {
	// Create a new object.
	obj := js.Global().Call("Object")

	// Set the objects.
	for k, v := range Args {
		obj.Set(k, v)
	}

	// Call the fetch API.
	return promiseHack(js.Global().Call("fetch", URL, obj))
}

func createSignal(ms int64) js.Value {
	// Create a instance of AbortController.
	controller := js.Global().Get("AbortController").New()

	// Run the setTimeout API on the controller abort function.
	js.Global().Call("setTimeout", controller.Get("abort"), ms)

	// Return the controllers signal.
	return controller.Get("signal")
}

func createReadableStream(r io.Reader) js.Value {
	// Create a new object.
	obj := js.Global().Call("Object")

	// Create a start function to start the data stream.
	start := func(controller js.Value) js.Func {
		var f js.Func
		f = js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
			b := make([]byte, 100)
			n, err := r.Read(b)
			if err != nil {
				controller.Call("close")
				return js.Undefined()
			}
			u := js.Global().Get("Uint8Array").New(n)
			js.CopyBytesToJS(u, b)
			controller.Call("enqueue", u)
			return f
		})
		return f
	}

	// Set the start attribute on the object.
	obj.Set("start", js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		return start(args[0])
	}))

	// Create the ReadableStream object.
	return js.Global().Get("ReadableStream").New(obj)
}

func goHeaders(fetch js.Value) http.Header {
	h := http.Header{}
	fetch.Get("headers").Call("forEach", js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		value := args[0]
		name := args[1]
		h.Set(name.String(), value.String())
		return js.Global().Get("true")
	}))
	return h
}

type blobReader struct {
	streamReader js.Value
}

// Read is used to read bytes like a normal IO reader in Go.
func (b *blobReader) Read(p []byte) (n int, err error) {
	result, err := promiseHack(b.streamReader.Call("read"))
	if err != nil {
		return 0, err
	}
	done := result.Get("done").Bool()
	if done {
		return 0, nil
	}
	return js.CopyBytesToGo(p, result.Get("value")), nil
}

func fetch2http(fetch js.Value) *http.Response {
	// Get the blob object.
	blob, err := promiseHack(fetch.Call("blob"))
	if err != nil {
		// Hmmmmmm, this is a bug.
		panic(err)
	}

	// Get the size and reader.
	var Reader io.ReadCloser
	Size := int64(blob.Get("size").Int())
	if Size == 0 {
		Reader = ioutil.NopCloser(strings.NewReader(""))
	} else {
		Reader = ioutil.NopCloser(&blobReader{streamReader: blob.Call("stream").Call("getReader")})
	}

	// Return the HTTP response.
	return &http.Response{
		Status:           fetch.Get("statusText").String(),
		StatusCode:       fetch.Get("status").Int(),
		Header:           goHeaders(fetch),
		Body:             Reader,
		ContentLength:    Size,
	}
}

// Run executes the request.
func (r *Request) Run() (*Response, error) {
	// Handle previous errors.
	if r.Error != nil {
		return nil, *r.Error
	}

	// Create the AbortController signal if needed.
	Signal := js.Undefined()
	if r.CurrentTimeout == nil && DefaultTimeout != 0 {
		Signal = createSignal(DefaultTimeout.Milliseconds())
	} else if r.CurrentTimeout != nil && *r.CurrentTimeout != 0 {
		Signal = createSignal(r.CurrentTimeout.Milliseconds())
	}

	// Defines the fetch arguments.
	Reader := r.CurrentReader
	if Reader == nil {
		Reader = strings.NewReader("")
	}
	FetchArgs := map[string]interface{}{
		"signal": Signal,
		"method": r.Method,
		"headers": strmap2obj(r.Headers),
		"body": createReadableStream(r.CurrentReader),
	}
	if r.Method == "GET" || r.Method == "HEAD" {
		delete(FetchArgs, "body")
	}

	// Call fetch.
	res, err := fetch(r.URL, FetchArgs)
	if err != nil {
		return nil, err
	}

	// Create the response object.
	return &Response{RawResponse: fetch2http(res)}, nil
}
