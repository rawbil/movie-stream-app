package main

import (
	"log"
	"os"

	"github.com/rawbil/movie-stream-app/internal/server"
)

func main() {
	dbConfig := &server.DBConfig{
		Dsn:  "",
		Addr: ":8080",
	}

	api := &server.Api{
		DBConfig: *dbConfig,
		DB:       nil,
	}

	if err := api.Run(api.Mount()); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
