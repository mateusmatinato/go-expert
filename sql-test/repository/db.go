package repository

import (
	"database/sql"

	userDomain "github.com/mateusmatinato/go-expert/sql-test/domain/user"

	_ "github.com/go-sql-driver/mysql"
)

type UserRepository interface {
	InsertUser(*userDomain.User) error
	DeleteUser(string) error
	SelectUser(*string) ([]userDomain.User, error)
}

type userRepositoryHandler struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) UserRepository {
	createTableUsers(db)

	return &userRepositoryHandler{
		db: db,
	}
}

func createTableUsers(db *sql.DB) {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(36) NOT NULL,
			name VARCHAR(255) NOT NULL,
			age INT NOT NULL,
			PRIMARY KEY (id)
		)
	`); err != nil {
		panic(err.Error())
	}
}

func (r *userRepositoryHandler) InsertUser(user *userDomain.User) error {
	_, err := r.db.Exec("INSERT INTO users (id, name, age) VALUES (?, ?, ?)", user.Id, user.Name, user.Age)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepositoryHandler) DeleteUser(id string) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepositoryHandler) SelectUser(id *string) ([]userDomain.User, error) {
	var where string
	var args []any
	if id != nil {
		where = "WHERE id = ?"
		args = append(args, *id)
	}

	rows, err := r.db.Query("SELECT * FROM users "+where, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []userDomain.User{}
	for rows.Next() {
		var user userDomain.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Age); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
