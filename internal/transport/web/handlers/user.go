package handlers

import (
	"slices"

	"github.com/Alieksieiev0/feed-templ/internal/services"
	"github.com/Alieksieiev0/feed-templ/internal/types"
	"github.com/Alieksieiev0/feed-templ/internal/view/search"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

const (
	usersStep = 10
)

func GetUsersHandler(serv services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Cookies("id")
		username := c.FormValue("username")
		if username == "" {
			return Render(c, search.Results([]types.User{}, id), templ.WithStatus(fiber.StatusOK))
		}

		users, r := serv.Search(
			c.Context(),
			services.DefaultParam("username", username),
			services.Limit(usersStep),
		)

		if r.Err != nil {
			return r.SendError(c)
		}

		index := indexById(id, users)
		if index >= 0 {
			users = slices.Delete(users, index, index+1)
		}

		return Render(c, search.Results(users, id), templ.WithStatus(r.StatusCode))
	}
}

func indexById(id string, users []types.User) int {
	index := -1
	if id == "" {
		return index
	}

	for i, u := range users {
		if u.Id == id {
			index = i
			break
		}
	}

	return index
}
