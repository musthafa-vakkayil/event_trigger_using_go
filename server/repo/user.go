package repo

import (
	"database/sql"
	"event_trigger/model"
	"fmt"
	"log"
)

func DeleteUser(db *sql.DB, user_id string) error {
	_, err := db.Exec("DELETE FROM public.users WHERE id = $1", user_id)
	return err
}

func GetUserByID(db *sql.DB, user_id string) (model.User, error) {
	var user model.User
	if err := db.QueryRow("SELECT id, username, email FROM public.users WHERE id=$1", user_id).Scan(
		&user.UserId, &user.Username, &user.Email,
	); err != nil {
		{
			return model.User{}, err
		}
	}

	return user, nil
}

func ListUsers(db *sql.DB) ([]model.User, error) {
	rows, err := db.Query("SELECT id, username, email FROM public.users")
	if err != nil {
		return []model.User{}, nil
	}

	var users []model.User

	for rows.Next() {
		var msg model.User
		if err := rows.Scan(
			&msg.UserId, &msg.Username, &msg.Email,
		); err != nil {
			return nil, err
		}

		users = append(users, msg)
	}

	return users, nil
}

func CreateUser(db *sql.DB, usr model.User) (string, error) {
	stmt := "INSERT INTO public.users(id, username, email, password_hash) VALUES ($1, $2, $3, $4) RETURNING id"

	var id string

	if err := db.QueryRow(stmt, usr.UserId, usr.Username, usr.Email, usr.PasswordHash).Scan(&id); err != nil {
		log.Print(err)
		return "", err
	}

	fmt.Println("New user successfully created, id:", id)
	return id, nil
}

func UpdateUser(db *sql.DB, usr model.User) error {
	stmt, err := db.Prepare("UPDATE public.users SET username = $1, email= $2, password_hash = $3 WHERE id = $4")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, qErr := stmt.Exec(usr.Username, usr.Email, usr.PasswordHash, usr.UserId)
	if qErr != nil {
		return err
	}
	return nil
}
