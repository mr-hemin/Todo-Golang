package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strconv"
	"todo/constant"
	"todo/contract"
	"todo/entity"
	"todo/filestore"
)

func main() {
	serializeMode := flag.String("serialize-mode", constant.MySerializationMode, "serialization mode to write data to file")
	command := flag.String("command", "no-command", "Command to execute")
	flag.Parse()

	fmt.Println("*** Todo application ***")

	switch *serializeMode {
	case constant.MySerializationMode:
		serializationMode = constant.MySerializationMode
	default:
		serializationMode = constant.JsonSerializationMode
	}

	var UserFileStore = filestore.New(userStoragePath, serializationMode)

	// Load userx storage from file
	users := UserFileStore.Load()

	userStorage = append(userStorage, users...)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println(userStorage)
		runCommand(UserFileStore, *command)
		fmt.Println("Users:", userStorage)
		fmt.Printf("Enter a command %v: ", commandList)
		scanner.Scan()
		*command = scanner.Text()
	}

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
var (
	userStorage     []entity.User
	taskStorage     []Task
	categoryStorage []Category

	authenticatedUser *entity.User
	serializationMode string
)

const (
	userStoragePath = "user.txt"
)

var commandList = []string{
	"create-task",
	"task-list",
	"create-category",
	"user-login",
	"register-user",
	"exit",
}

func runCommand(store contract.UserWriteStore, command string) {
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
		registerUser(store)
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

func registerUser(store contract.UserWriteStore) {
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

	user := entity.User{
		ID:       uint(len(userStorage)) + 1,
		Name:     name,
		Email:    email,
		Password: hashPassword(password),
	}

	userStorage = append(userStorage, user)

	// writeUserToFile(user)
	store.Save(user)

}

func hashPassword(password string) string {
	hash := md5.Sum([]byte(password))

	return hex.EncodeToString(hash[:])
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
		if user.Email == inputEmail && user.Password == hashPassword(inputPassword) {
			fmt.Println("User Logged In Successfully!")
			authenticatedUser = &user

			return
		}
	}

	fmt.Println("Email or Password isn't correct!")

}
