package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type DomainRouter interface {
	SubscribeRoutes(app *fiber.Router)
}

type Server struct {
	listenAddr string
	app        *fiber.App
}

func NewServer(port string) *Server {
	return &Server{
		listenAddr: port,
		app:        fiber.New(),
	}
}

func (s *Server) Ignite(routers ...DomainRouter) {
	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	s.app.Use(cors.New(
		cors.Config{
			AllowOrigins: "*",
		}))

	s.app.Use(logger.New())

	custom_router := s.app.Group("/")

	for _, router := range routers {
		router.SubscribeRoutes(&custom_router)
	}

	log.Fatal(s.app.Listen(":" + s.listenAddr))
}
