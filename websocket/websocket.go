package websocket

import (
	"github.com/gofiber/websocket/v2"
	"log"
	"time"
)

type client struct{} // Add more data to this type if needed

type Ws interface {
	RunHub()
}

type RoomInstance struct {
	Room    Room
	Client  Client
	Message Message
}

func NewWsInstance(roomId string) Ws {
	return &RoomInstance{Room: Room{RoomId: roomId}}
}

type Room struct {
	RoomId     string
	Clients    map[*websocket.Conn]Client
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
}

type Client struct {
	ClientId string
	Author   string
	Msg      Message
}

type Message struct {
	Author string
	Msg    chan string
	Time   time.Time
}

//var (
//	Clients    = make(map[*websocket.Conn]client) // Note: although large maps with pointer-like types (e.g. strings) as keys are slow, using pointers themselves as keys is acceptable and fast
//	Register   = make(chan *websocket.Conn)
//	Broadcast  = make(chan string)
//	Unregister = make(chan *websocket.Conn)
//)

func (r *RoomInstance) RunHub() {
	for {
		select {
		case connection := <-r.Room.Register:
			r.Room.Clients[connection] = r.Client

			log.Println("connection registered")

		case message := <-r.Message.Msg:
			log.Println("message received:", message)

			// Send the message to all clients
			for connection := range r.Room.Clients {
				if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					log.Println("write error:", err)

					connection.WriteMessage(websocket.CloseMessage, []byte{})
					connection.Close()
					delete(r.Room.Clients, connection)
				}
			}

		case connection := <-r.Room.Unregister:
			// Remove the client from the hub
			delete(r.Room.Clients, connection)

			log.Println("connection unregistered")
		}

	}

}
