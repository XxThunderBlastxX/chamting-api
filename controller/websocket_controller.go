package controller

import (
	"github.com/XxThunderBlastxX/chamting-api/database"
	"github.com/XxThunderBlastxX/chamting-api/models"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/nitishm/go-rejson/v4"
	"strings"
)

var (
	ServerInit = &models.Server{}
)

func WsRoute(conn *websocket.Conn) {
	clientId := conn.Query("id", strings.ReplaceAll(uuid.Must(uuid.NewRandom()).String(), "-", ""))

	client := models.Client{
		Id:   clientId,
		Conn: conn,
	}

	//ServerInit.Send(&client, "Server: Welcome your Id is "+client.Id)

	//Redis client instance
	models.RdbClient = database.RedisConnect(0)
	models.RdbChat = database.RedisConnect(1)

	//Redis json instance
	models.RJson = rejson.NewReJSONHandler()
	models.RJson.SetGoRedisClient(models.RdbChat)

	//closes the redis instances
	defer func() {
		if err := models.RdbClient.Close(); err != nil {
			return
		}
	}()
	defer func() {
		if err := models.RdbChat.Close(); err != nil {
			return
		}
	}()

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
