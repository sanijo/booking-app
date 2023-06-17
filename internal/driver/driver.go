package driver

import (
	"database/sql"
	"time"

    _ "github.com/jackc/pgx/v5"
    _ "github.com/jackc/pgx/v5/stdlib"
    _ "github.com/jackc/pgx/v5/pgconn"
)

// DB is a wrapper around sql.DB and holds the database connections.
// It makes it easier to switch to a different database later on.
type DB struct {
    SQL *sql.DB
}

// dbConn is package level variable 
var dbConn = &DB{}

// maxOpenDbConn is the maximum number of open connections to the database.
const maxOpenDbConn = 10

// maxIdleDbConn is the maximum number of idle connections to the database.
const maxIdleDbConn = 5

// maxDbLifetime is the maximum amount of time a connection may be reused.
const maxDbLifetime = 5 * time.Minute

// testDB tests the database connection using Ping.
func testDB(db *sql.DB) error {
    if err := db.Ping(); err != nil {
        return err
    }

    return nil
}

// NewDatabase opens a new database connection to the given DSN.
func NewDatabase(dsn string) (*sql.DB, error) {
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        return nil, err
    }

    if err = testDB(db); err != nil {
        return nil, err
    }

    return db, nil
}

// ConnectSQL creates database pool for postgresql.
func ConnectSQL(dsn string) (*DB, error) {
    db, err := NewDatabase(dsn)
    if err != nil {
        panic(err)
    }

    db.SetMaxOpenConns(maxOpenDbConn)
    db.SetMaxIdleConns(maxIdleDbConn)
    db.SetConnMaxLifetime(maxDbLifetime)

    dbConn.SQL = db

    if err = testDB(db); err != nil {
        return nil, err
    }

    return dbConn, nil
}

