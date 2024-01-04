package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func NewDb() (*sql.DB, error) {
	c, err := tryGetConfigDB()
	if err != nil {
		logrus.Error(err.Error())
		return nil, fmt.Errorf("something wrong in db environment: %w", err)
	}

	config := c.database

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.User, config.Pass, config.Host, config.Port, config.DbName)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
