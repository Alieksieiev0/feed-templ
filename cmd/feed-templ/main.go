package main

import (
	"flag"
	"log"

	"github.com/Alieksieiev0/feed-templ/internal/services"
	"github.com/Alieksieiev0/feed-templ/internal/transport/web"
	"github.com/gofiber/fiber/v2"
)

const (
	registerURL = "/api/auth/register"
	loginURL    = "/api/auth/login"
)

func main() {
	var (
		webServerAddr = flag.String("web-server", ":3005", "listen address of web server")
		apiAddr       = flag.String(
			"auth-service",
			"http://localhost:8080",
			"address of auth service",
		)
	)

	s := web.NewWebServer(fiber.New())
	authServ := services.NewAuthService(*apiAddr, registerURL, loginURL)
	err := s.Start(*webServerAddr, authServ)
	if err != nil {
		log.Fatal(err)
	}
}
