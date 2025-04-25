package config

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type pgsql struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string
}

func (p *pgsql) openPostgres() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.DBHost, p.DBPort, p.DBUser, p.DBPassword, p.DBName, p.DBSSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func newPostgres() *pgsql {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	return &pgsql{
		DBSSLMode:  os.Getenv("DB_SSL_MODE"),
		DBHost:     os.Getenv("DB_HOST"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBPort:     os.Getenv("DB_PORT"),
	}
}
