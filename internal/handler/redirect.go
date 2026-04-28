package handler

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/MrBorisT/url_shortener/internal/linkerr"
	"github.com/MrBorisT/url_shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

func Redirect(linksService *service.LinkService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortLink := strings.TrimSpace(chi.URLParam(r, "link"))
		originalURL, err := linksService.GetOriginalURL(r.Context(), shortLink)
		if err != nil {
			switch {
			case errors.Is(err, linkerr.ErrLinkNotFound):
				http.NotFound(w, r)
			case errors.Is(err, linkerr.ErrLinkDisabled):
				http.Error(w, http.StatusText(http.StatusGone), http.StatusGone)
			default:
				log.Println("getting original URL:", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}
		if err := linksService.IncrementClickCount(r.Context(), shortLink); err != nil {
			if errors.Is(err, linkerr.ErrLinkNotFound) {
				http.NotFound(w, r)
				return
			}
			log.Println("incrementing click count:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, originalURL, http.StatusFound)
	}
}
