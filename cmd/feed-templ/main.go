package main

import (
	"flag"
	"log"

	"github.com/Alieksieiev0/feed-templ/internal/services"
	"github.com/Alieksieiev0/feed-templ/internal/transport/web"
	"github.com/gofiber/fiber/v2"
)

func main() {
	var (
		webServerAddr = flag.String("web-server", ":3005", "listen address of web server")
		apiAddr       = flag.String(
			"api",
			"http://localhost:8080",
			"api address",
		)
	)

	s := web.NewWebServer(fiber.New(), *webServerAddr)
	authServ := services.NewAuthService(*apiAddr)
	feedServ := services.NewFeedService(*apiAddr)
	userServ := services.NewUserService(*apiAddr)
	err := s.Start(authServ, feedServ, userServ)
	if err != nil {
		log.Fatal(err)
	}
}
