package services

import (
	"context"
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
	GetRecentPosts(c context.Context, params ...Param) ([]types.Post, *Response)
	Post(c context.Context, id, token string, post *types.Post) *Response
	Subscribe(c context.Context, id, subId, token string) *Response
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
	params ...Param,
) ([]types.Post, *Response) {
	req, err := createRequest(c, http.MethodGet, fs.addr+postsURL, nil)
	if err != nil {
		return nil, NewResponse(fiber.StatusInternalServerError, err)
	}

	params = append(params, SortBy("created_at"), OrderBy("desc"))
	updateQuery(req, params)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, NewResponse(fiber.StatusInternalServerError, err)
	}

	if resp.StatusCode == http.StatusOK {
		posts := []types.Post{}
		return posts, parseResponse(resp, &posts)
	}

	return nil, NewResponse(resp.StatusCode, readResponseError(resp))
}

func (fs *feedService) Post(
	c context.Context,
	id, token string,
	post *types.Post,
) *Response {
	req, err := createRequest(c, http.MethodPut, fs.addr+fmt.Sprintf(postURL, id), post)
	if err != nil {
		return NewResponse(fiber.StatusInternalServerError, err)
	}
	req.Header.Add("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return NewResponse(fiber.StatusInternalServerError, err)
	}

	if resp.StatusCode == fiber.StatusCreated {
		return parseResponse(resp, post)
	}

	return NewResponse(resp.StatusCode, readResponseError(resp))
}

func (fs *feedService) Subscribe(c context.Context, id, subId, token string) *Response {
	req, err := createRequest(
		c,
		http.MethodPut,
		fs.addr+fmt.Sprintf(subscribeURL, id),
		&types.UserBase{Id: subId},
	)
	if err != nil {
		return NewResponse(fiber.StatusInternalServerError, err)
	}
	req.Header.Add("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return NewResponse(fiber.StatusInternalServerError, err)
	}

	if resp.StatusCode == fiber.StatusOK {
		return NewResponse(resp.StatusCode, nil)
	}

	return NewResponse(resp.StatusCode, readResponseError(resp))
}
