package src

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	conn *sql.DB
}

func CreateConn() (DB, error) {
	var db DB

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)
	log.Infof("Connecting to database : %s", os.Getenv("DB_DATABASE"))

	conn, err := sql.Open(os.Getenv("DB_CONNECTION"),
		connStr)

	conn.SetConnMaxLifetime(time.Minute * 3)
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10)

	if err != nil {
		return db, err
	}

	log.Infof("Connected to database")
	db.conn = conn
	return db, nil
}

func (d *DB) Prepare(query string) (*sql.Stmt, error) {
	return d.conn.Prepare(query)
}

func (d *DB) Close() error {
	return d.conn.Close()
}
