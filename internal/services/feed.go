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
	postsURL     = "/api/feed/posts"
	postURL      = "/api/feed/users/%s/posts"
	subscribeURL = "/api/feed/users/%s/subscribers"
)

type FeedService interface {
	GetRecentPosts(c context.Context, limit, offset string) ([]types.Post, int, error)
	Post(c context.Context, id, token string, post *types.Post) (int, error)
	Subscribe(c context.Context, id, subId, token string) (int, error)
}

func NewFeedService(addr string) FeedService {
	return &feedService{
		addr: addr,
	}
}

type feedService struct {
	addr string
}

func (fs *feedService) GetRecentPosts(
	c context.Context,
	limit, offset string,
) ([]types.Post, int, error) {
	req, err := createRequest(c, http.MethodGet, fs.addr+postsURL, nil)
	if err != nil {
		return nil, fiber.StatusInternalServerError, fmt.Errorf("couldnt process provided data")
	}

	q := req.URL.Query()
	q.Add("limit", limit)
	q.Add("offset", offset)
	q.Add("sort_by", "created_at")
	q.Add("order_by", "desc")
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fiber.StatusInternalServerError, fmt.Errorf("couldnt fetch posts")
	}

	if resp.StatusCode == http.StatusOK {
		posts := []types.Post{}
		err = json.NewDecoder(resp.Body).Decode(&posts)
		if err != nil {
			return nil, fiber.StatusInternalServerError, fmt.Errorf("couldnt verify received posts")
		}
		return posts, resp.StatusCode, nil
	}

	return nil, resp.StatusCode, readResponseError(resp)
}

func (fs *feedService) Post(
	c context.Context,
	id, token string,
	post *types.Post,
) (int, error) {
	req, err := createRequest(c, http.MethodPut, fs.addr+fmt.Sprintf(postURL, id), post)
	if err != nil {
		return fiber.StatusInternalServerError, fmt.Errorf("couldnt process provided data")
	}
	req.Header.Add("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fiber.StatusInternalServerError, fmt.Errorf("couldnt create post")
	}

	if resp.StatusCode == fiber.StatusCreated {
		err = json.NewDecoder(resp.Body).Decode(post)
		if err != nil {
			return fiber.StatusInternalServerError, fmt.Errorf("couldnt verify post creation")
		}
		return resp.StatusCode, nil
	}

	return resp.StatusCode, readResponseError(resp)
}

func (fs *feedService) Subscribe(c context.Context, id, subId, token string) (int, error) {
	req, err := createRequest(
		c,
		http.MethodPut,
		fs.addr+fmt.Sprintf(subscribeURL, id),
		&types.UserBase{Id: subId},
	)
	if err != nil {
		return fiber.StatusInternalServerError, fmt.Errorf("couldnt process provided data")
	}
	req.Header.Add("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fiber.StatusInternalServerError, fmt.Errorf("couldnt subscribe")
	}

	if resp.StatusCode == fiber.StatusOK {
		return resp.StatusCode, nil
	}

	return resp.StatusCode, readResponseError(resp)
}
