package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Alieksieiev0/feed-templ/internal/types"
	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	Register(c context.Context, user *types.User) (int, error)
	Login(c context.Context, user *types.User) (*types.UserToken, int, error)
}

func NewAuthService(addr, registerURL, loginURL string) AuthService {
	return &authService{
		addr:        addr,
		registerURL: registerURL,
		loginURL:    loginURL,
	}
}

type authService struct {
	addr        string
	registerURL string
	loginURL    string
}

func (as *authService) Register(c context.Context, user *types.User) (int, error) {
	resp, err := SendRequest(c, as.addr+as.registerURL, user)
	if err != nil {
		return fiber.StatusInternalServerError, fmt.Errorf("couldn`t process provided credentials")
	}

	if resp.StatusCode == fiber.StatusCreated {
		return resp.StatusCode, nil
	}

	return resp.StatusCode, ReadResponseError(resp)
}

func (as *authService) Login(c context.Context, user *types.User) (*types.UserToken, int, error) {
	resp, err := SendRequest(c, as.addr+as.loginURL, user)
	if err != nil {
		return nil, fiber.StatusInternalServerError, fmt.Errorf(
			"couldn`t process provided credentials",
		)
	}

	if resp.StatusCode == fiber.StatusOK {
		userToken := &types.UserToken{}
		err = json.NewDecoder(resp.Body).Decode(userToken)
		if err != nil {
			return nil, fiber.StatusInternalServerError, fmt.Errorf("couldn`t verify credentials")
		}
		return userToken, resp.StatusCode, nil
	}

	return nil, resp.StatusCode, ReadResponseError(resp)
}
