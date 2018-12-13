package storage

import (
	"backend/conf"
	"backend/utils"
	"github.com/jmoiron/sqlx"
)

var dbPrefix = "backend"

// NewDBConn new mysql connection
func NewDBConn(dbName string) *sqlx.DB {
	if dbName == "" {
		dbName = "master"
	}

	conn, err := sqlx.Open("mysql", conf.Conf.DB[dbPrefix][dbName])
	conn.SetMaxOpenConns(2000)
	conn.SetMaxIdleConns(1000)
	if err != nil {
		utils.ErrToPanic(err)
	}

	if err := conn.Ping(); err != nil {
		utils.ErrToPanic(err)
	}

	return conn
}
