package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	serializeMode := flag.String("serialize-mode", MySerializationMode, "serialization mode to write data to file")
	command := flag.String("command", "no-command", "Command to execute")
	flag.Parse()
	// Load user storage from file
	loadUserStorageFromeFile(*serializeMode)

	fmt.Println("*** Todo application ***")

	switch *serializeMode {
	case MySerializationMode:
		serializationMode = MySerializationMode
	default:
		serializationMode = JsonSerializationMode
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println(userStorage)
		runCommand(*command)
		fmt.Println("Users:", userStorage)
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
var (
	userStorage     []User
	taskStorage     []Task
	categoryStorage []Category

	authenticatedUser *User
	serializationMode string
)

const (
	userStoragePath       = "user.txt"
	MySerializationMode   = "myCustom"
	JsonSerializationMode = "json"
)

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
	var data []byte

	if serializationMode == MySerializationMode {
		// Serialize the user struct/object

		data = []byte(fmt.Sprintf("ID: %d, Name: %s, Email: %s, Password: %s\n", user.ID, user.Name, user.Email, user.Password))

	} else if serializationMode == JsonSerializationMode {
		var jsonErr error
		data, jsonErr = json.Marshal(user)

		if jsonErr != nil {
			fmt.Println("Can't marshal user struct to json:", jsonErr)

			return
		}

		data = append(data, []byte("\n")...)

	} else {
		fmt.Println("Invalid serialization mode")

		return
	}

	numberOfWrittenBytes, wrtErr := file.Write(data)
	if wrtErr != nil {
		fmt.Printf("Can't write to the file :%v\n", wrtErr)

		return
	}

	fmt.Println("Number of written bytes:", numberOfWrittenBytes)

}

func loadUserStorageFromeFile(serializeMode string) {
	file, err := os.Open(userStoragePath)
	if err != nil {
		fmt.Println("Can't open the file:", err)

		return
	}

	defer file.Close()

	var data = make([]byte, 10240)

	_, opnErr := file.Read(data)

	if opnErr != nil {
		fmt.Println("Can't read from the file:", opnErr)

		return
	}

	dataStrSlice := strings.Split(string(data), "\n")

	for _, userLine := range dataStrSlice {

		var userStruct = User{}

		switch serializeMode {
		case MySerializationMode:
			var dSrlzErr error
			userStruct, dSrlzErr = deSerializeFromMyMode(userLine)

			if dSrlzErr != nil {
				fmt.Println("can't deserialize user record to usr struct:", dSrlzErr)

				return
			}

		case JsonSerializationMode:
			if userLine[0] != '{' && userLine[len(userLine)-1] != '}' {

				continue
			}

			uMrshErr := json.Unmarshal([]byte(userLine), &userStruct)
			if uMrshErr != nil {
				fmt.Println("can't deserialize user record to usr struct with json mode:", uMrshErr)

				return
			}

		default:
			fmt.Println("invalid serializationMode!", serializeMode)

			return
		}

		userStorage = append(userStorage, userStruct)
	}

}

func deSerializeFromMyMode(userStr string) (User, error) {
	if userStr == "" {
		return User{}, errors.New("user string is empaty")
	}

	userFileds := strings.Split(userStr, ", ")
	var user = User{}

	for _, field := range userFileds {
		values := strings.Split(field, ": ")
		if len(values) != 2 {
			fmt.Println("Field isn't valid skipping...", len(values))
			return User{}, errors.New("strconv error")
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

	return user, nil
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
