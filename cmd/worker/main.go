package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/ngoctrng/bookz/internal/worker"
	"github.com/ngoctrng/bookz/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Pass, cfg.DB.Name, cfg.DB.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	w := worker.New(cfg, db)

	slog.Info("Starting worker server...")
	if err := w.Run(); err != nil {
		log.Fatalf("could not run worker: %v", err)
	}
}
