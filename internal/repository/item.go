package repository

import (
	"context"

	"github.com/garrychrstn/top-go/internal/repository/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemRepository interface {
	Create(ctx context.Context, name string, price pgtype.Numeric) (*db.Item, error)
	List(ctx context.Context) ([]db.Item, error)
	Update(ctx context.Context, id uuid.UUID, name string, price pgtype.Numeric) (*db.Item, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type itemRepo struct {
	query *db.Queries
}

func ItemInitRepo(pool *pgxpool.Pool) ItemRepository {
	return &itemRepo{query: db.New(pool)}
}

func (r *itemRepo) Create(ctx context.Context, name string, price pgtype.Numeric) (*db.Item, error) {
	item, err := r.query.ItemCreate(ctx, db.ItemCreateParams{
		Name:  name,
		Price: price,
	})
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *itemRepo) List(ctx context.Context) ([]db.Item, error) {
	return r.query.ItemList(ctx)
}

func (r *itemRepo) Update(ctx context.Context, id uuid.UUID, name string, price pgtype.Numeric) (*db.Item, error) {
	item, err := r.query.ItemUpdate(ctx, db.ItemUpdateParams{
		ID:    id,
		Name:  name,
		Price: price,
	})
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *itemRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.query.ItemDelete(ctx, id)
}
