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
	registerURL = "/api/auth/register"
	loginURL    = "/api/auth/login"
)

type AuthService interface {
	Register(c context.Context, user *types.User) (int, error)
	Login(c context.Context, user *types.User) (*types.UserToken, int, error)
}

func NewAuthService(addr string) AuthService {
	return &authService{
		addr: addr,
	}
}

type authService struct {
	addr string
}

func (as *authService) Register(c context.Context, user *types.User) (int, error) {
	req, err := createRequest(c, http.MethodPost, as.addr+registerURL, user)
	if err != nil {
		return fiber.StatusInternalServerError, fmt.Errorf("couldnt process provided credentials")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fiber.StatusInternalServerError, fmt.Errorf(
			"couldnt register user with provided credentials",
		)
	}

	if resp.StatusCode == fiber.StatusCreated {
		return resp.StatusCode, nil
	}

	return resp.StatusCode, readResponseError(resp)
}

func (as *authService) Login(c context.Context, user *types.User) (*types.UserToken, int, error) {
	req, err := createRequest(c, http.MethodPost, as.addr+loginURL, user)
	if err != nil {
		return nil, fiber.StatusInternalServerError, fmt.Errorf(
			"couldn`t process provided credentials",
		)
	}

	resp, err := http.DefaultClient.Do(req)
	fmt.Println(err)
	if err != nil {
		return nil, fiber.StatusInternalServerError, fmt.Errorf(
			"couldnt login user with provided credentials",
		)
	}

	if resp.StatusCode == fiber.StatusOK {
		userToken := &types.UserToken{}
		err = json.NewDecoder(resp.Body).Decode(userToken)
		if err != nil {
			return nil, fiber.StatusInternalServerError, fmt.Errorf("couldnt verify credentials")
		}
		return userToken, resp.StatusCode, nil
	}

	return nil, resp.StatusCode, readResponseError(resp)
}
