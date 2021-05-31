package handler

import (
	"encoding/json"
	"github.com/SemmiDev/go-product/internal/app/model"
	"github.com/SemmiDev/go-product/internal/app/service"
	"github.com/SemmiDev/go-product/internal/constant"
	"github.com/SemmiDev/go-product/internal/validation"
	"github.com/SemmiDev/go-product/internal/web"
	"net/http"

)

type AuthHandler interface {
	Login() http.HandlerFunc
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{authService}
}

type authHandler struct {
	authService service.AuthService
}

func (h *authHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.AuthRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			web.MarshalError(w, http.StatusBadRequest, constant.ErrRequestBody)
			return
		}

		err = validation.Struct(req)
		if err != nil {
			web.MarshalError(w, http.StatusBadRequest, err)
			return
		}

		res, err := h.authService.Login(r.Context(), req)
		if err != nil {
			switch err {
			case constant.ErrEmailNotRegistered, constant.ErrWrongPassword:
				web.MarshalError(w, http.StatusUnauthorized, err)
				return
			default:
				web.MarshalError(w, http.StatusInternalServerError, err)
				return
			}
		}

		web.MarshalPayload(w, http.StatusOK, res)
	}
}
