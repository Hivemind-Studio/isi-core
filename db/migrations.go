package db

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
)

func InitMigration(db *sql.DB) {
	if err := db.Ping(); err != nil {
		db.Close()
		panic(err)
	}

	if err := goose.SetDialect("mysql"); err != nil {
		defer db.Close()
		panic(err)
	}

	err := goose.Up(db, "db/migrations", goose.WithAllowMissing())
	if err != nil {
		db.Close()
		panic(err)
	}

	fmt.Println("Migrations applied successfully.")
}
