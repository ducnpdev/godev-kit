// Package httpserver implements HTTP server.
package httpserver

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/ducnpdev/godev-kit/config"
	"github.com/gin-gonic/gin"
)

const (
	_defaultAddr            = ":80"
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

// Option -.
type Option func(*Server)

// Port -.
func Port(port string) Option {
	return func(s *Server) {
		s.address = net.JoinHostPort("", port)
	}
}

// Prefork -.
// func Prefork(prefork bool) Option {
// 	return func(s *Server) {
// 		s.prefork = prefork
// 	}
// }

// ReadTimeout -.
func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = timeout
	}
}

// WriteTimeout -.
func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.writeTimeout = timeout
	}
}

// ShutdownTimeout -.
func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}

// Server -.
type Server struct {
	App    *gin.Engine
	notify chan error
	srv    *http.Server

	address         string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration
}

// // Option -.
// type Option func(*Server)

// // Port -.
// func Port(port string) Option {
// 	return func(s *Server) {
// 		s.address = port
// 	}
// }

// // ReadTimeout -.
// func ReadTimeout(timeout time.Duration) Option {
// 	return func(s *Server) {
// 		s.readTimeout = timeout
// 	}
// }

// // WriteTimeout -.
// func WriteTimeout(timeout time.Duration) Option {
// 	return func(s *Server) {
// 		s.writeTimeout = timeout
// 	}
// }

// // ShutdownTimeout -.
// func ShutdownTimeout(timeout time.Duration) Option {
// 	return func(s *Server) {
// 		s.shutdownTimeout = timeout
// 	}
// }

// New -.
func New(cfg *config.Config, opts ...Option) *Server {
	s := &Server{
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
	if cfg.App.MODE == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Set Gin mode

	// Create Gin engine with default middleware
	app := gin.New()

	// Add default middleware
	// app.Use(gin.Recovery())
	// app.Use(gin.Logger())

	s.App = app

	// Create HTTP server
	s.srv = &http.Server{
		Addr:         s.address,
		Handler:      s.App,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
	}

	return s
}

// Start -.
func (s *Server) Start() {
	go func() {
		s.notify <- s.srv.ListenAndServe()
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.srv.Shutdown(ctx)
}

// func errorResponse(c *gin.Context, code int, msg string) {
// 	c.JSON(code, response.Error{Error: msg})
// }
