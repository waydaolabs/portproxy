package main

import (
	"flag"

	"github.com/waydaolabs/portproxy/config"
	"github.com/waydaolabs/portproxy/handler"

	"github.com/gofiber/fiber/v2"
)

func init() {
	flag.Bool("multicore", false, "")
	flag.StringVar(&config.DB_FILE, "db", config.DB_FILE, "--db "+config.DB_FILE)
	flag.StringVar(&config.UrlBuildPath, "url", config.UrlBuildPath, "--url "+config.UrlBuildPath)
	flag.Parse()

}

func main() {
	app := fiber.New()

	app.Post(config.UrlBuildPath, handler.UrlBuild)
	app.Get(config.UrlBuildPath, handler.GetUrlBuild)
	app.All("*", handler.Proxy)

	app.Listen(":9000")
}
