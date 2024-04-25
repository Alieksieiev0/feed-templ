package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Alieksieiev0/feed-templ/internal/types"
	"github.com/gofiber/fiber/v2"
)

var (
	usersUrl = "/api/feed/users"
	userUrl  = "/api/feed/users/%s"
)

type UserService interface {
	Search(c context.Context, params ...Param) ([]types.User, *Response)
	GetById(c context.Context, id string) (*types.User, *Response)
}

func NewUserService(addr string) UserService {
	return &userService{
		addr: addr,
	}
}

type userService struct {
	addr string
}

func (us *userService) Search(c context.Context, params ...Param) ([]types.User, *Response) {
	req, err := createRequest(c, http.MethodGet, us.addr+usersUrl, nil)
	if err != nil {
		return nil, NewResponse(fiber.StatusInternalServerError, err)
	}

	updateQuery(req, params)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, NewResponse(fiber.StatusInternalServerError, err)
	}

	if resp.StatusCode == http.StatusOK {
		users := []types.User{}
		return users, parseResponse(resp, &users)
	}

	return nil, NewResponse(resp.StatusCode, readResponseError(resp))
}

func (us *userService) GetById(c context.Context, id string) (*types.User, *Response) {
	req, err := createRequest(c, http.MethodGet, us.addr+fmt.Sprintf(userUrl, id), nil)
	if err != nil {
		return nil, NewResponse(fiber.StatusInternalServerError, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, NewResponse(fiber.StatusInternalServerError, err)
	}

	if resp.StatusCode == http.StatusOK {
		user := &types.User{}
		return user, parseResponse(resp, user)
	}

	return nil, NewResponse(resp.StatusCode, readResponseError(resp))
}
