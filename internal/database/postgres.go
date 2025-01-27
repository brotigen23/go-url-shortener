package database

import "database/sql"

func CheckPostgresConnection(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	return db.Ping()
}
