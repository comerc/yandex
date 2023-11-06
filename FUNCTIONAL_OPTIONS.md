Паттерн функциональных опций

```go
package main

import (
	"errors"
	"net/http"
)

type options struct {
	port *int
}

type Option func(options *options) error

func WithPort(port int) Option {
	return func(options *options) error {
		if port < 0 {
			return errors.New("Booo")
		}
		options.port = &port
		return nil
	}
}

func WithPort(port int) Option {
	return func(options *options) error {
		if port < 0 {
			return errors.New("Booo")
		}
		options.port = &port
		return nil
	}
}

func NewServer(addr string, opts ...Option) (*http.Server, error) {

	var options options
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			return nil, err
		}
	}

	// ...
}

func main() {
	NewServer("localhost", WithPort(8080), WithTimeout(time.Second))
}
```
