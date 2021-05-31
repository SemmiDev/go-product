package service

import (
	"context"
	"database/sql"
	"github.com/SemmiDev/go-product/internal/app/model"
	"github.com/SemmiDev/go-product/internal/app/repository"
	"github.com/SemmiDev/go-product/internal/constant"
	"github.com/SemmiDev/go-product/internal/logger"
	"github.com/SemmiDev/go-product/internal/security/middleware"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type MerchantService interface {
	Create(ctx context.Context, req model.MerchantCreateRequest) (*model.MerchantResponse, error)
	List(ctx context.Context, req model.MerchantListRequest) ([]*model.MerchantResponse, error)
	Get(ctx context.Context, req model.MerchantGetRequest) (*model.MerchantResponse, error)
	Update(ctx context.Context, req model.MerchantUpdateRequest) (*model.MerchantResponse, error)
	UpdatePassword(ctx context.Context, req model.MerchantPasswordUpdateRequest) (*model.MerchantResponse, error)
	Delete(ctx context.Context, req model.MerchantDeleteRequest) error
}

func NewMerchantService(accountRepository repository.MerchantRepository) MerchantService {
	return &accountService{accountRepository}
}

type accountService struct {
	accountRepository repository.MerchantRepository
}

func (s *accountService) Create(ctx context.Context, req model.MerchantCreateRequest) (*model.MerchantResponse, error) {
	_, err := s.accountRepository.GetByEmail(ctx, req.Email)
	if err != nil && err != sql.ErrNoRows {
		logger.Log().Err(err).Msg("failed to get account by email")
		return nil, constant.ErrServer
	} else if err == nil {
		return nil, constant.ErrEmailRegistered
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log().Err(err).Msg("failed to generate from password")
		return nil, constant.ErrServer
	}

	account := &model.Merchant{
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(password),
		CreatedAt: time.Now(),
	}

	err = s.accountRepository.Create(ctx, account)
	if err != nil {
		logger.Log().Err(err).Msg("failed to create account")
		return nil, constant.ErrServer
	}

	return model.NewMerchantResponse(account), nil
}

func (s *accountService) List(ctx context.Context, req model.MerchantListRequest) ([]*model.MerchantResponse, error) {
	accounts, err := s.accountRepository.List(ctx, req.Limit, req.Offset, req.Name)
	if err != nil {
		logger.Log().Err(err).Msg("failed to list accounts")
		return nil, constant.ErrServer
	}

	return model.NewMerchantListResponse(accounts), nil
}

func (s *accountService) Get(ctx context.Context, req model.MerchantGetRequest) (*model.MerchantResponse, error) {
	account, err := s.accountRepository.Get(ctx, req.ID)
	if err != nil {
		return nil, s.switchErrMerchantNotFoundOrErrServer(err)
	}

	return model.NewMerchantResponse(account), nil
}

func (s *accountService) Update(ctx context.Context, req model.MerchantUpdateRequest) (*model.MerchantResponse, error) {
	if !middleware.IsMe(ctx, req.ID) {
		return nil, constant.ErrUnauthorized
	}

	account, err := s.accountRepository.GetByEmail(ctx, req.Email)
	if err != nil && err != sql.ErrNoRows {
		logger.Log().Err(err).Msg("failed to get account by email")
		return nil, constant.ErrServer
	} else if err == nil && account.ID != req.ID {
		return nil, constant.ErrEmailRegistered
	}

	account, err = s.accountRepository.Get(ctx, req.ID)
	if err != nil {
		return nil, s.switchErrMerchantNotFoundOrErrServer(err)
	}

	account.Name = req.Name
	account.Email = req.Email
	account.UpdatedAt.Time = time.Now()

	err = s.accountRepository.Update(ctx, account)
	if err != nil {
		return nil, s.switchErrMerchantNotFoundOrErrServer(err)
	}

	return model.NewMerchantResponse(account), nil
}

func (s *accountService) UpdatePassword(ctx context.Context, req model.MerchantPasswordUpdateRequest) (*model.MerchantResponse, error) {
	if !middleware.IsMe(ctx, req.ID) {
		return nil, constant.ErrUnauthorized
	}

	account, err := s.accountRepository.Get(ctx, req.ID)
	if err != nil {
		return nil, s.switchErrMerchantNotFoundOrErrServer(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(req.OldPassword))
	if err != nil {
		return nil, constant.ErrWrongPassword
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Log().Err(err).Msg("failed to generate from password")
		return nil, constant.ErrServer
	}

	account.Password = string(password)
	account.UpdatedAt.Time = time.Now()

	err = s.accountRepository.Update(ctx, account)
	if err != nil {
		return nil, s.switchErrMerchantNotFoundOrErrServer(err)
	}

	return model.NewMerchantResponse(account), nil
}

func (s *accountService) Delete(ctx context.Context, req model.MerchantDeleteRequest) error {
	if !middleware.IsMe(ctx, req.ID) {
		return constant.ErrUnauthorized
	}

	err := s.accountRepository.Delete(ctx, req.ID)
	if err != nil {
		return s.switchErrMerchantNotFoundOrErrServer(err)
	}

	return nil
}

func (s *accountService) switchErrMerchantNotFoundOrErrServer(err error) error {
	switch err {
	case sql.ErrNoRows:
		return constant.ErrMerchantNotFound
	default:
		logger.Log().Err(err).Msg("failed to execute operation account repository")
		return constant.ErrServer
	}
}
