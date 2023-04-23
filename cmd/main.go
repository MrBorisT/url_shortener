package main

import (
	"log"

	"github.com/MrBorisT/go_url_shortener/internal/config"
	"github.com/MrBorisT/go_url_shortener/internal/domain"
	"github.com/MrBorisT/go_url_shortener/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	if err := config.Init(); err != nil {
		log.Fatalln(err)
	}

	repo := repository.NewInMemoryRepository()
	urlShortenerService := domain.NewService(repo, config.ConfigData.MaxShortUrlLen)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil)
	})

	app.Get("/:url", func(c *fiber.Ctx) error {
		shortURL := c.Params("url")
		fullURL, err := urlShortenerService.GetFullURL(shortURL)
		if err != nil {
			return err
		}

		log.Println("redirecting to", fullURL)
		return c.Redirect(fullURL)
	})

	app.Post("/shorten", func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		fullURL := ""
		if err := c.BodyParser(&fullURL); err != nil {
			return err
		}

		shortURL, err := urlShortenerService.ShortenURL(fullURL)
		if err != nil {
			return err
		}

		return c.JSON(shortURL)
	})

	// use port with ":"
	// E.g.: ":3000"
	log.Fatal(app.Listen(config.ConfigData.Port))
}
