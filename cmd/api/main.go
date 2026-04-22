package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := newRouter()

	log.Println("started server on port", ":8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalln(err)
	}
}

func newRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/health", nil)

	r.Route("/links", func(r chi.Router) {
		r.Post("/", nil)
		r.Get("/", nil)
		r.Get("/{id}", nil)
		r.Patch("/{id}", nil)
		r.Delete("/{id}", nil)
		r.Post("/{id}/disable", nil)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", nil)
		r.Post("/login", nil)
	})

	return r
}
