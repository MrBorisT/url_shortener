package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/MrBorisT/url_shortener/internal/helper"
	mw "github.com/MrBorisT/url_shortener/internal/middleware"
	"github.com/MrBorisT/url_shortener/internal/models"
	"github.com/MrBorisT/url_shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

func ListLinks(linksStore *storage.LinksStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := mw.GetUserID(r.Context())
		if !ok {
			log.Println("user ID not found")
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, "something went wrong, try again later")
			return
		}

		encoder := json.NewEncoder(w)

		links, err := linksStore.ListLinks(r.Context(), userID)
		if err != nil {
			log.Println("listing links:", err)
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
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, "something went wrong, try again later")
			return
		}

		encoder := json.NewEncoder(w)

		linkID := strings.TrimSpace(chi.URLParam(r, "id"))
		link, err := linksStore.GetLink(r.Context(), userID, linkID)

		if err != nil {
			log.Println("getting link:", err)
		}

		if err := encoder.Encode(link); err != nil {
			log.Println("encoding links:", err)
		}
	}
}

func CreateLink(linksStore *storage.LinksStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := mw.GetUserID(r.Context())
		if !ok {
			log.Println("user ID not found")
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, "something went wrong, try again later")
			return
		}

		createLinkReq := models.CreateLinkRequest{}

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&createLinkReq); err != nil {
			_ = helper.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		newLink, err := linksStore.CreateLink(r.Context(), userID, createLinkReq)
		if err != nil {
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, "something went wrong, try again later")
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
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, "something went wrong, try again later")
			return
		}

		linkID := strings.TrimSpace(chi.URLParam(r, "linkID"))

		if err := linksStore.DeleteLink(r.Context(), userID, linkID); err != nil {
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, "something went wrong, try again later")
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func DisableLink(linksStore *storage.LinksStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := mw.GetUserID(r.Context())
		if !ok {
			log.Println("user ID not found")
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, "something went wrong, try again later")
			return
		}

		linkID := strings.TrimSpace(chi.URLParam(r, "linkID"))

		if err := linksStore.DisableLink(r.Context(), userID, linkID); err != nil {
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, "something went wrong, try again later")
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

func UpdateLink(linksStore *storage.LinksStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := mw.GetUserID(r.Context())
		if !ok {
			log.Println("user ID not found")
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, "something went wrong, try again later")
			return
		}

		linkID := strings.TrimSpace(chi.URLParam(r, "linkID"))
		updateLinkReq := models.UpdateLinkRequest{}

		if err := json.NewDecoder(r.Body).Decode(&updateLinkReq); err != nil {
			_ = helper.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		newLink, err := linksStore.UpdateLink(r.Context(), userID, linkID, updateLinkReq)
		if err != nil {
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, "something went wrong, try again later")
		}

		if newLink != nil {
			if err := json.NewEncoder(w).Encode(newLink); err != nil {
				log.Println("encoding updated link:", err)
				_ = helper.WriteJSONError(w, http.StatusInternalServerError, "something went wrong, try again later")
			}
		}
	}
}
