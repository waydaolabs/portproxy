package handler

import (
	"github.com/gofiber/fiber/v2"
)

func UrlBuild(c *fiber.Ctx) error {
	data := Url{}

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	SetUrls(data)
	return c.JSON(GetUrls())
}

func GetUrlBuild(c *fiber.Ctx) error {

	return c.JSON(GetUrls())
}
