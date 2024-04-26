package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"strconv"
	"time"

	"github.com/Alieksieiev0/feed-templ/internal/services"
	"github.com/Alieksieiev0/feed-templ/internal/types"
	"github.com/Alieksieiev0/feed-templ/internal/view/notify"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

const (
	YYYYMMDD = "2006-01-02"
	queryErr = "bad query params were found while preparing to load posts"
)

func GetNotificationsHandler(serv services.NotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		limit, err := strconv.Atoi(c.Query("limit", "10"))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).Send([]byte("Error: " + queryErr))
		}

		notifications, r := serv.Get(
			c.Context(),
			c.Cookies("id"),
			services.DefaultParam("created_at", time.Now().UTC().Format(YYYYMMDD)),
			services.Limit(limit),
			services.SortBy("created_at"),
			services.OrderBy("desc"),
		)
		if r.Err != nil {
			return r.SendError(c)
		}

		return Render(c, notify.Notifications(notifications), templ.WithStatus(r.StatusCode))
	}
}

func ListenHandler(serv services.NotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		ch := make(chan *types.Notification)
		err := serv.Listen(c.Context(), c.Cookies("id"), ch)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			for {
				notification := <-ch
				if notification == nil {
					break
				}

				err = json.NewEncoder(w).Encode(notification)
				if err != nil {
					log.Println(err)
					break
				}

				var html template.HTML
				html, err = templ.ToGoHTML(
					context.Background(),
					notify.Notification(*notification),
				)
				if err != nil {
					log.Println(err)
					break
				}

				fmt.Fprintf(w, "data: %v \n\n", html)

				err = w.Flush()
				if err != nil {
					log.Println(err)
					break
				}
			}
		}))

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		c.Status(fiber.StatusOK)
		return nil
	}
}

func ReviewHandler(serv services.NotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println(11111)
		id := c.Params("id")
		r := serv.Review(c.Context(), id)
		if r.Err != nil {
			return r.SendError(c)
		}

		notifications, r := serv.Get(
			c.Context(),
			c.Cookies("id"),
			services.Limit(10),
			services.DefaultParam("created_at", time.Now().UTC().Format(YYYYMMDD)),
			services.SortBy("created_at"),
			services.OrderBy("desc"),
		)
		if r.Err != nil {
			return r.SendError(c)
		}

		return Render(c, notify.Notifications(notifications), templ.WithStatus(fiber.StatusOK))
	}
}
