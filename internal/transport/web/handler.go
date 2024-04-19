package web

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"slices"

	"github.com/Alieksieiev0/feed-templ/internal/services"
	"github.com/Alieksieiev0/feed-templ/internal/types"
	"github.com/Alieksieiev0/feed-templ/internal/view/auth"
	"github.com/Alieksieiev0/feed-templ/internal/view/core"
	"github.com/Alieksieiev0/feed-templ/internal/view/feed"
	"github.com/Alieksieiev0/feed-templ/internal/view/notify"
	"github.com/Alieksieiev0/feed-templ/internal/view/pages"
	"github.com/Alieksieiev0/feed-templ/internal/view/search"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
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

func homeHandler(
	feedServ services.FeedService,
	notifServ services.NotificationServices,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		posts, statusCode, err := feedServ.GetRecentPosts(c.Context(), fmt.Sprint(postsStep), "0")

		if err != nil {
			return render(
				c,
				baseWithAuth(c, pages.ServerError("Error: "+err.Error())),
				templ.WithStatus(statusCode),
			)
		}

		id := c.Cookies("id")
		notifications := []types.Notification{}
		if id != "" {
			notifications, statusCode, err = notifServ.Get(c.Context(), id)
			if err != nil {
				return render(
					c,
					baseWithAuth(c, pages.ServerError("Error: "+err.Error())),
					templ.WithStatus(statusCode),
				)
			}
		}

		setLimitOffsetCookies(c, fmt.Sprint(postsStep*2), fmt.Sprint(postsStep))
		return render(
			c,
			baseWithAuth(c, pages.Home(isLoggedIn(c), posts, notifications)),
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

		fmt.Println("test")
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
		id := c.Cookies("id")
		username := c.FormValue("username")
		if username == "" {
			return render(c, search.Results([]types.User{}, id), templ.WithStatus(fiber.StatusOK))
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

		index := indexById(id, users)
		if index >= 0 {
			users = slices.Delete(users, index, index+1)
		}

		return render(c, search.Results(users, id), templ.WithStatus(statusCode))
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

func subscribeHandler(serv services.FeedService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		subId := c.Cookies("id")
		id := c.Params("id")
		jwt := c.Cookies("jwt")

		statusCode, err := serv.Subscribe(c.Context(), id, subId, jwt)
		if statusCode == fiber.StatusUnauthorized {
			clearCookies(c)
			redirect(c, "/signin", statusCode)
			return nil
		}

		if err != nil {
			return c.Status(statusCode).Send([]byte("Error: " + err.Error()))
		}

		return render(c, search.UnsubscribeButton(subId), templ.WithStatus(statusCode))
	}
}

func getNotificationsHandler(serv services.NotificationServices) fiber.Handler {
	return func(c *fiber.Ctx) error {
		notifications, statusCode, err := serv.Get(c.Context(), c.Cookies("id"))
		if err != nil {
			return c.Status(statusCode).Send([]byte("Error: " + err.Error()))
		}

		return render(c, notify.Notifications(notifications))
	}
}

func listenHandler(serv services.NotificationServices) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")

		ch := make(chan *types.Notification)
		err := serv.Listen(c.Context(), c.Cookies("id"), ch)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			for {
				notification := <-ch
				err := json.NewEncoder(w).Encode(notification)
				if err != nil {
					log.Println(err)
					break
				}

				err = w.Flush()
				if err != nil {
					log.Println(err)
					break
				}
			}
		}))

		return nil
	}
}
