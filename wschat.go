package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients []websocket.Conn

func main() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
		if err != nil {
			log.Fatal(err)
		}
		clients = append(clients, *conn)

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Fatal(err)
				return
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			for _, client := range clients {
				// Write message back to browser
				if err = client.WriteMessage(msgType, msg); err != nil {
					return
				}
			}

		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
