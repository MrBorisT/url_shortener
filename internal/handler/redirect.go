package handler

import (
	"errors"
	"log"
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
			if errors.Is(err, storage.ErrLinkNotFound) {
				http.NotFound(w, r)
				return
			}

			log.Println("getting original URL:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, normalizeURL(originalURL), http.StatusFound)
	}
}

func normalizeURL(raw string) string {
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return raw
	}

	return "https://" + raw
}
