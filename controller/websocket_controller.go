package controller

import (
	"github.com/XxThunderBlastxX/chamting-api/models"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

var ServerInit = &models.Server{}

func WsRoute(conn *websocket.Conn) {
	clientId := conn.Query("id", uuid.Must(uuid.NewRandom()).String())

	client := models.Client{
		Id:   clientId,
		Conn: conn,
	}

	ServerInit.Send(&client, "Server: Welcome your Id is "+client.Id)

	for {
		_, payLoad, err := conn.ReadMessage()
		if err != nil {
			ServerInit.RemoveClient(&client)
			return
		}

		//sending data to go routine
		models.Cli <- client
		models.PayLoad <- payLoad
	}
}
