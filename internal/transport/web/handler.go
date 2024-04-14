package web

import (
	"fmt"

	"github.com/Alieksieiev0/feed-templ/internal/services"
	"github.com/Alieksieiev0/feed-templ/internal/types"
	"github.com/Alieksieiev0/feed-templ/internal/view/auth"
	"github.com/Alieksieiev0/feed-templ/internal/view/core"
	"github.com/Alieksieiev0/feed-templ/internal/view/feed"
	"github.com/Alieksieiev0/feed-templ/internal/view/pages"
	"github.com/Alieksieiev0/feed-templ/internal/view/search"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

const (
	title             = "Congratulations! Registration successfully completed. You're all set to explore and engage with our platform."
	link              = "/signin"
	linkTitle         = "Click here to login."
	settingsErr       = "bad settings were found while preparing to load posts"
	incorrectCredsErr = "incorrect credentials were provided"
	badDataErr        = "bad data was provided"
	postsStep         = 10
)

type Pagination struct {
	Limit  int
	Offset int
}

func homeHandler(serv services.FeedService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		posts, statusCode, err := serv.GetRecentPosts(c.Context(), fmt.Sprint(postsStep), "0")

		if err != nil {
			return render(
				c,
				baseWithAuth(c, pages.ServerError("Error: "+err.Error())),
				templ.WithStatus(statusCode),
			)
		}

		setLimitOffsetCookies(c, fmt.Sprint(postsStep*2), fmt.Sprint(postsStep))
		return render(
			c,
			baseWithAuth(c, pages.Home(isLoggedIn(c), posts)),
			templ.WithStatus(statusCode),
		)
	}
}

func searchPage(c *fiber.Ctx) error {
	return render(c, baseWithAuth(c, pages.Search()))
}

func signinPage(c *fiber.Ctx) error {
	return render(c, baseWithAuth(c, auth.Signin()))
}

func signupPage(c *fiber.Ctx) error {
	return render(c, baseWithAuth(c, auth.Signup()))
}

func signupHandler(serv services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := &types.User{}
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).Send([]byte("Error: " + incorrectCredsErr))
		}

		statusCode, err := serv.Register(c.Context(), user)
		if err != nil {
			return c.Status(statusCode).Send([]byte("Error: " + err.Error()))
		}

		return render(
			c,
			core.Success(title, linkTitle, templ.URL(link)),
			templ.WithStatus(statusCode),
		)
	}
}

func signinHandler(serv services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := &types.User{}
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).Send([]byte("Error: " + incorrectCredsErr))
		}

		userToken, statusCode, err := serv.Login(c.Context(), user)
		if err != nil {
			return c.Status(statusCode).Send([]byte("Error: " + err.Error()))
		}

		setUserTokenCookies(c, userToken)
		redirect(c, "/", statusCode)
		return nil
	}
}

func createPostHandler(serv services.FeedService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		post := &types.Post{}
		if err := c.BodyParser(post); err != nil {
			return c.Status(fiber.StatusBadRequest).Send([]byte("Error: " + badDataErr))
		}

		statusCode, err := serv.Post(c.Context(), c.Cookies("id"), c.Cookies("jwt"), post)
		if statusCode == fiber.StatusUnauthorized {
			clearCookies(c)
			redirect(c, "/signin", statusCode)
			return nil
		}

		if err != nil {
			return c.Status(statusCode).Send([]byte("Error: " + err.Error()))
		}

		c.Set("HX-Reswap", "afterbegin")
		return render(c, feed.Post(*post), templ.WithStatus(statusCode))
	}
}

func getPostsHandler(serv services.FeedService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		limit, err := getIntCookie(c, "limit", "10")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).Send([]byte("Error: " + settingsErr))
		}

		offset, err := getIntCookie(c, "offset", "0")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).Send([]byte("Error: " + settingsErr))
		}

		posts, statusCode, err := serv.GetRecentPosts(
			c.Context(),
			fmt.Sprint(limit),
			fmt.Sprint(offset),
		)

		if err != nil {
			return c.Status(statusCode).Send([]byte("Error: " + err.Error()))
		}

		if len(posts) == 0 {
			c.Set("HX-Retarget", "this")
			c.Set("HX-Reswap", "outerHTML")
			return render(
				c,
				core.Warning("All posts are already loaded"),
				templ.WithStatus(statusCode),
			)
		}
		setLimitOffsetCookies(c, fmt.Sprint(limit+postsStep), fmt.Sprint(offset+postsStep))
		c.Set("HX-Reswap", "beforeend")
		return render(c, feed.Posts(posts), templ.WithStatus(statusCode))
	}
}

func getUsersHandler(serv services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.FormValue("username")
		if username == "" {
			return render(c, search.Results([]types.User{}), templ.WithStatus(fiber.StatusOK))
		}

		users, statusCode, err := serv.Search(
			c.Context(),
			username,
			"10",
			"0",
		)

		if err != nil {
			return c.Status(statusCode).Send([]byte("Error: " + err.Error()))
		}

		return render(c, search.Results(users), templ.WithStatus(statusCode))
	}
}
