package db

import (
	"database/sql"
	"log"
	"net/url"

	"RuAm/internal/pkg/config"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var DB *sql.DB

func InitDB(cfg config.Config) {
	conn, err := url.Parse(cfg.ServiceURL)
	if err != nil {
		log.Fatalf("could not parse service URL: %v", err)
	}
	conn.RawQuery = "sslmode=verify-ca;sslrootcert=ca.pem"

	DB, err = sql.Open("postgres", conn.String())
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("could not ping db: %v", err)
	}

	goose.SetDialect("postgres")
	err = goose.Up(DB, "./migrations")
	if err != nil {
		log.Fatalf("could not run migrations: %v", err)
	}

	log.Println("Connected to DB")
}
