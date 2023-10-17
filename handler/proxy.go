package handler

import (
	"net/url"
	"strings"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func Proxy(c *fiber.Ctx) (err error) {

	hs := strings.Split(c.Hostname(), ".")
	id := hs[0]

	urls := GetUrls()
	u, ok := urls[id]
	if !ok {
		return c.SendStatus(fiber.StatusNotFound)
	}

	uu := url.URL{
		Host:     u.Host,
		Path:     u.Path,
		RawQuery: u.Query,
	}

	uu.Path, err = url.JoinPath(uu.Path, c.Params("*"))

	if websocket.IsWebSocketUpgrade(c) {
		return websocket.New(func(c *websocket.Conn) {

		})(c)
	} else {
		a := fiber.AcquireAgent()

		defer fiber.ReleaseAgent(a)
		req := a.Request()
		req.Header.SetMethod(c.Method())
		for k, v := range c.GetReqHeaders() {
			req.Header.Add(k, v)
		}

		req.SetRequestURI("http://example.com")

		if err := a.Parse(); err != nil {
			panic(err)
		}

		code, body, _ := a.Bytes()

		c.Status(code)
		c.Write(body)
	}

	return nil

}
