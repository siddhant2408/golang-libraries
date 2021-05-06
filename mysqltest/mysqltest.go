// Package mysqltest provides test helpers for MySQL.
//
// The DSN can be overridden with the MYSQLTEST_DSN environment variable.
//
// If MySQL is not available, the test is skipped.
// It can be controlled with the MYSQLTEST_UNAVAILABLE_SKIP environment variable.
package mysqltest

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/mysqltracing"
	"github.com/siddhant2408/golang-libraries/sqltest"
)

var helper = &sqltest.Helper{
	Name:                  "MySQL",
	DSNDefault:            "root@tcp(localhost:3306)/",
	DSNEnvVar:             "MYSQLTEST_DSN",
	UnavailableSkipEnvVar: "MYSQLTEST_UNAVAILABLE_SKIP",
	Open:                  open,
	CreateQuery:           createQuery,
	DropQuery:             dropQuery,
}

// GetDSN calls the helper.
var GetDSN = helper.GetDSN

// GetDatabase calls the helper.
var GetDatabase = helper.GetDatabase

// CheckAvailable calls the helper.
var CheckAvailable = helper.CheckAvailable

func open(ctx context.Context, dsn string, dbName string) (*sql.DB, error) {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "parse DSN")
	}
	if dbName != "" {
		cfg.DBName = dbName
	}
	cr, err := mysqltracing.NewConnector(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "new connector")
	}
	db := sql.OpenDB(cr)
	return db, nil
}

func createQuery(dbName string) string {
	return fmt.Sprintf("CREATE DATABASE %s;", dbName)
}

func dropQuery(dbName string) string {
	return fmt.Sprintf("DROP DATABASE %s;", dbName)
}
