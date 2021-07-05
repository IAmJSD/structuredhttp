package structuredhttp

import (
	"sync"
	"sync/atomic"
)

// Batch is used to define a request batch.
type Batch []*Request

// All is used to execute a batch job and return the pure responses. If one fails, the first error is returned.
func (b Batch) All() ([]*Response, error) {
	// Get the results of each request.
	results := make([]interface{}, len(b))
	wg := sync.WaitGroup{}
	wg.Add(len(b))
	for i, v := range b {
		// Create a pointer to the result.
		ptr := &results[i]

		// Create a goroutine to do the request.
		go func(req *Request) {
			// Defer the request as done.
			defer wg.Done()

			// Run the request.
			res, err := req.Run()
			if err != nil {
				*ptr = err
				return
			}
			*ptr = res
		}(v)
	}
	wg.Wait()

	// Un-generic the results and return them.
	ungeneric := make([]*Response, len(results))
	for i, v := range results {
		if err, ok := v.(error); ok {
			return nil, err
		}
		ungeneric[i] = v.(*Response)
	}
	return ungeneric, nil
}

// BatchMapper is used to handle mapping a Batch to values.
type BatchMapper struct {
	b Batch
	f []func(interface{}) (interface{}, error)
}

// Map is used to handle adding a mapper function to the chain.
func (b *BatchMapper) Map(f func(interface{}) (interface{}, error)) *BatchMapper {
	b.f = append(b.f, f)
	return b
}

// Map is used to create a new batch mapper builder.
func (b Batch) Map(f func(*Response) (interface{}, error)) *BatchMapper {
	return &BatchMapper{
		b: b,
		f: []func(interface{}) (interface{}, error) {
			func(x interface{}) (interface{}, error) {
				return f(x.(*Response))
			},
		},
	}
}

// All is used to execute a batch job and return the mapped responses. If one fails, the first error is returned.
func (b *BatchMapper) All() ([]interface{}, error) {
	responses, err := b.b.All()
	if err != nil {
		return nil, err
	}
	done := make([]interface{}, len(responses))
	var errSet uintptr
	wg := sync.WaitGroup{}
	wg.Add(len(responses))
	for i, v := range responses {
		// Defines the interface for tracking the mapping.
		var iface interface{} = v

		// Get the pointer.
		ptr := &done[i]

		// Handle the rest for this value in a goroutine.
		go func() {
			// Defer marking this function as done.
			defer wg.Done()

			// Defines a locally scoped error.
			var locallyScopedErr error

			// Go through the chain.
			for _, f := range b.f {
				// Call the function and handle errors if they crop up.
				if iface, locallyScopedErr = f(iface); locallyScopedErr != nil {
					// Check we weren't raced and set error set to 1.
					if atomic.SwapUintptr(&errSet, 1) == 1 {
						// We were raced. Return here.
						return
					}

					// Set the error to the locally scoped one.
					err = locallyScopedErr

					// Return here.
					return
				}

				// Check if there is an error in another batch item. If so, return.
				if atomic.LoadUintptr(&errSet) != 0 {
					return
				}
			}

			// Set the pointer to the interface.
			*ptr = iface
		}()
	}
	wg.Wait()
	if errSet != 0 {
		return nil, err
	}
	return done, nil
}
