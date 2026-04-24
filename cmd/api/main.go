package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MrBorisT/url_shortener/internal/config"
	"github.com/MrBorisT/url_shortener/internal/handler"
	auth "github.com/MrBorisT/url_shortener/internal/jwt"
	mw "github.com/MrBorisT/url_shortener/internal/middleware"
	"github.com/MrBorisT/url_shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
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

	r := newRouter(userStore, authManager)

	log.Println("started server on port", ":8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
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

func newRouter(userStore *storage.UserStore, authManager *auth.JWTManager) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", handler.Health)

	r.Route("/api", func(r chi.Router) {
		r.Use(mw.JSONMiddleware)

		r.Route("/links", func(r chi.Router) {
			r.Use(mw.AuthMiddleware(authManager))
			r.Post("/", handler.CreateLink)
			r.Get("/", handler.ListLinks)
			r.Get("/{id}", handler.GetLink)
			r.Patch("/{id}", handler.UpdateLink)
			r.Delete("/{id}", handler.DeleteLink)
			r.Post("/{id}/disable", handler.DisableLink)
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", handler.Register(userStore))
			r.Post("/login", handler.Login(userStore, authManager))
		})
	})

	r.Get("/{link}", handler.Redirect)

	return r
}
