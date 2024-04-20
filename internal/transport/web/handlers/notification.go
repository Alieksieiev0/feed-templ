package handlers

import (
	"bufio"
	"encoding/json"
	"log"
	"strconv"

	"github.com/Alieksieiev0/feed-templ/internal/services"
	"github.com/Alieksieiev0/feed-templ/internal/types"
	"github.com/Alieksieiev0/feed-templ/internal/view/notify"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

const (
	queryErr = "bad query params were found while preparing to load posts"
)

func GetNotificationsHandler(serv services.NotificationServices) fiber.Handler {
	return func(c *fiber.Ctx) error {
		limit, err := strconv.Atoi(c.Query("limit", "10"))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).Send([]byte("Error: " + queryErr))
		}

		status := c.Query("status")
		if status == "" {
			return c.Status(fiber.StatusInternalServerError).Send([]byte("Error: " + queryErr))
		}

		notifications, r := serv.Get(
			c.Context(),
			c.Cookies("id"),
			services.Limit(limit),
			services.DefaultParam("status", status),
		)
		if r.Err != nil {
			return r.SendError(c)
		}

		return Render(c, notify.Notifications(notifications))
	}
}

func ListenHandler(serv services.NotificationServices) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")

		ch := make(chan *types.Notification)
		err := serv.Listen(c.Context(), c.Cookies("id"), ch)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			for {
				notification := <-ch
				err := json.NewEncoder(w).Encode(notification)
				if err != nil {
					log.Println(err)
					break
				}

				err = w.Flush()
				if err != nil {
					log.Println(err)
					break
				}
			}
		}))

		return nil
	}
}
