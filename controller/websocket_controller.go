package controller

import (
	webs "github.com/XxThunderBlastxX/chamting-api/websocket"
	"github.com/gofiber/websocket/v2"
	"log"
)

func WsRoute(conn *websocket.Conn) {
	var room webs.Room
	var msg webs.Message
	var wsRoutine webs.Ws

	go wsRoutine.RunHub()

	// When the function returns, unregister the client and close the connection
	defer func() {
		room.Unregister <- conn
		conn.Close()
	}()

	// Register the client
	room.Register <- conn

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("read error:", err)
			}

			return // Calls the deferred function, i.e. closes the connection on error
		}

		if messageType == websocket.TextMessage {
			// Broadcast the received message
			msg.Msg <- string(message)
		} else {
			log.Println("websocket message received of type", messageType)
		}
	}

}
