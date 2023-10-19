package handler

import (
	"log"
	"net/url"
	"strings"

	fastwebsocket "github.com/fasthttp/websocket"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/sync/errgroup"
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
			func(c *websocket.Conn) {
				wg := new(errgroup.Group)

				rc, _, err := fastwebsocket.DefaultDialer.Dial(uu.String(), nil)
				if err != nil {
					log.Println("dial error", err, "url", uu.String())
					return
				}
				defer rc.Close()

				wg.Go(func() error {
					for {
						messageType, message, err := rc.ReadMessage()
						if err != nil {
							log.Println("remote read error", err)
							return err
						}
						err = c.WriteMessage(messageType, message)

						if err != nil {
							log.Println("client write error", err)
							return err
						}
					}
				})
				wg.Go(func() error {
					for {
						messageType, message, err := c.ReadMessage()
						if err != nil {
							log.Println("client read error", err)
							return err
						}
						err = rc.WriteMessage(messageType, message)
						if err != nil {
							log.Println("remote write error", err)
							return err
						}
					}
				})
				err = wg.Wait()
				if err != nil {
					log.Println("ws error", err)
				}
			}(c)
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
				continue
			case "host":
				req.Header.SetHost(v)
			case "user-agent":
				a.UserAgent(v)
			default:
				req.Header.Add(k, v)
				log.Println("request heaer:", k, "=>", v)
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
			log.Println("response heaer:", string(key), "=>", string(value))
			c.Response().Header.Add(string(key), string(value))
		})
		c.Status(code)
		c.Write(body)
	}

	return nil

}
