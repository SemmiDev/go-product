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

type ProductHandler interface {
	Create() http.HandlerFunc
	List() http.HandlerFunc
	Get() http.HandlerFunc
	Update() http.HandlerFunc
	Delete() http.HandlerFunc
}

type productHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) ProductHandler {
	return &productHandler{productService}
}

func (h *productHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.ProductCreateRequest
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

		res, err := h.productService.Create(r.Context(), req)
		if err != nil {
			switch err {
			case constant.ErrUnauthorized:
				web.MarshalError(w, http.StatusUnauthorized, err)
				return
			default:
				web.MarshalError(w, http.StatusInternalServerError, err)
				return
			}
		}

		web.MarshalPayload(w, http.StatusCreated, res)
	}
}

func (h *productHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, offset, err := web.GetPagination(r)
		if err != nil {
			web.MarshalError(w, http.StatusBadRequest, err)
			return
		}

		req := model.ProductListRequest{
			Limit:  limit,
			Offset: offset,
			Title:  web.GetUrlQueryString(r, "name"),
		}

		res, err := h.productService.List(r.Context(), req)
		if err != nil {
			web.MarshalError(w, http.StatusInternalServerError, err)
			return
		}

		web.MarshalPayload(w, http.StatusOK, res)
	}
}

func (h *productHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := web.GetUrlPathInt64(r, "product_id")
		if err != nil {
			web.MarshalError(w, http.StatusBadRequest, err)
			return
		}

		req := model.ProductGetRequest{ID: id}
		res, err := h.productService.Get(r.Context(), req)
		if err != nil {
			switch err {
			case constant.ErrProductNotFound:
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

func (h *productHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := web.GetUrlPathInt64(r, "product_id")
		if err != nil {
			web.MarshalError(w, http.StatusBadRequest, err)
			return
		}

		req := model.ProductUpdateRequest{ID: id}
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

		res, err := h.productService.Update(r.Context(), req)
		if err != nil {
			switch err {
			case constant.ErrUnauthorized:
				web.MarshalError(w, http.StatusUnauthorized, err)
				return
			case constant.ErrProductNotFound:
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

func (h *productHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := web.GetUrlPathInt64(r, "product_id")
		if err != nil {
			web.MarshalError(w, http.StatusBadRequest, err)
			return
		}

		req := model.ProductDeleteRequest{ID: id}
		err = h.productService.Delete(r.Context(), req)
		if err != nil {
			switch err {
			case constant.ErrUnauthorized:
				web.MarshalError(w, http.StatusUnauthorized, err)
				return
			case constant.ErrProductNotFound:
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
