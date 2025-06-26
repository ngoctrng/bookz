package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/hibiken/asynq"
	"github.com/ngoctrng/bookz/internal/httpserver"
	"github.com/ngoctrng/bookz/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	sentrygo "github.com/getsentry/sentry-go"
	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = sentrygo.Init(sentrygo.ClientOptions{
		Dsn:              cfg.SentryDSN,
		Environment:      cfg.AppEnv,
		AttachStacktrace: true,
	})
	if err != nil {
		log.Fatalf("cannot init sentry: %v", err)
	}
	defer sentrygo.Flush(5 * time.Second)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Pass, cfg.DB.Name, cfg.DB.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	server := httpserver.New(cfg, db, client)
	log.Fatal(server.Start())
}
