package main

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	Rol      string `json:"rol"`
	Token    string `json:"token"`
}

func (o *User) login(db *sql.DB, username string, password string) error {
	statement := fmt.Sprintf("SELECT username, password FROM users WHERE username = '%s' AND password = '%x'", username, password)
	fmt.Println(statement)
	return db.QueryRow(statement).Scan(&o.Username, &o.Password)
}

func (o *User) getUserByUsername(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", o.Username)
	return db.QueryRow(statement).Scan(&o.ID, &o.Username, &o.Password , &o.Name, &o.LastName, &o.Rol, &o.Token)
}

func (o *User) getUserByID(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT * FROM users WHERE id = %d", o.ID)
	return db.QueryRow(statement).Scan(&o.ID, &o.Username, &o.Password , &o.Name, &o.LastName, &o.Rol, &o.Token)
}

func (o *User) insertUser(db *sql.DB) error {
	statement := fmt.Sprintf(
		"INSERT INTO users (username, password, name, lastname, rol) "+
			"VALUES('%s', '%x', '%s', '%s', '%s')",
		o.Username, o.Password, o.Name, o.LastName, o.Rol)

	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

func (o *User) updateUser(db *sql.DB) error {
	statement := fmt.Sprintf(
		"UPDATE users "+
			"SET username=%d, password='%s', name='%s', lastname='%s', rol='%s' "+
			"WHERE id=%d",
		o.ID, o.Username, o.Name, o.LastName, o.Rol, o.ID)
	_, err := db.Exec(statement)
	return err
}

func (o *User) updateUserTokenByUsername(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE users SET token = '%s' WHERE username = '%s'", o.Token, o.Username)
	_, err := db.Exec(statement)
	return err
}

func (o *User) deleteUser(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM users WHERE id=%d", o.ID)
	_, err := db.Exec(statement)
	return err
}

func getUsers(db *sql.DB) ([]User, error) {
	var list []User
	rows, err := db.Query("SELECT * from users")

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var o User
		if err := rows.Scan(&o.ID, &o.Username, &o.Password , &o.Name, &o.LastName, &o.Rol, &o.Token); err != nil {
			return nil, err
		}
		list = append(list, o)
	}
	return list, nil
}
