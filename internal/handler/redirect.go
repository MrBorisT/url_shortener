package handler

import (
	"net/http"
	"strings"

	"github.com/MrBorisT/url_shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

func Redirect(linksStore *storage.LinksStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortLink := strings.TrimSpace(chi.URLParam(r, "link"))
		originalURL, err := linksStore.GetOriginalURL(r.Context(), shortLink)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
	}
}
