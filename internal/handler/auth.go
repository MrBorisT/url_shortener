package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/MrBorisT/url_shortener/internal/autherr"
	"github.com/MrBorisT/url_shortener/internal/helper"
	"github.com/MrBorisT/url_shortener/internal/models"
	"github.com/MrBorisT/url_shortener/internal/service"
)

func Register(userService *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userRequest := models.UserRequest{}

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&userRequest); err != nil {
			log.Println("decoding user:", err)
			_ = helper.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if err := userService.RegisterUser(r.Context(), userRequest); err != nil {
			switch err {
			case autherr.ErrEmptyUserEmail, autherr.ErrIncorrectUserEmail, autherr.ErrEmptyUserPassword, autherr.ErrShortUserPassword, autherr.ErrLongUserPassword:
				_ = helper.WriteJSONError(w, http.StatusBadRequest, err.Error())
				return
			case autherr.ErrUserAlreadyExists:
				_ = helper.WriteJSONError(w, http.StatusConflict, "user with this email already exists")
			default:
				log.Println("registering user:", err)
				_ = helper.WriteJSONError(w, http.StatusInternalServerError, helper.ErrInternal)
			}
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func Login(userService *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userRequest := models.UserRequest{}

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&userRequest); err != nil {
			log.Println("decoding user:", err)
			_ = helper.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		userID, err := userService.AuthenticateUser(r.Context(), userRequest)
		if err != nil {
			switch {
			case errors.Is(err, autherr.ErrInvalidCredentials):
				_ = helper.WriteJSONError(w, http.StatusUnauthorized, "invalid email or password")
			default:
				log.Println("logging in user:", err)
				_ = helper.WriteJSONError(w, http.StatusInternalServerError, helper.ErrInternal)
			}
			return
		}
		token, err := userService.GenerateJWT(userID)
		if err != nil {
			log.Println("generate jwt:", err)
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, helper.ErrInternal)
			return
		}

		if err = json.NewEncoder(w).Encode(models.JWTToken{Token: token}); err != nil {
			log.Println("encoding jwt:", err)
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, helper.ErrInternal)
			return
		}
	}
}
