package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/mateusmatinato/go-expert/sql-test/domain/user"
	"github.com/mateusmatinato/go-expert/sql-test/repository"
	"github.com/mateusmatinato/go-expert/sql-test/validator"
)

func main() {
	db, err := sql.Open("mysql", "mateus:mateus@tcp(localhost:3306)/goexpert")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	userRepo := repository.NewRepository(db)

	args := os.Args[1:]
	if len(args) == 0 {
		printHelp()
		return
	}

	cmd := args[0]
	params := args[1:]

	switch cmd {
	case "INSERT":
		if err := validator.ValidateUserParams(params); err != nil {
			panic(err.Error())
		}

		user, err := user.NewUser(params[0], params[1])
		if err != nil {
			panic(err.Error())
		}

		if err = userRepo.InsertUser(user); err != nil {
			panic(err.Error())
		}
		fmt.Println("User inserted successfully")
		fmt.Printf("%v\n", *user)
	case "SELECT":
		var id *string
		if len(params) > 0 {
			id = &params[0]
		}

		users, err := userRepo.SelectUser(id)
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
		if err := validator.ValidateDeleteParams(params); err != nil {
			panic(err.Error())
		}

		user, err := userRepo.SelectUser(&params[0])
		if err != nil {
			panic(err.Error())
		}
		if len(user) == 0 {
			fmt.Println("User not found")
			return
		}

		if err = userRepo.DeleteUser(user[0].Id); err != nil {
			panic(err.Error())
		}

		fmt.Println("User deleted successfully")
		fmt.Printf("%v\n", user[0])
	default:
		printHelp()
	}
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
