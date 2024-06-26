package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/Alieksieiev0/feed-templ/internal/types"
	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
)

func createRequest(c context.Context, method, url string, v any) (*http.Request, error) {
	reader := &bytes.Reader{}
	if v != nil {
		body, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		reader.Reset(body)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func updateQuery(req *http.Request, params []Param) {
	if len(params) == 0 {
		return
	}
	q := req.URL.Query()
	ApplyParams(q, params)
	req.URL.RawQuery = q.Encode()
}

func parseResponse[T any](resp *http.Response, entity *T) *Response {
	err := json.NewDecoder(resp.Body).Decode(entity)
	if err != nil {
		return NewResponse(fiber.StatusInternalServerError, err)
	}
	return NewResponse(resp.StatusCode, err)
}

func readResponseError(resp *http.Response) error {
	respErr := &types.ResponseError{}
	err := json.NewDecoder(resp.Body).Decode(respErr)
	if err != nil {
		return fmt.Errorf("couldnt process results of operation")
	}
	return fmt.Errorf(respErr.Error)
}

func createWebsocketRequest(host, path string) (*websocket.Conn, error) {
	URL := url.URL{Scheme: "ws", Host: host, Path: path}
	conn, _, err := websocket.DefaultDialer.Dial(URL.String(), nil)
	return conn, err
}

func listenWebsocket[T any](conn *websocket.Conn, ch chan<- *T) error {
	go func() {
		defer close(ch)
		for {
			entity := new(T)
			err := conn.ReadJSON(entity)
			if err != nil {
				log.Println(err)
				return
			}
			ch <- entity
		}
	}()

	return nil
}
