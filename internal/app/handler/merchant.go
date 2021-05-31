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

type MerchantHandler interface {
	Create() http.HandlerFunc
	List() http.HandlerFunc
	Get() http.HandlerFunc
	Update() http.HandlerFunc
	UpdatePassword() http.HandlerFunc
	Delete() http.HandlerFunc
}

func NewMerchantHandler(merchantService service.MerchantService) MerchantHandler {
	return &merchantHandler{merchantService}
}

type merchantHandler struct {
	merchantService service.MerchantService
}

func (h *merchantHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.MerchantCreateRequest
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

		res, err := h.merchantService.Create(r.Context(), req)
		if err != nil {
			switch err {
			case constant.ErrEmailRegistered:
				web.MarshalError(w, http.StatusConflict, err)
				return
			default:
				web.MarshalError(w, http.StatusInternalServerError, err)
				return
			}
		}

		web.MarshalPayload(w, http.StatusCreated, res)
	}
}

func (h *merchantHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, offset, err := web.GetPagination(r)
		if err != nil {
			web.MarshalError(w, http.StatusBadRequest, err)
			return
		}

		req := model.MerchantListRequest{
			Limit:  limit,
			Offset: offset,
			Name:   web.GetUrlQueryString(r, "name"),
		}

		res, err := h.merchantService.List(r.Context(), req)
		if err != nil {
			web.MarshalError(w, http.StatusInternalServerError, err)
			return
		}

		web.MarshalPayload(w, http.StatusOK, res)
	}
}

func (h *merchantHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := web.GetUrlPathInt64(r, "merchant_id")
		if err != nil {
			web.MarshalError(w, http.StatusBadRequest, err)
			return
		}

		req := model.MerchantGetRequest{ID: id}
		res, err := h.merchantService.Get(r.Context(), req)
		if err != nil {
			switch err {
			case constant.ErrMerchantNotFound:
				web.MarshalError(w, http.StatusNotFound, err)
				return
			default:
				web.MarshalError(w, http.StatusInternalServerError, err)
				return
			}
		}

		web.MarshalPayload(w, http.StatusOK, res)
	}
}

func (h *merchantHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := web.GetUrlPathInt64(r, "merchant_id")
		if err != nil {
			web.MarshalError(w, http.StatusBadRequest, err)
			return
		}

		req := model.MerchantUpdateRequest{ID: id}
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			web.MarshalError(w, http.StatusBadRequest, constant.ErrRequestBody)
			return
		}

		err = validation.Struct(req)
		if err != nil {
			web.MarshalError(w, http.StatusBadRequest, err)
			return
		}

		res, err := h.merchantService.Update(r.Context(), req)
		if err != nil {
			switch err {
			case constant.ErrUnauthorized:
				web.MarshalError(w, http.StatusUnauthorized, err)
				return
			case constant.ErrEmailRegistered:
				web.MarshalError(w, http.StatusConflict, err)
				return
			case constant.ErrMerchantNotFound:
				web.MarshalError(w, http.StatusNotFound, err)
				return
			default:
				web.MarshalError(w, http.StatusInternalServerError, err)
				return
			}
		}

		web.MarshalPayload(w, http.StatusOK, res)
	}
}

func (h *merchantHandler) UpdatePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := web.GetUrlPathInt64(r, "merchant_id")
		if err != nil {
			web.MarshalError(w, http.StatusBadRequest, err)
			return
		}

		req := model.MerchantPasswordUpdateRequest{ID: id}
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			web.MarshalError(w, http.StatusBadRequest, constant.ErrRequestBody)
			return
		}

		err = validation.Struct(req)
		if err != nil {
			web.MarshalError(w, http.StatusBadRequest, err)
			return
		}

		res, err := h.merchantService.UpdatePassword(r.Context(), req)
		if err != nil {
			switch err {
			case constant.ErrUnauthorized, constant.ErrWrongPassword:
				web.MarshalError(w, http.StatusUnauthorized, err)
				return
			case constant.ErrMerchantNotFound:
				web.MarshalError(w, http.StatusNotFound, err)
				return
			default:
				web.MarshalError(w, http.StatusInternalServerError, err)
				return
			}
		}

		web.MarshalPayload(w, http.StatusOK, res)
	}
}

func (h *merchantHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := web.GetUrlPathInt64(r, "merchant_id")
		if err != nil {
			web.MarshalError(w, http.StatusBadRequest, err)
			return
		}

		req := model.MerchantDeleteRequest{ID: id}
		err = h.merchantService.Delete(r.Context(), req)
		if err != nil {
			switch err {
			case constant.ErrUnauthorized:
				web.MarshalError(w, http.StatusUnauthorized, err)
				return
			case constant.ErrMerchantNotFound:
				web.MarshalError(w, http.StatusNotFound, err)
				return
			default:
				web.MarshalError(w, http.StatusInternalServerError, err)
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
