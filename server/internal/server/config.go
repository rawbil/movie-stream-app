package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	repository "github.com/rawbil/movie-stream-app/internal/adapters/sqlc"
	"github.com/rawbil/movie-stream-app/internal/auth"
	"github.com/rawbil/movie-stream-app/internal/auth/authutils"
	"github.com/rawbil/movie-stream-app/internal/movies"
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

	repo := repository.New(api.DB)
	movieService := movies.NewService(*repo, api.DB)
	movieHandler := movies.NewHandler(movieService)

	authService := auth.NewService(*repo, api.DB)
	authHandler := auth.NewHandler(authService)

	// groups
	api_v1 := r.Group("/api/v1")
	movies := api_v1.Group("/movies")
	movie_genres := movies.Group("/genres")
	auth := api_v1.Group("/auth")

	api_v1.GET("/health", func(c *gin.Context) {
		if err := api.DB.Ping(); err != nil {
			utils.Log.Error("ERROR CONNECTING TO DB", "error", err)
			return
		}
		utils.Log.Info("server and db ok")
		c.JSON(http.StatusOK, gin.H{
			"message": "Server and DB OK",
		})
	})

	//! /api/v1/auth
	auth.POST("/register", authHandler.RegisterUser)
	auth.POST("/login", authHandler.LoginUser)
	auth.POST("/create-role", authutils.AuthMiddleware(*repo), authHandler.CreateRole)

	//! /api/v1/movies/genres
	movie_genres.POST("/add", authutils.AuthMiddleware(*repo), movieHandler.CreateGenre)
	movie_genres.PATCH("/update",authutils.AuthMiddleware(*repo), movieHandler.UpdateGenre)
	movie_genres.GET("/list", movieHandler.ListGenres)
	movies.GET("/all", movieHandler.ListMovies)

	//! /api/v1/movies
	movies.POST("/create",authutils.AuthMiddleware(*repo), movieHandler.CreateMovie)
	movies.GET("/one", movieHandler.GetMovie)

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
