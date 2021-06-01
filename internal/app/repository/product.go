package repository

import (
	"context"
	"fmt"
	"github.com/SemmiDev/go-product/internal/app/model"
	"github.com/SemmiDev/go-product/internal/config"
	"github.com/SemmiDev/go-product/internal/db/mysql"
	"github.com/SemmiDev/go-product/internal/db/redis"
	"github.com/go-redis/cache/v8"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	List(ctx context.Context, limit, offset int, title string) ([]*model.Product, error)
	Get(ctx context.Context, id int64) (*model.Product, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id int64) error
}

type productRepository struct {
	mysqlClient mysql.Client
	redisClient redis.Client
}

func NewProductRepository(mysqlClient mysql.Client, redisClient redis.Client) ProductRepository {
	return &productRepository{mysqlClient, redisClient}
}

func (r *productRepository) Create(ctx context.Context, product *model.Product) error {
	res, err := r.mysqlClient.Conn().ExecContext(ctx, `
	INSERT INTO
		product (name, price, merchant_id, created_at)
	VALUES
		(?, ?, ?, ?)
	`, product.Name, product.Price, product.MerchantID, product.CreatedAt)
	if err != nil {
		return err
	}

	product.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}

	temp, err := r.Get(ctx, product.ID)
	*product = *temp
	return nil
}

func (r productRepository) List(ctx context.Context, limit, offset int, title string) ([]*model.Product, error) {
	var products []*model.Product
	rows, err := r.mysqlClient.Conn().QueryContext(ctx, `
	SELECT
		product.id, product.name, product.price, product.created_at, product.updated_at, product.merchant_id,
		merchant.id, merchant.name, merchant.email, merchant.password, merchant.created_at, merchant.updated_at
	FROM
		product
	INNER JOIN
		merchant
	ON
		product.merchant_id = merchant.id
	WHERE
		product.name LIKE ?
	LIMIT
		? OFFSET ?
	`, "%"+title+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		product := new(model.Product)
		err := rows.Scan(
			&product.ID, &product.Name, &product.Price, &product.CreatedAt, &product.UpdatedAt, &product.MerchantID,
			&product.Merchant.ID, &product.Merchant.Name, &product.Merchant.Email, &product.Merchant.Password, &product.Merchant.CreatedAt, &product.Merchant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r productRepository) Get(ctx context.Context, id int64) (*model.Product, error) {
	product := new(model.Product)
	err := r.redisClient.Cache().Get(ctx, fmt.Sprintf("product_%d", id), product)
	if err != nil && err != cache.ErrCacheMiss {
		return nil, err
	} else if err == nil {
		return product, nil
	}

	err = r.mysqlClient.Conn().QueryRowContext(ctx, `
	SELECT
		product.id, product.name, product.price, product.created_at, product.updated_at, product.merchant_id,
		merchant.id, merchant.name, merchant.email, merchant.password, merchant.created_at, merchant.updated_at
	FROM
		product
	INNER JOIN
		merchant
	ON
		product.merchant_id = merchant.id
	WHERE
		product.id = ?
	`, id,
	).Scan(
		&product.ID, &product.Name, &product.Price, &product.CreatedAt, &product.UpdatedAt, &product.MerchantID,
		&product.Merchant.ID, &product.Merchant.Name, &product.Merchant.Email, &product.Merchant.Password, &product.Merchant.CreatedAt, &product.Merchant.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return product, r.redisClient.Cache().Set(&cache.Item{
		Ctx:   ctx,
		Key:   fmt.Sprintf("product_%d", id),
		Value: product,
		TTL:   config.Cfg().RedisTTL,
	})
}

func (r productRepository) Update(ctx context.Context, product *model.Product) error {
	_, err := r.mysqlClient.Conn().ExecContext(ctx, `
	UPDATE
		product
	SET
		name = ?, price = ?, updated_at = ?
	WHERE
		id = ?
	`, product.Name, product.Price, product.UpdatedAt.Time, product.ID)
	if err != nil {
		return err
	}

	err = r.redisClient.Cache().Delete(ctx, fmt.Sprintf("product_%d", product.ID))
	if err != nil && err != cache.ErrCacheMiss {
		return err
	}

	temp, err := r.Get(ctx, product.ID)
	*product = *temp
	return err
}

func (r productRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.mysqlClient.Conn().ExecContext(ctx, `
	DELETE FROM
		product
	WHERE
		id = ?
	`, id)
	if err != nil {
		return err
	}

	err = r.redisClient.Cache().Delete(ctx, fmt.Sprintf("product_%d", id))
	if err != nil && err != cache.ErrCacheMiss {
		return err
	}

	return nil
}
