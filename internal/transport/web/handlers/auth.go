package handlers

import (
	"github.com/Alieksieiev0/feed-templ/internal/services"
	"github.com/Alieksieiev0/feed-templ/internal/types"
	"github.com/Alieksieiev0/feed-templ/internal/view/core"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

const (
	title             = "Congratulations! Registration successfully completed. You're all set to explore and engage with our platform."
	link              = "/signin"
	linkTitle         = "Click here to login."
	incorrectCredsErr = "incorrect credentials were provided"
)

func SignupHandler(serv services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := &types.User{}
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).Send([]byte("Error: " + incorrectCredsErr))
		}

		r := serv.Register(c.Context(), user)
		if r.Err != nil {
			return r.SendError(c)
		}

		return Render(
			c,
			core.Success(title, linkTitle, templ.URL(link)),
			templ.WithStatus(r.StatusCode),
		)
	}
}

func SigninHandler(serv services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := &types.User{}
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).Send([]byte("Error: " + incorrectCredsErr))
		}

		userToken, r := serv.Login(c.Context(), user)
		if r.Err != nil {
			return r.SendError(c)
		}

		setUserTokenCookies(c, userToken)
		r.Redirect(c, "/")
		return nil
	}
}

func setUserTokenCookies(c *fiber.Ctx, userToken *types.UserToken) {
	c.Cookie(sessionCookie("jwt", userToken.Token.Value))
	c.Cookie(sessionCookie("id", userToken.User.Id))
	c.Cookie(sessionCookie("username", userToken.User.Username))
	c.Cookie(sessionCookie("email", userToken.User.Email))
}
