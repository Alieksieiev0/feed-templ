package services

import (
	"context"
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
	Get(c context.Context, userId string, params ...Param) ([]types.Notification, *Response)
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
	params ...Param,
) ([]types.Notification, *Response) {
	req, err := createRequest(
		c,
		http.MethodGet,
		ns.addr+fmt.Sprintf(userNotificationsURL, userId),
		nil,
	)

	if err != nil {
		return nil, NewResponse(fiber.StatusInternalServerError, err)
	}

	updateQuery(req, params)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, NewResponse(fiber.StatusInternalServerError, err)
	}

	if resp.StatusCode == http.StatusOK {
		notifications := []types.Notification{}
		return notifications, parseResponse(resp, &notifications)
	}

	return nil, NewResponse(resp.StatusCode, readResponseError(resp))
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
