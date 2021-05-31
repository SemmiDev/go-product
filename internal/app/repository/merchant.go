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

type MerchantRepository interface {
	Create(ctx context.Context, merchant *model.Merchant) error
	List(ctx context.Context, limit, offset int, name string) ([]*model.Merchant, error)
	Get(ctx context.Context, id int64) (*model.Merchant, error)
	GetByEmail(ctx context.Context, email string) (*model.Merchant, error)
	Update(ctx context.Context, merchant *model.Merchant) error
	Delete(ctx context.Context, id int64) error
}

func NewMerchantRepository(mysqlClient mysql.Client, redisClient redis.Client) MerchantRepository {
	return &merchantRepository{mysqlClient, redisClient}
}

type merchantRepository struct {
	mysqlClient mysql.Client
	redisClient redis.Client
}

func (r *merchantRepository) Create(ctx context.Context, merchant *model.Merchant) error {
	res, err := r.mysqlClient.Conn().ExecContext(ctx, `
	INSERT INTO
		merchant (name, email, password, created_at)
	VALUES
		(?, ?, ?, ?)
	`, merchant.Name, merchant.Email, merchant.Password, merchant.CreatedAt)
	if err != nil {
		return err
	}

	merchant.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}

	temp, err := r.Get(ctx, merchant.ID)
	*merchant = *temp
	return err
}

func (r *merchantRepository) List(ctx context.Context, limit, offset int, name string) ([]*model.Merchant, error) {
	var merchants []*model.Merchant
	rows, err := r.mysqlClient.Conn().QueryContext(ctx, `
	SELECT
		id, name, email, created_at, updated_at
	FROM
		merchant
	WHERE
		name LIKE ?
	LIMIT
		? OFFSET ?
	`, "%"+name+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		merchant := new(model.Merchant)
		err := rows.Scan(&merchant.ID, &merchant.Name, &merchant.Email, &merchant.CreatedAt, &merchant.UpdatedAt)
		if err != nil {
			return nil, err
		}
		merchants = append(merchants, merchant)
	}

	return merchants, nil
}

func (r *merchantRepository) Get(ctx context.Context, id int64) (*model.Merchant, error) {
	merchant := new(model.Merchant)
	err := r.redisClient.Cache().Get(ctx, fmt.Sprintf("merchant_%d", id), merchant)
	if err != nil && err != cache.ErrCacheMiss {
		return nil, err
	} else if err == nil {
		return merchant, nil
	}

	err = r.mysqlClient.Conn().QueryRowContext(ctx, `
	SELECT
		id, name, email, password, created_at, updated_at
	FROM
		merchant
	WHERE
		id = ?
	`, id,
	).Scan(&merchant.ID, &merchant.Name, &merchant.Email, &merchant.Password, &merchant.CreatedAt, &merchant.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return merchant, r.redisClient.Cache().Set(&cache.Item{
		Ctx:   ctx,
		Key:   fmt.Sprintf("merchant_%d", id),
		Value: merchant,
		TTL:   config.Cfg().RedisTTL,
	})
}

func (r *merchantRepository) GetByEmail(ctx context.Context, email string) (*model.Merchant, error) {
	merchant := new(model.Merchant)
	err := r.redisClient.Cache().Get(ctx, fmt.Sprintf("merchant_%s", email), merchant)
	if err != nil && err != cache.ErrCacheMiss {
		return nil, err
	} else if err == nil {
		return merchant, nil
	}

	err = r.mysqlClient.Conn().QueryRowContext(ctx, `
	SELECT
		id, name, email, password, created_at, updated_at
	FROM
		merchant
	WHERE
		email = ?
	`, email,
	).Scan(&merchant.ID, &merchant.Name, &merchant.Email, &merchant.Password, &merchant.CreatedAt, &merchant.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return merchant, r.redisClient.Cache().Set(&cache.Item{
		Ctx:   ctx,
		Key:   fmt.Sprintf("merchant_%s", email),
		Value: merchant,
		TTL:   config.Cfg().RedisTTL,
	})
}

func (r *merchantRepository) Update(ctx context.Context, merchant *model.Merchant) error {
	_, err := r.mysqlClient.Conn().ExecContext(ctx, `
	UPDATE
		merchant
	SET
		name = ?, email = ?, password = ?, updated_at = ?
	WHERE
		id = ?
	`, merchant.Name, merchant.Email, merchant.Password, merchant.UpdatedAt.Time, merchant.ID)
	if err != nil {
		return err
	}

	err = r.redisClient.Cache().Delete(ctx, fmt.Sprintf("merchant_%d", merchant.ID))
	if err != nil && err != cache.ErrCacheMiss {
		return err
	}

	temp, err := r.Get(ctx, merchant.ID)
	*merchant = *temp
	return err
}

func (r *merchantRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.mysqlClient.Conn().ExecContext(ctx, `
	DELETE FROM
		merchant
	WHERE
		id = ?
	`, id)
	if err != nil {
		return err
	}

	err = r.redisClient.Cache().Delete(ctx, fmt.Sprintf("merchant_%d", id))
	if err != nil && err != cache.ErrCacheMiss {
		return err
	}

	return nil
}
