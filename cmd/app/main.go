package main

import (
	"flag"
	"log"
	"os"

	"github.com/antonmisa/1cctl/config"
	"github.com/antonmisa/1cctl/internal/app"
)

func main() {
	var prepare bool

	flag.BoolVar(&prepare, "prepare", false, "creating default environment and config")

	flag.Parse()

	// Just prepare env, config and exit
	if prepare {
		err := config.Prepare()
		if err != nil {
			log.Fatalf("Prepare error: %s", err)
		}
		os.Exit(0)
	}

	// Configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run -.
	app.Run(cfg)
}
