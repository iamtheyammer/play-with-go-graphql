package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var db = func() *sql.DB {
	db, err := sql.Open("sqlite3", "links.db")
	if err != nil {
		panic(fmt.Errorf("error opening SQLite3 database: %w", err))
	}

	return db
}()

// PrepareDB prepares the database for use. Call in main().
func PrepareDB() error {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT UNIQUE NOT NULL, password TEXT NOT NULL, token TEXT, refresh_token TEXT)")
	if err != nil {
		return fmt.Errorf("error preparing create users table sql: %w", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}

	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS links (id INTEGER PRIMARY KEY AUTOINCREMENT, creator_user_id INTEGER NOT NULL REFERENCES users(id), title TEXT NOT NULL, address TEXT NOT NULL)")
	if err != nil {
		return fmt.Errorf("error preparing create links table sql: %w", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("error creating links table sql: %w", err)
	}

	return nil
}


