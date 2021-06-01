package service

import (
	"context"
	"database/sql"
	"github.com/SemmiDev/go-product/internal/app/model"
	"github.com/SemmiDev/go-product/internal/app/repository"
	"github.com/SemmiDev/go-product/internal/constant"
	"github.com/SemmiDev/go-product/internal/logger"
	"github.com/SemmiDev/go-product/internal/security/middleware"
	"time"
)

type ProductService interface {
	Create(ctx context.Context, req model.ProductCreateRequest) (*model.ProductResponse, error)
	List(ctx context.Context, req model.ProductListRequest) ([]*model.ProductResponse, error)
	Get(ctx context.Context, req model.ProductGetRequest) (*model.ProductResponse, error)
	Update(ctx context.Context, req model.ProductUpdateRequest) (*model.ProductResponse, error)
	Delete(ctx context.Context, req model.ProductDeleteRequest) error
}

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &productService{productRepository}
}

func (s *productService) Create(ctx context.Context, req model.ProductCreateRequest) (*model.ProductResponse, error) {
	claimsID, valid := middleware.GetClaimsID(ctx)
	if !valid {
		return nil, constant.ErrUnauthorized
	}

	product := &model.Product{
		Name:       req.Name,
		Price:      req.Price,
		CreatedAt:  time.Now(),
		MerchantID: claimsID,
	}

	err := s.productRepository.Create(ctx, product)
	if err != nil {
		logger.Log().Err(err).Msg("failed to create product")
		return nil, constant.ErrServer
	}

	return model.NewProductResponse(product), nil
}

func (s *productService) List(ctx context.Context, req model.ProductListRequest) ([]*model.ProductResponse, error) {
	products, err := s.productRepository.List(ctx, req.Limit, req.Offset, req.Title)
	if err != nil {
		logger.Log().Err(err).Msg("failed to list products")
		return nil, constant.ErrServer
	}

	return model.NewProductListResponse(products), nil
}

func (s *productService) Get(ctx context.Context, req model.ProductGetRequest) (*model.ProductResponse, error) {
	product, err := s.productRepository.Get(ctx, req.ID)
	if err != nil {
		return nil, s.switchErrProductNotFoundOrErrServer(err)
	}
	return model.NewProductResponse(product), nil
}

func (s *productService) Update(ctx context.Context, req model.ProductUpdateRequest) (*model.ProductResponse, error) {
	product, err := s.productRepository.Get(ctx, req.ID)
	if err != nil {
		return nil, s.switchErrProductNotFoundOrErrServer(err)
	}

	if !middleware.IsMe(ctx, product.MerchantID) {
		return nil, constant.ErrUnauthorized
	}

	product.Name = req.Name
	product.Price = req.Price
	product.UpdatedAt.Time = time.Now()

	err = s.productRepository.Update(ctx, product)
	if err != nil {
		return nil, s.switchErrProductNotFoundOrErrServer(err)
	}

	return model.NewProductResponse(product), nil
}

func (s *productService) Delete(ctx context.Context, req model.ProductDeleteRequest) error {
	product, err := s.productRepository.Get(ctx, req.ID)
	if err != nil {
		return s.switchErrProductNotFoundOrErrServer(err)
	}

	if !middleware.IsMe(ctx, product.MerchantID) {
		return constant.ErrUnauthorized
	}

	err = s.productRepository.Delete(ctx, req.ID)
	if err != nil {
		return s.switchErrProductNotFoundOrErrServer(err)
	}

	return nil
}

func (s *productService) switchErrProductNotFoundOrErrServer(err error) error {
	switch err {
	case sql.ErrNoRows:
		return constant.ErrPostNotFound
	default:
		logger.Log().Err(err).Msg("failed to execute operation product repository")
		return constant.ErrServer
	}
}
