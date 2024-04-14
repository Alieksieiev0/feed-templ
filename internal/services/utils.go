package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Alieksieiev0/feed-templ/internal/types"
)

func readResponseError(resp *http.Response) error {
	respErr := &types.ResponseError{}
	err := json.NewDecoder(resp.Body).Decode(respErr)
	fmt.Println(err)
	if err != nil {
		return fmt.Errorf("couldn`t process results of operation")
	}
	return fmt.Errorf(respErr.Error)
}

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
