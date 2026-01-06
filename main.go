package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	fmt.Println("*** Todo application ***")

	command := flag.String("command", "no-command", "Command to execute")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		runCommand(*command)
		fmt.Printf("Enter a command %v: ", commandList)
		scanner.Scan()
		*command = scanner.Text()
	}

}

type User struct {
	ID       uint
	Name     string
	Email    string
	Password string
}

type Task struct {
	ID       uint
	Title    string
	DueDate  string
	Category string
	IsDone   bool
	UserID   uint
}

// Global variable - just reachable in this packge
var userStorage []User
var taskStorage []Task
var authenticatedUser *User

var commandList = []string{
	"create-task",
	"task-list",
	"create-category",
	"register-user",
	"exit",
}

func runCommand(command string) {
	if command != "register-user" && command != "exit" && authenticatedUser == nil {
		userLogin()

		
		if authenticatedUser == nil {
			return
		}

	}


	switch command {
	case "create-task":
		createTask()
	case "create-category":
		createCategory()
	case "register-user":
		registerUser()
	case "task-list":
		taskList()
	case "user-login":
		userLogin()
	case "exit":
		os.Exit(0)
	default:
		fmt.Println("Command not recognized")

	}
}

func createTask() {
	fmt.Println("AuthenricatedUser Email:", authenticatedUser.Email)
	scanner := bufio.NewScanner(os.Stdin)
	var title, dueDate, category string

	fmt.Printf("Enter task title: ")
	scanner.Scan()
	title = scanner.Text()

	fmt.Printf("Enter task dueDate: ")
	scanner.Scan()
	dueDate = scanner.Text()

	fmt.Printf("Enter task category: ")
	scanner.Scan()
	category = scanner.Text()

	task := Task{
		ID:       uint(len(taskStorage)) + 1,
		Title:    title,
		DueDate:  dueDate,
		Category: category,
		IsDone:   false,
		UserID:   authenticatedUser.ID,
	}

	taskStorage = append(taskStorage, task)

	fmt.Println("Task:", title, dueDate, category)

}

func taskList() {
	for _, task := range taskStorage {
		if task.UserID == authenticatedUser.ID {
			fmt.Println(task)
		}

	}
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
	var name, email, password string

	fmt.Printf("Enter user name: ")
	scanner.Scan()
	name = scanner.Text()

	fmt.Printf("Enter user Email: ")
	scanner.Scan()
	email = scanner.Text()

	fmt.Printf("Enter user Password: ")
	scanner.Scan()
	password = scanner.Text()

	user := User{
		ID:       uint(len(userStorage)) + 1,
		Name:     name,
		Email:    email,
		Password: password,
	}

	userStorage = append(userStorage, user)

	fmt.Println("User:", user.ID, user.Name, user.Email, user.Password)
}

// Get user email and password
// Checks if there is a user record with corresponding data: allow user to continue
func userLogin() {
	fmt.Println("* User Login *")
	scanner := bufio.NewScanner(os.Stdin)
	var inputEmail, inputPassword string

	fmt.Printf("Enter Email: ")
	scanner.Scan()
	inputEmail = scanner.Text()

	fmt.Printf("Enter Password: ")
	scanner.Scan()
	inputPassword = scanner.Text()

	for _, user := range userStorage {
		if user.Email == inputEmail && user.Password == inputPassword {
			fmt.Println("User Logged In Successfully!")
			authenticatedUser = &user

			return
		}
	}

	fmt.Println("Email or Password isn't correct!")

}
