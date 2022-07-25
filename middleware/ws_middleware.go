package middleware

import (
	ws "github.com/XxThunderBlastxX/chamting-api/websocket"
	"github.com/gofiber/fiber/v2"
)

func WsGetRoomId() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		roomId := ctx.Query("room_id", "nilvalue")
		//log.Println(roomId)
		ws.NewWsInstance(roomId)

		return ctx.Next()
	}

}
