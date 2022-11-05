package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB Holds The Database Connection
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDConn = 10
const maxIdleDConn = 5
const maxDbLifetime = 5 * time.Minute

// ConnectSQL Creates Database Pool For Postgres
func ConnectSQL(dsn string) (*DB, error) {
	d, err := newDatabase(dsn)
	if err != nil {
		panic(err)
	}

	d.SetMaxOpenConns(maxOpenDConn)
	d.SetMaxIdleConns(maxIdleDConn)
	d.SetConnMaxLifetime(maxDbLifetime)

	dbConn.SQL = d
	err = testDB(d)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

// testDB Tries to Pin the Database
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

// newDatabase Creates a new Database for the application
func newDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
