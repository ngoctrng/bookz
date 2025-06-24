package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/ngoctrng/bookz/pkg/config"
	"github.com/ngoctrng/bookz/pkg/migration"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("cannot load config: %v\n", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Pass, cfg.DB.Name, cfg.DB.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	conn, err := db.DB()
	if err != nil {
		log.Fatalf("cannot get database connection: %v", err)
	}

	total, err := migration.Run(conn)
	if err != nil {
		log.Fatalf("cannot execute migration: %v\n", err)
	}

	slog.Info(fmt.Sprintf("applied %d migrations\n", total))
}
