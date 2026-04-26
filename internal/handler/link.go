package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/MrBorisT/url_shortener/internal/helper"
	mw "github.com/MrBorisT/url_shortener/internal/middleware"
	"github.com/MrBorisT/url_shortener/internal/models"
	"github.com/MrBorisT/url_shortener/internal/service"
	"github.com/MrBorisT/url_shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

func ListLinks(linksStore *storage.LinksStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := mw.GetUserID(r.Context())
		if !ok {
			log.Println("user ID not found")
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, helper.ErrInternal)
			return
		}

		encoder := json.NewEncoder(w)

		links, err := linksStore.ListLinks(r.Context(), userID)
		if err != nil {
			log.Println("listing links:", err)
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, helper.ErrInternal)
			return
		}

		if err := encoder.Encode(links); err != nil {
			log.Println("encoding links:", err)
		}
	}
}

func GetLink(linksStore *storage.LinksStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := mw.GetUserID(r.Context())
		if !ok {
			log.Println("user ID not found")
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, helper.ErrInternal)
			return
		}

		encoder := json.NewEncoder(w)

		linkID := strings.TrimSpace(chi.URLParam(r, "id"))
		link, err := linksStore.GetLink(r.Context(), userID, linkID)

		if errors.Is(err, storage.ErrLinkNotFound) {
			_ = helper.WriteJSONError(w, http.StatusNotFound, "link not found")
			return
		}

		if err != nil {
			log.Println("getting link:", err)
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, helper.ErrInternal)
			return
		}

		if err := encoder.Encode(link); err != nil {
			log.Println("encoding links:", err)
		}
	}
}

func CreateLink(linkService *service.LinkService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := mw.GetUserID(r.Context())
		if !ok {
			log.Println("user ID not found")
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, helper.ErrInternal)
			return
		}

		createLinkReq := models.CreateLinkRequest{}

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&createLinkReq); err != nil {
			_ = helper.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		newLink, err := linkService.CreateLink(r.Context(), userID, createLinkReq)
		if err != nil {
			if errors.Is(err, service.ErrInvalidOriginalURL) {
				_ = helper.WriteJSONError(w, http.StatusBadRequest, "original_url is required")
				return
			}

			log.Println("creating link:", err)
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, helper.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(newLink); err != nil {
			log.Println("creating new link:", err)
		}
	}
}

func DeleteLink(linksStore *storage.LinksStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := mw.GetUserID(r.Context())
		if !ok {
			log.Println("user ID not found")
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, helper.ErrInternal)
			return
		}

		linkID := strings.TrimSpace(chi.URLParam(r, "id"))

		if err := linksStore.DeleteLink(r.Context(), userID, linkID); err != nil {
			switch {
			case errors.Is(err, storage.ErrEmptyOriginalURL):
				_ = helper.WriteJSONError(w, http.StatusBadRequest, "original_url is required")
			case errors.Is(err, storage.ErrLinkNotFound):
				_ = helper.WriteJSONError(w, http.StatusNotFound, "link not found")
			default:
				_ = helper.WriteJSONError(w, http.StatusInternalServerError, helper.ErrInternal)
			}
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func DisableLink(linksStore *storage.LinksStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := mw.GetUserID(r.Context())
		if !ok {
			log.Println("user ID not found")
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, helper.ErrInternal)
			return
		}

		linkID := strings.TrimSpace(chi.URLParam(r, "id"))

		if err := linksStore.DisableLink(r.Context(), userID, linkID); err != nil {
			switch {
			case errors.Is(err, storage.ErrEmptyOriginalURL):
				_ = helper.WriteJSONError(w, http.StatusBadRequest, "original_url is required")
			case errors.Is(err, storage.ErrLinkNotFound):
				_ = helper.WriteJSONError(w, http.StatusNotFound, "link not found")
			default:
				log.Println("disabling link:", err)
				_ = helper.WriteJSONError(w, http.StatusInternalServerError, helper.ErrInternal)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

func UpdateLink(linksStore *storage.LinksStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := mw.GetUserID(r.Context())
		if !ok {
			log.Println("user ID not found")
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, helper.ErrInternal)
			return
		}

		linkID := strings.TrimSpace(chi.URLParam(r, "id"))
		updateLinkReq := models.UpdateLinkRequest{}

		if err := json.NewDecoder(r.Body).Decode(&updateLinkReq); err != nil {
			_ = helper.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		newLink, err := linksStore.UpdateLink(r.Context(), userID, linkID, updateLinkReq)
		if err != nil {
			switch {
			case errors.Is(err, storage.ErrEmptyOriginalURL):
				_ = helper.WriteJSONError(w, http.StatusBadRequest, "original_url is required")
			case errors.Is(err, storage.ErrLinkNotFound):
				_ = helper.WriteJSONError(w, http.StatusNotFound, "link not found")
			default:
				log.Println("updating link:", err)
				_ = helper.WriteJSONError(w, http.StatusInternalServerError, helper.ErrInternal)
			}
			return
		}

		if newLink != nil {
			if err := json.NewEncoder(w).Encode(newLink); err != nil {
				log.Println("encoding updated link:", err)
				_ = helper.WriteJSONError(w, http.StatusInternalServerError, helper.ErrInternal)
			}
		}
	}
}
