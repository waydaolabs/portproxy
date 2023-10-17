package handler

import (
	"log"
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
	if u.SSL {
		uu.Scheme = "s"
	}

	uu.Path, err = url.JoinPath(uu.Path, c.Params("*"))

	queries := uu.Query()
	for k, v := range c.Queries() {
		queries.Add(k, v)
	}
	uu.RawQuery = queries.Encode()

	if websocket.IsWebSocketUpgrade(c) {
		uu.Scheme = "ws" + uu.Scheme
		return websocket.New(func(c *websocket.Conn) {

		})(c)
	} else {
		uu.Scheme = "http" + uu.Scheme
		log.Println("request url", uu.String())
		a := fiber.AcquireAgent()

		req := a.Request()
		// set request body
		req.SetBody(c.BodyRaw())
		req.Header.SetMethod(c.Method())
		for k, v := range c.GetReqHeaders() {
			switch strings.ToLower(k) {
			case "connection":
			case "host":
				req.Header.SetHost(v)
			case "user-agent":
				a.UserAgent(v)
			default:
				req.Header.Add(k, v)
				log.Println("heaer:", k, "=>", v)
			}

		}

		req.SetRequestURI(uu.String())
		resp := fiber.AcquireResponse()
		defer fiber.ReleaseResponse(resp)
		a.SetResponse(resp)
		if err := a.Parse(); err != nil {
			panic(err)
		}

		code, body, _ := a.Bytes()
		resp.Header.VisitAll(func(key, value []byte) {
			c.Response().Header.Add(string(key), string(value))
		})
		c.Status(code)
		c.Write(body)
	}

	return nil

}
