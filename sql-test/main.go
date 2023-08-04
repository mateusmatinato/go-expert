package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func main() {
	db, err := sql.Open("mysql", "mateus:mateus@tcp(localhost:3306)/goexpert")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	args := os.Args[1:]
	if len(args) == 0 {
		printHelp()
		return
	}

	cmd := args[0]
	params := args[1:]
	switch cmd {
	case "INSERT":
		if err := validateUserParams(params); err != nil {
			panic(err.Error())
		}

		user, err := NewUser(params[0], params[1])
		if err != nil {
			panic(err.Error())
		}

		if err = insertUser(db, user); err != nil {
			panic(err.Error())
		}
		fmt.Println("User inserted successfully")
		fmt.Printf("%v\n", *user)
	case "SELECT":
		var id *string
		if len(params) > 0 {
			id = &params[0]
		}

		users, err := selectUser(db, id)
		if err != nil {
			panic(err.Error())
		}
		if len(users) == 0 {
			fmt.Println("No user found")
			return
		}

		for _, u := range users {
			fmt.Printf("%v\n", u)
		}
	case "DELETE":
		if err := validateDeleteParams(params); err != nil {
			panic(err.Error())
		}

		user, err := selectUser(db, &params[0])
		if err != nil {
			panic(err.Error())
		}
		if len(user) == 0 {
			fmt.Println("User not found")
			return
		}

		if err = deleteUser(db, user[0].Id); err != nil {
			panic(err.Error())
		}

		fmt.Println("User deleted successfully")
		fmt.Printf("%v\n", user[0])
	default:
		printHelp()
	}
}

type User struct {
	Id   string
	Name string
	Age  uint64
}

func NewUser(name string, ageStr string) (*User, error) {
	age, err := strconv.ParseUint(ageStr, 10, 64)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &User{
		Id:   id.String(),
		Name: name,
		Age:  age,
	}, nil
}

func insertUser(db *sql.DB, user *User) error {
	_, err := db.Exec("INSERT INTO users (id, name, age) VALUES (?, ?, ?)", user.Id, user.Name, user.Age)
	if err != nil {
		return err
	}
	return nil
}

func deleteUser(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func selectUser(db *sql.DB, id *string) ([]User, error) {
	var where string
	var args []any
	if id != nil {
		where = "WHERE id = ?"
		args = append(args, *id)
	}

	rows, err := db.Query("SELECT * FROM users "+where, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Name, &user.Age); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func validateDeleteParams(params []string) error {
	if len(params) != 1 {
		return fmt.Errorf("expected 1 param, got %d", len(params))
	}
	return nil
}

func validateUserParams(params []string) error {
	if len(params) != 2 {
		return fmt.Errorf("expected 2 params, got %d", len(params))
	}
	return nil
}

func printHelp() {
	fmt.Println("\nAvailable commands: INSERT/SELECT/DELETE")
	fmt.Println("INSERT: Expects two required params: name and age")
	fmt.Println("Example: INSERT Mateus 25")
	fmt.Println()
	fmt.Println("SELECT: Expects one optional param: id (default is select all)")
	fmt.Println("Example: SELECT 12345-abcde-67890-fghij")
	fmt.Println()
	fmt.Println("DELETE: Expects one required param: id")
	fmt.Println("Example: DELETE 12345-abcde-67890-fghij")
}
