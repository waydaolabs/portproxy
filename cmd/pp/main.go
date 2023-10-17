package main

import (
	"portproxy/handler"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/urlBuildJJF7_2Tqm-zh4dbF", handler.UrlBuild)
	app.Get("/urlBuildJJF7_2Tqm-zh4dbF", handler.GetUrlBuild)
	app.All("*", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return websocket.New(func(c *websocket.Conn) {

			})(c)
		} else {
			return c.SendString(string(c.Context().Method()) + ":" + c.Hostname() + "/" + c.Params("*"))
		}

	})

	app.Listen(":9000")
}
