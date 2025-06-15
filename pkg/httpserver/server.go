// Package httpserver implements HTTP server.
package httpserver

import (
	"net"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

const (
// _defaultAddr            = ":80"
// _defaultReadTimeout     = 5 * time.Second
// _defaultWriteTimeout    = 5 * time.Second
// _defaultShutdownTimeout = 3 * time.Second
)

// Option -.
type OptionFiber func(*ServerFiber)

// Port -.
func PortFiber(port string) OptionFiber {
	return func(s *ServerFiber) {
		s.address = net.JoinHostPort("", port)
	}
}

// Prefork -.
// Prefork -.
func Prefork(prefork bool) OptionFiber {
	return func(s *ServerFiber) {
		s.prefork = prefork
	}
}

// ReadTimeout -.
func ReadTimeoutc(timeout time.Duration) OptionFiber {
	return func(s *ServerFiber) {
		s.readTimeout = timeout
	}
}

// // WriteTimeout -.
func WriteTimeoutOptionFiber(timeout time.Duration) OptionFiber {
	return func(s *ServerFiber) {
		s.writeTimeout = timeout
	}
}

// ShutdownTimeout -.
func ShutdownTimeoutOptionFiber(timeout time.Duration) OptionFiber {
	return func(s *ServerFiber) {
		s.shutdownTimeout = timeout
	}
}

// Server -.
type ServerFiber struct {
	App    *fiber.App
	notify chan error

	address         string
	prefork         bool
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration
}

// New -.
func NewFiber(opts ...OptionFiber) *ServerFiber {
	s := &ServerFiber{
		App:             nil,
		notify:          make(chan error, 1),
		address:         _defaultAddr,
		readTimeout:     _defaultReadTimeout,
		writeTimeout:    _defaultWriteTimeout,
		shutdownTimeout: _defaultShutdownTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	app := fiber.New(fiber.Config{
		Prefork:      s.prefork,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		JSONDecoder:  json.Unmarshal,
		JSONEncoder:  json.Marshal,
	})

	s.App = app

	return s
}

// Start -.
func (s *ServerFiber) StartFiber() {
	go func() {
		s.notify <- s.App.Listen(s.address)
		close(s.notify)
	}()
}

// Notify -.
func (s *ServerFiber) NotifyFiber() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *ServerFiber) ShutdownFiber() error {
	return s.App.ShutdownWithTimeout(s.shutdownTimeout)
}
