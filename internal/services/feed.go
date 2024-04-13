package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Alieksieiev0/feed-templ/internal/types"
	"github.com/gofiber/fiber/v2"
)

type FeedService interface {
	Post(ctx context.Context, token string, post *types.Post) (int, error)
}

type feedService struct {
	addr    string
	postURL string
}

func (fs *feedService) Post(ctx context.Context, token string, post *types.Post) (int, error) {
	resp, err := SendRequest(
		ctx,
		fs.addr+fs.postURL,
		post,
		map[string]string{"Authorization": token},
	)
	if err != nil {
		return fiber.StatusInternalServerError, fmt.Errorf("couldn`t create post")
	}

	if resp.StatusCode == fiber.StatusCreated {
		err = json.NewDecoder(resp.Body).Decode(post)
		if err != nil {
			return fiber.StatusInternalServerError, fmt.Errorf("couldn`t verify post creation")
		}
		return resp.StatusCode, nil
	}

	return resp.StatusCode, ReadResponseError(resp)

}
