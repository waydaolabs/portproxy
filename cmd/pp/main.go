package main

import (
	"portproxy/handler"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/urlBuildJJF7_2Tqm-zh4dbF", handler.UrlBuild)
	app.Get("/urlBuildJJF7_2Tqm-zh4dbF", handler.GetUrlBuild)
	app.All("*", handler.Proxy)

	app.Listen(":9000")
}
