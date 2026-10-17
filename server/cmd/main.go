package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	repository "github.com/rawbil/movie-stream-app/internal/adapters/sqlc"
	"github.com/rawbil/movie-stream-app/internal/db"
	"github.com/rawbil/movie-stream-app/internal/seed"
	"github.com/rawbil/movie-stream-app/internal/server"
	"github.com/rawbil/movie-stream-app/internal/utils"
)

func registerRDSCert() error {
	caCert, err := os.ReadFile("/etc/ssl/certs/rds-global-bundle.pem")
	if err != nil {
		return err
	}

	rootCAs, err := x509.SystemCertPool()
	if err != nil || rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
		return fmt.Errorf("failed to append RDS CA certificates")
	}

	return mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs:    rootCAs,
		MinVersion: tls.VersionTLS12,
	})
}

func main() {

	if err := godotenv.Load(); err != nil {
		utils.Log.Error("ERROR LOADING ENV VARIABLES")
		return
	}

	if err := registerRDSCert(); err != nil {
		log.Fatal(err)
	}

	cfg := mysql.Config{
		User:                 utils.DBConfigFunc().Username,
		Passwd:               utils.DBConfigFunc().Passwd,
		Addr:                 utils.DBConfigFunc().Addr,
		DBName:               utils.DBConfigFunc().Name,
		ParseTime:            utils.DBConfigFunc().ParseTime,
		AllowNativePasswords: true,
		TLSConfig: "custom",
	}

	dbConfig := &server.DBConfig{
		Dsn:  cfg.FormatDSN(),
		Addr: utils.ServerConfigFunc().ServerAddr,
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

	//~ Seed Data
	if err := seed.SeedPermissions(db); err != nil {
		utils.Log.Error("Error seeding permissions", "error", err)
		return
	}

	if err := seed.SeedRoles(db); err != nil {
		utils.Log.Error("Error seeding roles", "error", err)
		return
	}

	if err := seed.SeedRolePermissions(context.Background(), *repository.New(db), db); err != nil {
		utils.Log.Error("Error seeding role permissions", "error", err)
		return
	}

	if err := seed.SeedMovieRankings(db); err != nil {
		utils.Log.Error("Error seeding rankings", "error", err)
		return
	}

	if err := api.Run(api.Mount()); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
