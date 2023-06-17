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

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
    return &postgresDbRepo{
        App: a,
        DB: conn,
    }
}
