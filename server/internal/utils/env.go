package utils

import (
	"fmt"
	"os"
)

type DBConfig struct {
	Name      string
	Passwd    string
	Addr      string
	Username  string
	ParseTime bool
}

type ServerConfig struct {
	ServerAddr string
}

func DBConfigFunc() *DBConfig {
	return &DBConfig{
		Name:      GetEnv("DB_NAME", "stream_app"),
		Passwd:    GetEnv("DB_PASSWD", ""),
		Addr:      fmt.Sprintf("%s:%s", GetEnv("DB_HOST", "localhost"), GetEnv("DB_PORT", "3306")),
		Username:  GetEnv("DB_USERNAME", "root"),
		ParseTime: GetEnv("DB_PARSE_TIME", "true") == "true",
	}
}

func ServerConfigFunc() *ServerConfig {
	return &ServerConfig{
		ServerAddr: GetEnv("SERVER_ADDR", ":8080"),
	}
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
