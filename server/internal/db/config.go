package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rawbil/movie-stream-app/internal/server"
	"github.com/rawbil/movie-stream-app/internal/utils"
)

func DBConfigFunc(config *server.DBConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.Dsn)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	utils.Log.Info("DB CONNECTED SUCCESSFULLY!!")

	return db, nil
}
