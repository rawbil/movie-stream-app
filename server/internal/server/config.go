package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rawbil/movie-stream-app/internal/utils"
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
		if err := api.DB.Ping(); err != nil {
			utils.Log.Error("ERROR CONNECTING TO DB", "error", err)
			return
		}
		utils.Log.Info("server and db ok")
		c.JSON(http.StatusOK, gin.H{
			"message": "Server and DB OK",
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

	utils.Log.Info(fmt.Sprintf("SERVER IS LISTENING ON PORT: %s", api.DBConfig.Addr))
	return srv.ListenAndServe()
}
