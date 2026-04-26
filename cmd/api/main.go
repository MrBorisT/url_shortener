package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/MrBorisT/url_shortener/internal/config"
	"github.com/MrBorisT/url_shortener/internal/handler"
	auth "github.com/MrBorisT/url_shortener/internal/jwt"
	mw "github.com/MrBorisT/url_shortener/internal/middleware"
	"github.com/MrBorisT/url_shortener/internal/service"
	"github.com/MrBorisT/url_shortener/internal/storage"
)

func main() {
	env := os.Getenv("APP_ENV")

	switch env {
	case "production":
		// nothing to load
	default:
		_ = godotenv.Load()
	}

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Error loading configuration:", err)
	}
	pool, err := newPool(config)
	if err != nil {
		log.Fatalln("Unable to create database pool:", err)
	}

	defer pool.Close()
	userStore := storage.NewUserStore(pool)
	authManager := auth.NewJWTManager(config)
	linksStore := storage.NewLinksStore(pool)
	linkService := service.NewLinkService(linksStore)
	authService := service.NewAuthService(userStore, authManager)

	r := newRouter(linkService, authService)

	log.Println("started server on port", config.Port)
	if err := http.ListenAndServe(config.Port, r); err != nil {
		log.Fatalln(err)
	}
}

func compileDSN(config *config.Config) string {
	return "host=" + config.DBHost +
		" port=" + config.DBPort +
		" dbname=" + config.DBName +
		" user=" + config.DBUser +
		" password=" + config.DBPassword +
		" sslmode=" + config.DBSSLMode
}

func newPool(config *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dsn := compileDSN(config)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

func newRouter(
	linkService *service.LinkService,
	authService *service.AuthService) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", handler.Health)

	r.Route("/api", func(r chi.Router) {
		r.Use(mw.JSONMiddleware)

		r.Route("/links", func(r chi.Router) {
			r.Use(mw.AuthMiddleware(authService))
			r.Post("/", handler.CreateLink(linkService))
			r.Get("/", handler.ListLinks(linkService))
			r.Get("/{id}", handler.GetLink(linkService))
			r.Patch("/{id}", handler.UpdateLink(linkService))
			r.Delete("/{id}", handler.DeleteLink(linkService))
			r.Post("/{id}/disable", handler.DisableLink(linkService))
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", handler.Register(authService))
			r.Post("/login", handler.Login(authService))
		})
	})

	r.Get("/{link}", handler.Redirect(linkService))

	return r
}
