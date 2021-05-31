package model

import (
	"database/sql"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

type Merchant struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

func (a *Merchant) GenerateClaims() jwt.MapClaims {
	return jwt.MapClaims{"id": a.ID}
}

type MerchantCreateRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8"`
}

type MerchantListRequest struct {
	Limit  int
	Offset int
	Name   string
}

type MerchantGetRequest struct {
	ID int64
}

type MerchantUpdateRequest struct {
	ID    int64  `json:"-"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type MerchantPasswordUpdateRequest struct {
	ID          int64  `json:"-"`
	OldPassword string `json:"old_password" validate:"required,gte=8"`
	NewPassword string `json:"new_password" validate:"required,gte=8"`
}

type MerchantDeleteRequest struct {
	ID int64
}

type MerchantResponse struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func NewMerchantResponse(payload *Merchant) *MerchantResponse {
	res := &MerchantResponse{
		ID:        payload.ID,
		Name:      payload.Name,
		Email:     payload.Email,
		CreatedAt: payload.CreatedAt,
	}
	if payload.UpdatedAt.Valid {
		res.UpdatedAt = &payload.UpdatedAt.Time
	}
	return res
}

func NewMerchantListResponse(payloads []*Merchant) []*MerchantResponse {
	res := make([]*MerchantResponse, len(payloads))
	for i, payload := range payloads {
		res[i] = NewMerchantResponse(payload)
	}
	return res
}