package mysqlconn

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"log/slog"
	"time"
)

func Init(host string, port string, uname string, pwd string, dbname string) *sqlx.DB {

	defer func() {
		if r := recover(); r != nil {
			slog.Error("Errors")
			fmt.Println("Recovered from panic:", r)
		}
	}()

	dsnString := fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true", uname, pwd, host, port, dbname)

	db, err := sqlx.Connect("mysql", dsnString)

	if err != nil {
		log.Fatalln(err)
	}

	db.SetMaxOpenConns(300)
	db.SetMaxIdleConns(300)
	db.SetConnMaxLifetime(3 * time.Minute)

	return db
}
