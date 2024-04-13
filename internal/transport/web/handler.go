package web

import (
	"github.com/Alieksieiev0/feed-templ/internal/services"
	"github.com/Alieksieiev0/feed-templ/internal/types"
	"github.com/Alieksieiev0/feed-templ/internal/view/auth"
	"github.com/Alieksieiev0/feed-templ/internal/view/core"
	"github.com/Alieksieiev0/feed-templ/internal/view/layout"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

const (
	title     = "Congratulations! Registration successfully completed. You're all set to explore and engage with our platform."
	link      = "/signin"
	linkTitle = "Click here to login."
)

func homePage(c *fiber.Ctx) error {
	return render(c, baseWithAuth(c, core.Home(isLoggedIn(c), []types.Post{})))
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
			return c.Status(fiber.StatusBadRequest).
				Send([]byte("Error: incorrect credentials were provided"))
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
			return c.Status(fiber.StatusBadRequest).
				Send([]byte("Error: incorrect credentials were provided"))
		}

		userToken, statusCode, err := serv.Login(c.Context(), user)
		if err != nil {
			return c.Status(statusCode).Send([]byte("Error: " + err.Error()))
		}

		setCookies(c, userToken)
		c.Set("HX-Redirect", "/")
		c.Status(statusCode)

		return nil
	}
}

func postsHandler(serv services.FeedService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		post := &types.Post{}
		if err := c.BodyParser(post); err != nil {
			return c.Status(fiber.StatusBadRequest).
				Send([]byte("Error: bad data was provided"))
		}

		token := c.Cookies("jwt")
		if token == "" {
			redirect(c, "/signin", fiber.StatusBadRequest)
			return nil
		}

		statusCode, err := serv.Post(c.Context(), token, post)
		if statusCode == fiber.StatusUnauthorized {
			redirect(c, "/signin", statusCode)
			return nil
		}

		if err != nil {
			return c.Status(statusCode).Send([]byte("Error: " + err.Error()))
		}

		c.Set("HX-Reswap", "beforebegin")
		return render(c, core.Post(*post), templ.WithStatus(statusCode))
	}
}

func baseWithAuth(c *fiber.Ctx, contents templ.Component) templ.Component {
	return layout.Base(isLoggedIn(c), contents)
}

func isLoggedIn(c *fiber.Ctx) bool {
	return c.Cookies("jwt") != ""
}

func redirect(c *fiber.Ctx, url string, statusCode int) {
	c.Set("HX-Redirect", "/signin")
	c.Status(fiber.StatusBadRequest)
}

func setCookies(c *fiber.Ctx, userToken *types.UserToken) {
	tokenCookie := &fiber.Cookie{
		Name:        "jwt",
		Value:       userToken.Token.Value,
		SessionOnly: true,
	}
	c.Cookie(tokenCookie)

	idCookie := &fiber.Cookie{
		Name:        "id",
		Value:       userToken.User.Id,
		SessionOnly: true,
	}
	c.Cookie(idCookie)

	usernameCookie := &fiber.Cookie{
		Name:        "username",
		Value:       userToken.User.Username,
		SessionOnly: true,
	}
	c.Cookie(usernameCookie)

	emailCookie := &fiber.Cookie{
		Name:        "email",
		Value:       userToken.User.Username,
		SessionOnly: true,
	}
	c.Cookie(emailCookie)
}

func render(
	c *fiber.Ctx,
	component templ.Component,
	options ...func(*templ.ComponentHandler),
) error {
	componentHandler := templ.Handler(component)
	for _, o := range options {
		o(componentHandler)
	}
	return adaptor.HTTPHandler(componentHandler)(c)
}
