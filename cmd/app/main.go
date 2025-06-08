package main

import (
	"log"

	"github.com/ducnpdev/godev-kit/config"
	"github.com/ducnpdev/godev-kit/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
