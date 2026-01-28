package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"todo/delivery/deliveryparam"
)

func main() {
	fmt.Println("command", os.Args[0])

	message := "default message"

	if len(os.Args) > 1 {
		message = os.Args[1]
	}

	connection, err := net.Dial("tcp", "127.0.0.1:1996")
	if err != nil {
		log.Fatalln("can't dial the given address:", err)
	}

	defer connection.Close()

	fmt.Println("local address:", connection.LocalAddr())

	req := deliveryparam.Request{
		Command: message,
	}

	if req.Command == "create-task" {
		req.CreateTaskRequest = deliveryparam.CreateTaskRequest{
			Title:      req.CreateTaskRequest.Title,
			DueDate:    req.CreateTaskRequest.DueDate,
			CategoryID: req.CreateTaskRequest.CategoryID,
		}
	}

	serializedData, mErr := json.Marshal(req.Command)
	if mErr != nil {
		log.Fatalln("can't marshal request:", mErr)
	}

	numberOfWrittenBytes, wErr := connection.Write([]byte(serializedData))
	if wErr != nil {
		log.Fatalln("can't write data to connection:", wErr)
	}

	fmt.Println("numberOfWrittenBytes:", numberOfWrittenBytes)

	var data = make([]byte, 1024)
	_, rErr := connection.Read(data)
	if rErr != nil {
		log.Fatalln("can't read data from connection:", rErr)
	}

	fmt.Println("server response:", string(data))

}
