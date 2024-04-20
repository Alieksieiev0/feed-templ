package services

import (
	"context"
	"net/http"

	"github.com/Alieksieiev0/feed-templ/internal/types"
	"github.com/gofiber/fiber/v2"
)

const (
	registerURL = "/api/auth/register"
	loginURL    = "/api/auth/login"
)

type AuthService interface {
	Register(c context.Context, user *types.User) *Response
	Login(c context.Context, user *types.User) (*types.UserToken, *Response)
}

func NewAuthService(addr string) AuthService {
	return &authService{
		addr: addr,
	}
}

type authService struct {
	addr string
}

func (as *authService) Register(c context.Context, user *types.User) *Response {
	req, err := createRequest(c, http.MethodPost, as.addr+registerURL, user)
	if err != nil {
		return NewResponse(fiber.StatusInternalServerError, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return NewResponse(fiber.StatusInternalServerError, err)
	}

	if resp.StatusCode == fiber.StatusCreated {
		return NewResponse(resp.StatusCode, nil)
	}

	return NewResponse(resp.StatusCode, readResponseError(resp))
}

func (as *authService) Login(c context.Context, user *types.User) (*types.UserToken, *Response) {
	req, err := createRequest(c, http.MethodPost, as.addr+loginURL, user)
	if err != nil {
		return nil, NewResponse(fiber.StatusInternalServerError, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, NewResponse(fiber.StatusInternalServerError, err)
	}

	if resp.StatusCode == fiber.StatusOK {
		userToken := &types.UserToken{}
		return userToken, parseResponse(resp, userToken)
	}

	return nil, NewResponse(resp.StatusCode, readResponseError(resp))
}
