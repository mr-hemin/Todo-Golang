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
	"todo/repository/filestore"
	"todo/repository/memorystore"
	"todo/service/task"
)

func main() {
	taskMemoryRepo := memorystore.NewTaskStore()

	taskService := task.NewService(taskMemoryRepo)

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
		runCommand(UserFileStore, *command, &taskService)
		fmt.Printf("Enter a command %v: ", commandList)
		scanner.Scan()
		*command = scanner.Text()
	}

}

// Global variable - just reachable in this packge
var (
	userStorage     []entity.User
	categoryStorage []entity.Category

	authenticatedUser *entity.User
	serializationMode string
)

const (
	userStoragePath = "user.txt"
)

var commandList = []string{
	"register-user",
	"login-iser",
	"list-users",
	"create-category",
	"create-task",
	"list-tasks",
	"exit",
}

func runCommand(store contract.UserWriteStore, command string, taskService *task.Service) {
	if command != "register-user" && command != "exit" && authenticatedUser == nil {
		loginUser()

		if authenticatedUser == nil {
			return
		}

	}

	switch command {
	case "register-user":
		registerUser(store)
	case "login-user":
		loginUser()
	case "list-users":
		listUsers()
	case "create-category":
		createCategory(taskService)
	case "create-task":
		createTask(taskService)
	case "list-tasks":
		listTasks(taskService)
	case "exit":
		os.Exit(0)
	default:
		fmt.Println("Command not recognized")

	}
}

func createCategory(taskService *task.Service) {
	scanner := bufio.NewScanner(os.Stdin)
	var title, color string

	fmt.Printf("Enter category title: ")
	scanner.Scan()
	title = scanner.Text()

	fmt.Printf("Enter category color: ")
	scanner.Scan()
	color = scanner.Text()

	category := entity.Category{
		ID:     uint(len(categoryStorage) + 1),
		Title:  title,
		Color:  color,
		UserID: authenticatedUser.ID,
	}

	categoryStorage = append(categoryStorage, category)

	fmt.Println("Category:", title, color)

}

func createTask(taskService *task.Service) {
	fmt.Println("AuthenricatedUser Email:", authenticatedUser.Email)
	scanner := bufio.NewScanner(os.Stdin)
	var title, dueDate, categoryID string

	fmt.Printf("Enter task title: ")
	scanner.Scan()
	title = scanner.Text()

	fmt.Printf("Enter task dueDate: ")
	scanner.Scan()
	dueDate = scanner.Text()

	fmt.Printf("Enter task category-id: %v", categoryStorage)
	scanner.Scan()
	categoryID = scanner.Text()
	_, err := strconv.Atoi(categoryID)

	if err != nil {
		fmt.Printf("Category ID isn't a valid integer: %v\n", err)

		return

	}

	intCategoryID, _ := strconv.Atoi(categoryID)

	response, err := taskService.Create(task.CreateRequest{
		Title:               title,
		DueDate:             dueDate,
		CategoryID:          uint(intCategoryID),
		AuthenticatedUserID: authenticatedUser.ID,
	})

	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println("Task created successfuly. Task:", response)

}

func listTasks(taskService *task.Service) {
	userTasks, err := taskService.List(task.ListRequest{UserID: authenticatedUser.ID})
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println("user tasks:", userTasks)

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
func loginUser() {
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

func listUsers() {
	fmt.Println("*** Users ***")

	for _, user := range userStorage {
		fmt.Printf("%+v\n", user)
	}

	fmt.Println()
}
