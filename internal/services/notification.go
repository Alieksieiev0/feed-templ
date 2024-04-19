package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Alieksieiev0/feed-templ/internal/types"
	"github.com/gofiber/fiber/v2"
)

const (
	userNotificationsURL = "/api/notify/notifications/%s"
	userListenURL        = "/api/notify/listen/%s"
)

type NotificationServices interface {
	Get(c context.Context, userId string) ([]types.Notification, int, error)
	Listen(c context.Context, userId string, ch chan<- *types.Notification) error
}

func NewNotificationService(addr string) NotificationServices {
	return &notificationService{
		addr: addr,
	}
}

type notificationService struct {
	addr string
}

func (ns *notificationService) Get(
	c context.Context,
	userId string,
) ([]types.Notification, int, error) {
	fmt.Println(ns.addr)
	fmt.Println(userId)
	fmt.Printf(userNotificationsURL, userId)
	req, err := createRequest(
		c,
		http.MethodGet,
		ns.addr+fmt.Sprintf(userNotificationsURL, userId),
		nil,
	)

	if err != nil {
		return nil, fiber.StatusInternalServerError, fmt.Errorf("couldnt process provided data")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fiber.StatusInternalServerError, fmt.Errorf("couldnt fetch notifications")
	}

	fmt.Println(resp.StatusCode)
	if resp.StatusCode == http.StatusOK {
		notifications := []types.Notification{}
		err = json.NewDecoder(resp.Body).Decode(&notifications)
		if err != nil {
			return nil, fiber.StatusInternalServerError, fmt.Errorf(
				"couldnt verify received notifications",
			)
		}
		return notifications, resp.StatusCode, nil
	}

	return nil, resp.StatusCode, readResponseError(resp)
}

func (ns *notificationService) Listen(
	c context.Context,
	userId string,
	ch chan<- *types.Notification,
) error {
	conn, err := createWebsocketRequest(ns.addr, userListenURL)
	if err != nil {
		return err
	}
	return listenWebsocket(conn, ch)
}
