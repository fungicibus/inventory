package main

import (
	"context"
	"embed"
	"encoding/json"
	"os"
	"os/signal"

	"github.com/fungicibus/inventory/config"
	v1 "github.com/fungicibus/inventory/internal/api/v1"
	"github.com/fungicibus/inventory/internal/logger"
	"github.com/fungicibus/inventory/internal/server"
	"github.com/fungicibus/inventory/internal/storage/migrations"
	"github.com/fungicibus/inventory/internal/storage/pg"
)

var Tag string
var Commit string

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	log := logger.New()

	version := getVersion()

	cfg, err := config.GetDefault()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get config")
	}
	log.SetLevel(cfg.LogLevel)
	cfg.AppVersion = version

	cfgContent, _ := json.Marshal(cfg)
	log.Debug().RawJSON("config", cfgContent).Msg("config")

	pg, err := pg.New(cfg.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create postgres adapter")
	}
	if err := migrations.Up(embedMigrations, pg.DB()); err != nil {
		log.Fatal().Err(err).Msg("failed to up migrations")
	}

	v1 := v1.New(cfg, log, pg)

	server := server.New(cfg, log, v1.GetHandler())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	go func() {
		err := server.Run(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	<-ctx.Done()
	pg.Close()
	server.Shutdown()
}

func getVersion() string {
	tag, commit := Tag, Commit

	if Tag == "" {
		tag = "tag"
	}
	if Commit == "" {
		commit = "commit"
	}
	return tag + "-" + commit
}
