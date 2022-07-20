package websocket

import (
	"github.com/gofiber/websocket/v2"
	"log"
)

type client struct{} // Add more data to this type if needed

var (
	Clients    = make(map[*websocket.Conn]client) // Note: although large maps with pointer-like types (e.g. strings) as keys are slow, using pointers themselves as keys is acceptable and fast
	Register   = make(chan *websocket.Conn)
	Broadcast  = make(chan string)
	Unregister = make(chan *websocket.Conn)
)

func RunHub() {
	for {
		select {
		case connection := <-Register:
			Clients[connection] = client{}
			log.Println("connection registered")

		case message := <-Broadcast:
			log.Println("message received:", message)

			// Send the message to all clients
			for connection := range Clients {
				if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					log.Println("write error:", err)

					connection.WriteMessage(websocket.CloseMessage, []byte{})
					connection.Close()
					delete(Clients, connection)
				}
			}

		case connection := <-Unregister:
			// Remove the client from the hub
			delete(Clients, connection)

			log.Println("connection unregistered")
		}
	}
}
