package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"todo/delivery/tcp/deliveryparam"
	"todo/repository/memorystore"
	"todo/service/task"
)

func main() {
	const (
		network = "tcp"
		address = "127.0.0.1:1996"
	)

	// Create new listener
	listener, err := net.Listen(network, address)
	if err != nil {
		log.Fatalln("can't listen on given address:", address, err)
	}

	defer listener.Close()

	taskMemoryRepo := memorystore.NewTaskStore()
	taskService := task.NewService(taskMemoryRepo)

	fmt.Println("Server is listening on:", listener.Addr())

	for {
		//  Listen for new connection
		connection, cErr := listener.Accept()
		if cErr != nil {
			log.Println("can't listen to new connection:", cErr)

			continue
		}

		// process request
		var rawRequest = make([]byte, 1024)
		numberOfReadBytes, rErr := connection.Read(rawRequest)
		if rErr != nil {
			log.Println("can't read data from connection:", rErr)

			continue
		}

		fmt.Printf("Client address: %s, numberOfReadBytes: %d, data: %s\n",
			connection.RemoteAddr(), numberOfReadBytes, string(rawRequest))

		req := &deliveryparam.Request{}
		if uErr := json.Unmarshal(rawRequest[:numberOfReadBytes], req); uErr != nil {
			log.Println("bad request:", uErr)

			continue
		}

		switch req.Command {
		case "create-task":
			response, cErr := taskService.Create(task.CreateRequest{
				Title:               "",
				DueDate:             "",
				CategoryID:          0,
				AuthenticatedUserID: 0,
			})

			if cErr != nil {
				_, wErr := connection.Write([]byte(cErr.Error()))
				if wErr != nil {
					log.Println("can't write data to connection:", wErr)

					continue
				}

			}

			data, mErr := json.Marshal(response)
			if mErr != nil {
				_, wErr := connection.Write([]byte(mErr.Error()))
				if wErr != nil {
					log.Println("can't marshal response:", wErr)

					continue
				}

				continue
			}

			_, wErr := connection.Write(data)
			if wErr != nil {
				log.Println("can't write data to connection:", wErr)

				continue
			}

		}

		connection.Close()

	}

}
