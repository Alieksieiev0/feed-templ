package services

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Alieksieiev0/feed-templ/internal/types"
	"github.com/gofiber/fiber/v2"
)

const (
	notificationURL = "/api/notify/notifications/%s"
	reviewURL       = "/api/notify/review/%s"
	listenURL       = "/api/notify/listen/%s"
)

type NotificationService interface {
	Get(c context.Context, userId string, params ...Param) ([]types.Notification, *Response)
	Listen(c context.Context, userId string, ch chan<- *types.Notification) error
	Review(c context.Context, id string) *Response
}

func NewNotificationService(addr string) NotificationService {
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
		ns.addr+fmt.Sprintf(notificationURL, userId),
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
	URL, err := url.Parse(ns.addr)
	if err != nil {
		return err
	}

	conn, err := createWebsocketRequest(URL.Host, fmt.Sprintf(listenURL, userId))
	if err != nil {
		return err
	}
	return listenWebsocket(conn, ch)
}

func (ns *notificationService) Review(c context.Context, id string) *Response {
	req, err := createRequest(c, http.MethodPut, ns.addr+fmt.Sprintf(reviewURL, id), nil)
	fmt.Println(ns.addr + fmt.Sprintf(reviewURL, id))

	if err != nil {
		fmt.Println(err)
		return NewResponse(fiber.StatusInternalServerError, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return NewResponse(fiber.StatusInternalServerError, err)
	}

	if resp.StatusCode == http.StatusOK {
		return NewResponse(resp.StatusCode, nil)
	}

	fmt.Println(resp.StatusCode)
	return NewResponse(resp.StatusCode, readResponseError(resp))
}
