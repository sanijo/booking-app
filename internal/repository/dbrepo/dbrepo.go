package dbrepo

import (
	"database/sql"

	"github.com/sanijo/rent-app/internal/config"
	"github.com/sanijo/rent-app/internal/repository"
)


type postgresDbRepo struct {
    App *config.AppConfig
    DB  *sql.DB
}

type testDBRepo struct {
    App *config.AppConfig
    DB  *sql.DB
}

// NewPostgresRepo creates a new repository
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
    return &postgresDbRepo{
        App: a,
        DB: conn,
    }
}

// NewTestingRepo creates a new repository for testing
func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
    return &testDBRepo{
        App: a,
    }
}
