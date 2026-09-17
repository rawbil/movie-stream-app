package server

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Api struct {
	DB       *sql.DB
	DBConfig DBConfig
}

type DBConfig struct {
	Dsn  string
	Addr string
}

func (api *Api) Mount() http.Handler {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server & DB OK",
		})
	})

	return r
}

func (api *Api) Run(m http.Handler) error {
	srv := http.Server{
		Addr:         api.DBConfig.Addr,
		Handler:      m,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server is listening on port: %s", api.DBConfig.Addr)
	return srv.ListenAndServe()
}
