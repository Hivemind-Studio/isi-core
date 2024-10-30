package db

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
)

func InitMigration(db *sql.DB) {
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	if err := goose.SetDialect("mysql"); err != nil {
		panic(err)
	}

	err := goose.Up(db, "db/migrations", goose.WithAllowMissing())
	if err != nil {
		panic(err)
	}

	fmt.Println("Migrations applied successfully.")
}
