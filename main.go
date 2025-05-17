package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// estructura para recibir mensajes
type Message struct {
	User    string `json:"username"`
	Message string `json:"message"`
}

// Upgrader is used to upgrade HTTP connections to WebSocket connections

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool) //connected clients
var broadcast = make(chan []byte)            // broadcast channel
var mutex = &sync.Mutex{}                    // protect client maps

func wsHandler(w http.ResponseWriter, r *http.Request) {

	// Upgrade the http connection to a websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("error upgrading: ", err)
		return
	}
	defer conn.Close()
	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	for {

		_, message, err := conn.ReadMessage()

		if err != nil {
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}

		var msg Message
		err = json.Unmarshal(message, &msg)

		if err != nil {
			fmt.Println("No se pudo descifrar el msg")
			continue
		}
		fmt.Printf("User: %s  says: %s \n", msg.User, msg.Message)
		encoded, err := json.Marshal(msg)

		if err != nil {
			fmt.Println("Error al parsear en go ")
		}
		broadcast <- encoded

	}
}
func handleMessage() {
	for {
		// grab the next message from the broadcast channel
		message := <-broadcast

		// extraemos el json

		// send the message to all clients conencted
		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

/*
	func handleConnection(conn *websocket.Conn) {
		// listen for incoming messages
		for {
			_, message, err := conn.ReadMessage()

			if err != nil {
				fmt.Println("Error reading line: ", err)
				break
			}
			fmt.Printf("Received: %s \n", message)
			// echo the message back to the client
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				fmt.Println("Error writing message", err)
				break
			}
		}
	}
*/
func main() {
	http.HandleFunc("/ws", wsHandler)
	go handleMessage()
	fmt.Println("Web socket started in port :8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Error starting server", err)
	}

}
