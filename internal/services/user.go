package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Alieksieiev0/feed-templ/internal/types"
	"github.com/gofiber/fiber/v2"
)

var (
	usersUrl = "/api/feed/users"
)

type UserService interface {
	Search(c context.Context, query, limit, offset string) ([]types.User, int, error)
}

func NewUserService(addr string) UserService {
	return &userService{
		addr: addr,
	}
}

type userService struct {
	addr string
}

func (us *userService) Search(
	c context.Context,
	query, limit, offset string,
) ([]types.User, int, error) {
	req, err := createRequest(c, http.MethodGet, us.addr+usersUrl, nil)
	if err != nil {
		return nil, fiber.StatusInternalServerError, fmt.Errorf("couldnt process provided data")
	}

	q := req.URL.Query()
	q.Add("limit", limit)
	q.Add("offset", offset)
	q.Add("username", query)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fiber.StatusInternalServerError, fmt.Errorf("couldnt fetch posts")
	}

	if resp.StatusCode == http.StatusOK {
		users := []types.User{}
		err = json.NewDecoder(resp.Body).Decode(&users)
		if err != nil {
			return nil, fiber.StatusInternalServerError, fmt.Errorf("couldnt verify received posts")
		}
		return users, resp.StatusCode, nil
	}

	return nil, resp.StatusCode, readResponseError(resp)
}
