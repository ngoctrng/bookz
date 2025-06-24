package migration

import (
	"database/sql"
	"embed"
	"errors"

	migrate "github.com/rubenv/sql-migrate"
)

//go:embed migrations/*
var dbMigrations embed.FS

func Run(db *sql.DB) (int, error) {
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigrations,
		Root:       "migrations",
	}

	total, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return 0, errors.Join(err, errors.New("cannot execute migration"))
	}

	return total, nil
}
