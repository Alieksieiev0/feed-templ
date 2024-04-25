package handlers

import (
	"fmt"

	"github.com/Alieksieiev0/feed-templ/internal/services"
	"github.com/Alieksieiev0/feed-templ/internal/view/auth"
	"github.com/Alieksieiev0/feed-templ/internal/view/pages"
	"github.com/Alieksieiev0/feed-templ/internal/view/profile"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func HomePageHandler(feedServ services.FeedService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		posts, r := feedServ.GetPosts(
			c.Context(),
			services.Limit(postsStep),
			services.SortBy("created_at"),
			services.OrderBy("desc"),
		)
		if r.Err != nil {
			return Render(
				c,
				baseWithAuth(c, pages.ServerError("Error: "+r.Err.Error())),
				templ.WithStatus(r.StatusCode),
			)
		}

		setLimitOffsetCookies(c, fmt.Sprint(postsStep*2), fmt.Sprint(postsStep))
		return Render(
			c,
			baseWithAuth(c, pages.Home(isLoggedIn(c), posts)),
			templ.WithStatus(r.StatusCode),
		)
	}
}

func SearchPage(c *fiber.Ctx) error {
	return Render(c, baseWithAuth(c, pages.Search()))
}

func SigninPage(c *fiber.Ctx) error {
	return Render(c, baseWithAuth(c, auth.Signin()))
}

func SignupPage(c *fiber.Ctx) error {
	return Render(c, baseWithAuth(c, auth.Signup()))
}

func ProfilePage(userServ services.UserService, feedServ services.FeedService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			id = c.Cookies("id")
		}

		user, r := userServ.GetById(c.Context(), id)
		if r.Err != nil {
			return r.SendError(c)
		}

		posts, r := feedServ.GetPosts(c.Context(),
			services.Limit(postsStep),
			services.SortBy("created_at"),
			services.OrderBy("desc"),
			services.DefaultParam("user_id", id),
		)
		if r.Err != nil {
			return r.SendError(c)
		}

		return Render(
			c,
			baseWithAuth(c, profile.Card(*user, c.Cookies("id"), posts)),
			templ.WithStatus(fiber.StatusOK),
		)
	}
}
