package sso

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
)

// Names of tables in the database that are used by this package.
// This package does not create the tables; see mysql/schema.sql.
const (
	memberTable       = "fsso_members"
	emailAuthTable    = "fsso_auth_email"
	googleAuthTable   = "fsso_auth_goo"
	facebookAuthTable = "fsso_auth_fb"
	activeTable       = "fsso_active"
	refreshTable      = "fsso_refresh"
	emailVerifyTable  = "fsso_email_verify"
)

var db *sql.DB

// InitDB sets the db handle for use in all subsequent db operations.
// Must be called once.
// The dsn string should look like "user:pass@host/dbname".
func InitDB(dsn string) {
	d, err := sql.Open("mysql", dsn)
	if err == nil {
		err = d.Ping()
	}
	if err != nil {
		// Give up completely if our db connection fails.
		log.Fatal(err)
	}
	db = d
}

// isDuplicate tests err to see if a db error is a unique constraint violation
func isDuplicate(err error) bool {
	if err0, ok := err.(*mysql.MySQLError); ok {
		// 1062 is mysql for unique contstraint violation
		return err0.Number == 1062
	}
	return false
}
