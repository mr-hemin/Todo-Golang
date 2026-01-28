package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"todo/delivery/deliveryparam"
	"todo/repository/memorystore"
	"todo/service/task"
)

func main() {
	const (
		network = "tcp"
		address = "127.0.0.1:1996"
	)

	// create new listener
	listener, err := net.Listen(network, address)
	if err != nil {
		log.Fatalln("can't listen on given address", address, err)
	}
	defer listener.Close()

	fmt.Println("server listening on", listener.Addr())

	taskMemoryRepo := memorystore.NewTaskStore()

	taskService := task.NewService(taskMemoryRepo)

	for {
		// listen for new connection
		connection, aErr := listener.Accept()
		if aErr != nil {
			log.Println("can't listen to new connection", aErr)

			continue
		}

		// process request
		var rawRequest = make([]byte, 1024)
		numberOfReadBytes, rErr := connection.Read(rawRequest)
		if rErr != nil {
			log.Println("can't read data from connection", rErr)

			continue
		}

		fmt.Printf("client address: %s, numOfReadBytes: %d, data: %s\n",
			connection.RemoteAddr(), numberOfReadBytes, string(rawRequest))

		req := &deliveryparam.Request{}
		if uErr := json.Unmarshal(rawRequest[:numberOfReadBytes], req); uErr != nil {
			log.Println("bad request", uErr)

			continue
		}

		switch req.Command {
		case "create-task":
			response, cErr := taskService.Create(task.CreateRequest{
				Title:               req.CreateTaskRequest.Title,
				DueDate:             req.CreateTaskRequest.DueDate,
				CategoryID:          req.CreateTaskRequest.CategoryID,
				AuthenticatedUserID: 0,
			})
			if cErr != nil {
				_, wErr := connection.Write([]byte(cErr.Error()))
				if wErr != nil {
					log.Println("can't write data to connection", rErr)

					continue
				}
			}

			data, mErr := json.Marshal(&response)
			if mErr != nil {
				_, wErr := connection.Write([]byte(mErr.Error()))
				if wErr != nil {
					log.Println("can't marshal response", rErr)

					continue
				}

				continue
			}

			_, wErr := connection.Write(data)
			if wErr != nil {
				log.Println("can't write data to connection", rErr)

				continue
			}
		}

		connection.Close()
	}
}
