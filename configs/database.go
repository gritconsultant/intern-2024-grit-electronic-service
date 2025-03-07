package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

var (
	dbMap = make(map[string]*bun.DB)
	db    *bun.DB
)

type Option struct {
	Dsn      string
	Host     string
	Port     int64
	Database string
	Username string
	Password string
	TimeZone string
}

func Database() *bun.DB {

	godotenv.Load()

	db, err := defaultConfig()
	if err != nil {
		panic(err)
	}

	db.AddQueryHook(bundebug.NewQueryHook())

	return db

}

func Open(ctx context.Context) error {
	var err error
	for _, db := range dbMap {
		if errClose := db.Ping(); errClose != nil {
			err = errors.Join(err, errClose)
		}
	}
	return err
}

func Close(ctx context.Context) error {
	var err error
	for _, db := range dbMap {
		if errClose := db.Close(); errClose != nil {
			err = errors.Join(err, errClose)
		}
	}
	return err
}

func defaultConfig() (*bun.DB, error) {
	port, _ := strconv.ParseInt(os.Getenv("DB_PORT"), 10, 64)

	op := Option{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		Database: os.Getenv("DB_DATABASE"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		TimeZone: "Asia/Bangkok",
	}

	op.Dsn = fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable TimeZone=%s", op.Host, op.Port, op.Database, op.Username, op.Password, op.TimeZone)

	config, err := pgx.ParseConfig(op.Dsn)
	if err != nil {
		panic(err)
	}
	sqldb := stdlib.OpenDB(*config)
	db := bun.NewDB(sqldb, pgdialect.New())

	return db, db.Ping()
}
