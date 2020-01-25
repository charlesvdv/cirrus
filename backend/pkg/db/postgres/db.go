package postgres

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type Config struct {
	DBName   string
	User     string
	Password string
	Host     string
	Port     uint16
	sslMode  string
}

func Open(conf Config) (*sql.DB, error) {
	connStr := formatConnectionString(conf)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return db, err
	}

	err = setupSchema(db)
	if err != nil {
		return db, err
	}

	return db, nil
}

func formatConnectionString(conf Config) string {
	// TODO: sane default
	return formatConnectionStringKV("dbname", conf.DBName) +
		formatConnectionStringKV("user", conf.User) +
		formatConnectionStringKV("password", conf.Password) +
		formatConnectionStringKV("host", conf.Host) +
		formatConnectionStringKV("port", strconv.FormatUint(uint64(conf.Port), 10)) +
		formatConnectionStringKV("sslmode", conf.sslMode)
}

func formatConnectionStringKV(key, value string) string {
	value = strings.ReplaceAll(value, "'", "\\'")
	return fmt.Sprintf("%s='%s' ", key, value)
}
