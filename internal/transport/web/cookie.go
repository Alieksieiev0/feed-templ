package web

import (
	"strconv"

	"github.com/Alieksieiev0/feed-templ/internal/types"
	"github.com/gofiber/fiber/v2"
)

func setUserTokenCookies(c *fiber.Ctx, userToken *types.UserToken) {
	c.Cookie(sessionCookie("jwt", userToken.Token.Value))
	c.Cookie(sessionCookie("id", userToken.User.Id))
	c.Cookie(sessionCookie("username", userToken.User.Username))
	c.Cookie(sessionCookie("email", userToken.User.Email))
}

func setLimitOffsetCookies(c *fiber.Ctx, limit, offset string) {
	c.Cookie(sessionCookie("limit", limit))
	c.Cookie(sessionCookie("offset", offset))
}

func getIntCookie(c *fiber.Ctx, name, def string) (int, error) {
	v, err := strconv.Atoi(c.Cookies(name, def))
	if err != nil {
		return 0, err
	}
	return v, nil
}

func sessionCookie(name, value string) *fiber.Cookie {
	return &fiber.Cookie{
		Name:        name,
		Value:       value,
		SessionOnly: true,
	}
}

func clearCookies(c *fiber.Ctx) {
	c.ClearCookie("jwt", "id", "username", "email")
}
