package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func NewDb() (*sql.DB, error) {
	config := newConfig().database
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.User, config.Pass, config.Host, config.Port, config.DbName)
	log.Println(connStr)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
