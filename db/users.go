package db

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"name"`
	Password string `json:"password"`
}

func InsertUser(username string, realPassword string) (int, error) {
	hashedPw, err := hashPassword(realPassword)
	if err != nil {
		return 0, fmt.Errorf("error hashing password when inserting user: %w", err)
	}

	stmt, err := db.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("error preparing insert user sql: %w", err)
	}

	result, err := stmt.Exec(username, hashedPw)
	if err != nil {
		return 0, fmt.Errorf("error inserting user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert id: %w", err)
	}

	return int(id), err
}

type GetUserRequest struct {
	ID       int
	Username string
}

func getUser(req *GetUserRequest) (*User, error) {
	query := "SELECT id, username, password FROM users WHERE "
	var arg interface{}

	if req.ID > 0 {
		query += "id = ?"
		arg = req.ID
	} else {
		query += "username = ?"
		arg = req.Username
	}

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error preparing get user by id sql: %w", err)
	}

	row := stmt.QueryRow(arg)

	var u User

	err = row.Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("error scanning get user by id row: %w", err)
	}

	return &u, nil
}

func GetUserByID(userId int) (*User, error) {
	return getUser(&GetUserRequest{ID: userId})
}

func GetUserByUsername(username string) (*User, error) {
	return getUser(&GetUserRequest{Username: username})
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", fmt.Errorf("error generating hash from password: %w", err)
	}

	return string(bytes), nil
}

func ComparePasswordToHash(hashPassword string, realPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(realPassword))
	return err == nil
}
