package model

import (
	"database/sql"
	"time"
)

type Product struct {
	ID        int64
	Name      string
	Price     float32
	CreatedAt time.Time
	UpdatedAt sql.NullTime

	MerchantID int64
	Merchant   Merchant
}

type ProductCreateRequest struct {
	Name      string `json:"name" validate:"required"`
	Price     float32 `json:"price" validate:"required"`
}

type ProductListRequest struct {
	Limit  int
	Offset int
	Title  string
}

type ProductGetRequest struct {
	ID int64
}

type ProductUpdateRequest struct {
	ID    int64  `json:"-"`
	Name      string `json:"name" validate:"required"`
	Price     float32 `json:"price" validate:"required"`
}

type ProductDeleteRequest struct {
	ID int64
}

type ProductResponse struct {
	ID        int64      `json:"id"`
	Name      string `json:"name"`
	Price     float32 `json:"price"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	MerchantID int64             `json:"merchant_id"`
	Merchant   *MerchantResponse `json:"merchant"`
}

func NewProductResponse(payload *Product) *ProductResponse {
	res := &ProductResponse{
		ID:         payload.ID,
		Name: payload.Name,
		Price: payload.Price,
		CreatedAt:  payload.CreatedAt,
		MerchantID: payload.MerchantID,
		Merchant:   NewMerchantResponse(&payload.Merchant),
	}
	if payload.UpdatedAt.Valid {
		res.UpdatedAt = &payload.UpdatedAt.Time
	}
	return res
}

func NewProductListResponse(payloads []*Product) []*ProductResponse {
	res := make([]*ProductResponse, len(payloads))
	for i, payload := range payloads {
		res[i] = NewProductResponse(payload)
	}
	return res
}
