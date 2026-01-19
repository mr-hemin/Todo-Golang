package filestore

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"todo/constant"
	"todo/entity"
)

type FileStore struct {
	filePath          string
	serializationMode string
}

// Constructor
func New(path, serializationMode string) FileStore {
	return FileStore{filePath: path, serializationMode: serializationMode}
}

func (f FileStore) Save(u entity.User) {
	f.writeUserToFile(u)
}

func (f FileStore) Load() []entity.User {
	var uStorage []entity.User

	file, err := os.Open(f.filePath)
	if err != nil {
		fmt.Println("Can't open the file:", err)

		return nil
	}

	defer file.Close()

	var data = make([]byte, 1024)

	_, opnErr := file.Read(data)

	if opnErr != nil {
		fmt.Println("Can't read from the file:", opnErr)

		return nil
	}

	dataStrSlice := strings.Split(string(data), "\n")

	for _, userLine := range dataStrSlice {

		var userStruct = entity.User{}

		switch f.serializationMode {
		case constant.MySerializationMode:
			var dSrlzErr error
			userStruct, dSrlzErr = deSerializeFromMyMode(userLine)

			if dSrlzErr != nil {
				fmt.Println("can't deserialize user record to usr struct:", dSrlzErr)

				return nil
			}

		case constant.JsonSerializationMode:
			if userLine[0] != '{' && userLine[len(userLine)-1] != '}' {

				continue
			}

			uMrshErr := json.Unmarshal([]byte(userLine), &userStruct)
			if uMrshErr != nil {
				fmt.Println("can't deserialize user record to usr struct with json mode:", uMrshErr)

				return nil
			}

		default:
			fmt.Println("invalid serializationMode!", f.serializationMode)

			return nil
		}

		uStorage = append(uStorage, userStruct)
	}

	return uStorage

}

func (f FileStore) writeUserToFile(user entity.User) {
	// Save user data  in user.txt file
	var file *os.File

	file, err := os.OpenFile(f.filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Can't create or open file:", err)

		return
	}

	defer file.Close()
	var data []byte

	if f.serializationMode == constant.MySerializationMode {
		// Serialize the user struct/object

		data = []byte(fmt.Sprintf("ID: %d, Name: %s, Email: %s, Password: %s\n", user.ID, user.Name, user.Email, user.Password))

	} else if f.serializationMode == constant.JsonSerializationMode {
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

func deSerializeFromMyMode(userStr string) (entity.User, error) {
	if userStr == "" {
		return entity.User{}, errors.New("user string is empaty")
	}

	userFileds := strings.Split(userStr, ", ")
	var user = entity.User{}

	for _, field := range userFileds {
		values := strings.Split(field, ": ")
		if len(values) != 2 {
			fmt.Println("Field isn't valid skipping...", len(values))
			return entity.User{}, errors.New("strconv error")
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
