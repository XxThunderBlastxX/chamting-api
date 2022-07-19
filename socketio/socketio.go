package socketio

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	socketio "github.com/googollee/go-socket.io"
	"log"
)

func SocketIo() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		server := socketio.NewServer(nil)

		server.OnConnect("/", func(conn socketio.Conn) error {
			conn.SetContext(context.TODO())
			log.Println("New Connection:", conn.ID())
			return nil
		})

		server.OnEvent("/", "notice", func(conn socketio.Conn, msg string) {
			fmt.Println("notice:", msg)
			conn.Emit("reply", "have : "+msg)
		})

		server.OnEvent("/chat", "msg", func(conn socketio.Conn, msg string) string {
			conn.SetContext(msg)
			return "receive : " + msg
		})

		server.OnEvent("/", "bye", func(conn socketio.Conn) string {
			last := conn.Context().(string)

			conn.Emit("bye", last)
			conn.Close()
			return last
		})

		server.OnError("/", func(conn socketio.Conn, err error) {
			fmt.Println("error : ", err)
		})

		server.OnDisconnect("/", func(conn socketio.Conn, reason string) {
			fmt.Println("closed : ", reason)
		})

		go server.Serve()

		defer server.Close()

		return nil
	}

}
