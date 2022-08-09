package controller

import (
	"github.com/XxThunderBlastxX/chamting-api/models"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

var server = &models.Server{}

func WsRoute(conn *websocket.Conn) {

	client := models.Client{
		Id:   uuid.Must(uuid.NewRandom()).String(),
		Conn: conn,
	}

	server.Send(&client, "Server: Welcome your Id is "+client.Id)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			server.RemoveClient(&client)
			return
		}
		server.ProcessMessage(client, messageType, p)
	}
}
