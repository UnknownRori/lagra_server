package models

import (
	"encoding/hex"

	"github.com/UnknownRori/lagra_server/src"

	"github.com/google/uuid"
)

type FindUsersByUuid struct {
	Uuid string `param:"uuid" query:"uuid" form:"uuid" json:"uuid" xml:"uuid" validate:"required"`
}

type NewUser struct {
	Username string `json:"username" xml:"username" validate:"required"`
	Password string `json:"password" xml:"password" validate:"required"`
	Email    string `json:"email" xml:"email" validate:"required,email"`
	Role     string `json:"role" xml:"role" validate:"required"`
}

type LoginUser struct {
	Username string `json:"username" xml:"username" validate:"required"`
	Password string `json:"password" xml:"password" validate:"required"`
}

type User struct {
	Uuid     string `json:"uuid" xml:"uuid"`
	Username string `json:"username" xml:"username"`
	Password string `json:"password" xml:"password"`
	Email    string `json:"email" xml:"email"`
	Role     string `json:"role" xml:"role"`
}

type ReturnUser struct {
	Uuid     string `json:"uuid" xml:"uuid"`
	Username string `json:"username" xml:"username"`
	Email    string `json:"email" xml:"email"`
	Role     string `json:"role" xml:"role"`
}

type DisplayUser struct {
	Uuid     string `json:"uuid" xml:"uuid"`
	Username string `json:"username" xml:"username"`
}

func CreateUser(db *src.DB, user NewUser) (uuid.UUID, error) {
	uuid := uuid.New()
	password := src.CreateHash([]byte(user.Password))
	stmt, err := db.Prepare("INSERT INTO users (uuid, username, email, password, role) VALUES (?, ?, ?, ?, ?)")
	defer stmt.Close()

	if err != nil {
		return uuid, err
	}

	_, err = stmt.Exec(uuid, user.Username, user.Email, hex.EncodeToString(password), user.Role)

	if err != nil {
		return uuid, err
	}

	return uuid, nil
}

func FetchUserByUuid(db *src.DB, uuid string) (User, error) {
	var user User
	stmt, err := db.Prepare("SELECT * FROM users WHERE uuid = ?")
	if err != nil {
		return user, err
	}

	query := stmt.QueryRow(uuid)
	err = query.Scan(&user.Uuid, &user.Username, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return user, err
	}

	return user, nil
}

func FetchUserByUsername(db *src.DB, username string) (User, error) {
	var user User
	stmt, err := db.Prepare("SELECT * FROM users WHERE username = ?")
	if err != nil {
		return user, err
	}

	query := stmt.QueryRow(username)
	err = query.Scan(&user.Uuid, &user.Username, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return user, err
	}

	return user, nil
}

func FetchUserByUsernameOrUuid(db *src.DB, usernameOrUuid string) (User, error) {
	var user User
	stmt, err := db.Prepare("SELECT * FROM users WHERE username = ? OR uuid = ?")
	if err != nil {
		return user, err
	}

	query := stmt.QueryRow(usernameOrUuid, usernameOrUuid)
	err = query.Scan(&user.Uuid, &user.Username, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return user, err
	}

	return user, nil
}
