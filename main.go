package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	logger "github.com/rtfmkiesel/kisslog"
	flag "github.com/spf13/pflag"

	"github.com/rtfmkiesel/umami-importer/config"
	"github.com/rtfmkiesel/umami-importer/db"
	"github.com/rtfmkiesel/umami-importer/umami"
)

var version string = "@DEV" // Adjusted by the Makefile

func main() {
	if err := logger.InitDefault("umami-importer" + version); err != nil {
		panic(err)
	}
	log := logger.New("main")

	configPath := flag.StringP("config", "c", "./config.yaml", "config path")
	flag.BoolVarP(&logger.FlagDebug, "verbose", "v", false, "enable verbose/debug output")
	flag.Parse()

	cfg, err := config.LoadFromFile(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	client := umami.NewClient(cfg.Umami)

	if err := db.Open(cfg.Database.Path); err != nil {
		log.Fatal(err)
	}
	defer db.Close() //nolint:errcheck

	ctx, stopImport := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		defer close(done)
		for _, c := range cfg.Imports {
			select {
			case <-ctx.Done():
				return
			default:
				if err := client.Import(c, ctx); err != nil {
					log.Error("Failed import for '%s': %s", c.Website.BaseURL, err)
				}
			}
		}
	}()

	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-chanSignal:
		log.Warning("Received stop signal, waiting for current imports to finish...")
		// Sends a trigger to the import functions
		// to stop after the current file
		stopImport()

		<-done
	case <-done:
		log.Info("Done")
	}
}
