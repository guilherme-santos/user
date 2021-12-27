package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	migratemysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewConnection(host string, port int, user, password, database, migrationDir string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true", user, password, host, port, database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(1 * time.Hour)
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(10)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	cfg := &migratemysql.Config{
		DatabaseName: database,
	}
	driver, err := migratemysql.WithInstance(db, cfg)
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+migrationDir, "mysql", driver)
	if err != nil {
		return nil, err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, err
	}
	return db, nil
}

func IsDuplicateError(err error, field string) bool {
	var merr *mysql.MySQLError
	if errors.As(err, &merr) {
		if merr.Number == 1062 && strings.Contains(merr.Message, fmt.Sprintf("'%s'", field)) {
			return true
		}
	}
	return false
}
