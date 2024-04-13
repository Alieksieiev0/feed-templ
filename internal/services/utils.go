package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Alieksieiev0/feed-templ/internal/types"
)

func ReadResponseError(resp *http.Response) error {
	respErr := &types.ResponseError{}
	err := json.NewDecoder(resp.Body).Decode(respErr)
	if err != nil {
		return fmt.Errorf("couldn`t process results of operation")
	}
	return fmt.Errorf(respErr.Error)
}

func SendRequest(
	ctx context.Context,
	url string,
	v any,
	headers ...map[string]string,
) (*http.Response, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	if len(headers) > 0 {
		for k, v := range headers[0] {
			req.Header.Add(k, v)
		}
	}

	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}
