package main

import (
	"flag"
	"log"

	"github.com/Alieksieiev0/feed-templ/internal/transport/web"
	"github.com/gofiber/fiber/v2"
)

func main() {
	var (
		webServerAddr = flag.String("web-server", ":3005", "listen address of web server")
	)

	s := web.NewWebServer(fiber.New())
	err := s.Start(*webServerAddr)
	if err != nil {
		log.Fatal(err)
	}
}
