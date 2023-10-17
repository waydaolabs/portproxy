package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func UrlBuild(c *fiber.Ctx) error {
	data := Url{}

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if len(data.Id) < 1 {
		return errors.New("id error")
	}
	if len(data.Host) < 1 {
		return errors.New("host error")
	}

	SetUrls(data)
	return c.JSON(GetUrls())
}

func GetUrlBuild(c *fiber.Ctx) error {
	return c.JSON(GetUrls())
}
