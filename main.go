package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	fmt.Println("*** Todo application ***")

	command := flag.String("command", "no command", "Command to execute")
	flag.Parse()

	for {
		runCommand(*command)
		fmt.Printf("Enter a command %v: ", commandList)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		*command = scanner.Text()
	}

	fmt.Printf("UserStorage: %+v\n", userStorage)

}

type User struct {
	ID       int
	Email    string
	Password string
}

var userStorage []User
var commandList = []string{
	"create-task",
	"create-category",
	"register-user",
	"user-login",
	"exit",
}

func runCommand(command string) {
	switch command {
	case "create-task":
		createTask()
	case "create-category":
		createCategory()
	case "register-user":
		registerUser()
	case "user-login":
		userLogin()
	case "exit":
		os.Exit(0)
	default:
		fmt.Println("Command not recognized")

	}
}

func createTask() {
	scanner := bufio.NewScanner(os.Stdin)
	var name, dueDate, category string

	fmt.Printf("Enter task title: ")
	scanner.Scan()
	name = scanner.Text()

	fmt.Printf("Enter task dueDate: ")
	scanner.Scan()
	dueDate = scanner.Text()

	fmt.Printf("Enter task category: ")
	scanner.Scan()
	category = scanner.Text()

	fmt.Println("Task:", name, dueDate, category)
}

func createCategory() {
	scanner := bufio.NewScanner(os.Stdin)
	var title, color string

	fmt.Printf("Enter category title: ")
	scanner.Scan()
	title = scanner.Text()

	fmt.Printf("Enter category color: ")
	scanner.Scan()
	color = scanner.Text()

	fmt.Println("Category:", title, color)

}

func registerUser() {
	scanner := bufio.NewScanner(os.Stdin)
	var email, password string

	fmt.Printf("Enter user Email: ")
	scanner.Scan()
	email = scanner.Text()

	fmt.Printf("Enter user Password: ")
	scanner.Scan()
	password = scanner.Text()

	user := User{
		ID:       len(userStorage) + 1,
		Email:    email,
		Password: password,
	}

	userStorage = append(userStorage, user)

	fmt.Println("User:", user.ID, user.Email, user.Password)
}

func userLogin() {
	scanner := bufio.NewScanner(os.Stdin)
	var inputEmail, inputPassword string

	fmt.Printf("Enter Email: ")
	scanner.Scan()
	inputEmail = scanner.Text()

	fmt.Printf("Enter Password: ")
	scanner.Scan()
	inputPassword = scanner.Text()

	for _, user := range userStorage {
		if user.Email == inputEmail {
			if user.Password == inputPassword {
				fmt.Println("User Logged In Successfully!")
				return
			}
			fmt.Println("Wrong password!")
			return
		}
	}
	fmt.Println("User Not Found!")
}
