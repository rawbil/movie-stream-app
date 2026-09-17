package main

import (
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/rawbil/movie-stream-app/internal/db"
	"github.com/rawbil/movie-stream-app/internal/server"
	"github.com/rawbil/movie-stream-app/internal/utils"
)

func main() {

	if err := godotenv.Load(); err != nil {
		utils.Log.Error("ERROR LOADING ENV VARIABLES")
		return
	}

	cfg := mysql.Config{
		User:                 utils.DBConfigFunc().Username,
		Passwd:               utils.DBConfigFunc().Passwd,
		Addr:                 utils.DBConfigFunc().Addr,
		DBName:               utils.DBConfigFunc().Name,
		ParseTime:            utils.DBConfigFunc().ParseTime,
		AllowNativePasswords: true,
	}
	
	dbConfig := &server.DBConfig{
		Dsn:  cfg.FormatDSN(),
		Addr: ":8080",
	}

	db, err := db.DBConfigFunc(dbConfig)
	if err != nil {
		utils.Log.Error("DB ERROR", "error", err)
		os.Exit(1)
	}

	api := &server.Api{
		DBConfig: *dbConfig,
		DB:       db,
	}

	if err := api.Run(api.Mount()); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
