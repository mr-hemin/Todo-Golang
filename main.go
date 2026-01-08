package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Load user storage from file
	loadUserStorageFromeFile()

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
	ID         uint
	Title      string
	DueDate    string
	CategoryID uint
	IsDone     bool
	UserID     uint
}

type Category struct {
	ID     uint
	Title  string
	Color  string
	UserID uint
}

// Global variable - just reachable in this packge
var userStorage []User
var taskStorage []Task
var categoryStorage []Category
var authenticatedUser *User

const userStoragePath = "user.txt"

var commandList = []string{
	"create-task",
	"task-list",
	"create-category",
	"user-login",
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

	fmt.Printf("Enter task category-id: %v", categoryStorage)
	scanner.Scan()
	category = scanner.Text()
	categoryID, err := strconv.Atoi(category)

	if err != nil {
		fmt.Printf("Category ID isn't a valid integer: %v\n", err)

		return

	}

	isFound := false

	for _, ctg := range categoryStorage {
		if ctg.ID == uint(categoryID) && ctg.UserID == authenticatedUser.ID {
			isFound = true

			break
		}
	}

	if !isFound {
		fmt.Println("Category ID not found!")

		return
	}

	task := Task{
		ID:         uint(len(taskStorage)) + 1,
		Title:      title,
		DueDate:    dueDate,
		CategoryID: uint(categoryID),
		IsDone:     false,
		UserID:     authenticatedUser.ID,
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

	category := Category{
		ID:     uint(len(categoryStorage) + 1),
		Title:  title,
		Color:  color,
		UserID: authenticatedUser.ID,
	}

	categoryStorage = append(categoryStorage, category)

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

	writeUserToFile(user)

}

func writeUserToFile(user User) {
	// Save user data  in user.txt file
	var file *os.File

	file, err := os.OpenFile(userStoragePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Can't create or open file:", err)

		return
	}

	defer file.Close()

	data := fmt.Sprintf("ID: %d, Name: %s, Email: %s, Password: %s\n", user.ID, user.Name, user.Email, user.Password)
	numberOfWrittenBytes, wrtErr := file.Write([]byte(data))
	if wrtErr != nil {
		fmt.Printf("Can't write to the file :%v\n", wrtErr)

		return
	}

	fmt.Println("Number of written bytes:", numberOfWrittenBytes)

}

func loadUserStorageFromeFile() {
	file, err := os.Open(userStoragePath)
	if err != nil {
		fmt.Println("Can't open the file:", err)
	}

	defer file.Close()

	var data = make([]byte, 10240)

	_, opnErr := file.Read(data)

	if opnErr != nil {
		fmt.Println("Can't read from the file:", opnErr)
	}

	dataStrSlice := strings.Split(string(data), "\n")

	fmt.Println(dataStrSlice)
	fmt.Println()

	for _, userLine := range dataStrSlice {
		if userLine == "" {
			continue
		}

		userFileds := strings.Split(userLine, ", ")
		var user = User{}

		for _, field := range userFileds {
			values := strings.Split(field, ": ")
			if len(values) != 2 {
				fmt.Println("Field isn't valid skipping...", len(values))
				continue
			}

			fieldName := strings.ReplaceAll(values[0], " ", "")
			fieldValue := values[1]

			switch fieldName {
			case "ID":
				id, err := strconv.Atoi(fieldValue)
				if err != nil {
					fmt.Println("strconv error:", err)
				}

				user.ID = uint(id)

			case "Name":
				user.Name = fieldValue
			case "Email":
				user.Email = fieldValue
			case "Password":
				user.Password = fieldValue
			}

		}

		fmt.Printf("User: %+v\n", user)

		fmt.Println()
	}

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
