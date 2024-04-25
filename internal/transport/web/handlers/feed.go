package handlers

import (
	"fmt"
	"strconv"

	"github.com/Alieksieiev0/feed-templ/internal/services"
	"github.com/Alieksieiev0/feed-templ/internal/types"
	"github.com/Alieksieiev0/feed-templ/internal/view/core"
	"github.com/Alieksieiev0/feed-templ/internal/view/feed"
	"github.com/Alieksieiev0/feed-templ/internal/view/search"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

const (
	settingsErr = "bad settings were found while preparing to load posts"
	badDataErr  = "bad data was provided"
	postsStep   = 10
)

func GetPostsHandler(serv services.FeedService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		limit, err := getIntCookie(c, "limit", "10")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).Send([]byte("Error: " + settingsErr))
		}

		offset, err := getIntCookie(c, "offset", "0")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).Send([]byte("Error: " + settingsErr))
		}

		posts, r := serv.GetPosts(
			c.Context(),
			services.Limit(limit),
			services.Offset(offset),
			services.SortBy("created_at"),
			services.OrderBy("desc"),
		)

		if r.Err != nil {
			return r.SendError(c)
		}

		if len(posts) == 0 {
			c.Set("HX-Retarget", "this")
			c.Set("HX-Reswap", "outerHTML")
			return Render(
				c,
				core.Warning("All posts are already loaded"),
				templ.WithStatus(r.StatusCode),
			)
		}
		setLimitOffsetCookies(c, fmt.Sprint(limit+postsStep), fmt.Sprint(offset+postsStep))
		c.Set("HX-Reswap", "beforeend")
		return Render(c, feed.Posts(posts), templ.WithStatus(r.StatusCode))
	}
}

func CreatePostHandler(serv services.FeedService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		post := &types.Post{}
		if err := c.BodyParser(post); err != nil {
			return c.Status(fiber.StatusBadRequest).Send([]byte("Error: " + badDataErr))
		}

		r := serv.Post(c.Context(), c.Cookies("id"), c.Cookies("jwt"), post)
		if r.StatusCode == fiber.StatusUnauthorized {
			redirectToAuth(c, r)
			return nil
		}

		if r.Err != nil {
			return r.SendError(c)
		}

		c.Set("HX-Reswap", "afterbegin")
		return Render(c, feed.Post(*post), templ.WithStatus(r.StatusCode))
	}
}

func SubscribeHandler(serv services.FeedService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		subId := c.Cookies("id")
		jwt := c.Cookies("jwt")
		id := c.Params("id")

		r := serv.Subscribe(c.Context(), id, subId, jwt)
		if r.StatusCode == fiber.StatusUnauthorized {
			redirectToAuth(c, r)
			return nil
		}

		if r.Err != nil {
			return r.SendError(c)
		}

		return Render(c, search.UnsubscribeButton(id), templ.WithStatus(r.StatusCode))
	}
}

func UnsubscribeHandler(serv services.FeedService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		subId := c.Cookies("id")
		jwt := c.Cookies("jwt")
		id := c.Params("id")

		r := serv.Unsubscribe(c.Context(), id, subId, jwt)
		if r.StatusCode == fiber.StatusUnauthorized {
			redirectToAuth(c, r)
			return nil
		}

		if r.Err != nil {
			return r.SendError(c)
		}

		return Render(c, search.SubscribeButton(id), templ.WithStatus(r.StatusCode))
	}
}

func getIntCookie(c *fiber.Ctx, name, def string) (int, error) {
	v, err := strconv.Atoi(c.Cookies(name, def))
	if err != nil {
		return 0, err
	}
	return v, nil
}
