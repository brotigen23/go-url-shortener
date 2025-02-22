package database

import "database/sql"

// Проверяет соединение с базой данных
func CheckPostgresConnection(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	return db.Ping()
}
