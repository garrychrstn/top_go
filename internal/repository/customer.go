package repository

import (
	"context"

	"github.com/garrychrstn/top-go/internal/repository/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerRepository interface {
	Create(ctx context.Context, name string, phoneNumber string, address pgtype.Text) (*db.Customer, error)
	GetByID(ctx context.Context, id uuid.UUID) (*db.Customer, error)
	GetByName(ctx context.Context, name string) (*db.Customer, error)
	List(ctx context.Context) ([]db.Customer, error)
	Update(ctx context.Context, id uuid.UUID, name string, phoneNumber string, address pgtype.Text) (*db.Customer, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type customerRepo struct {
	query *db.Queries
}

func CustomerInitRepo(pool *pgxpool.Pool) CustomerRepository {
	return &customerRepo{query: db.New(pool)}
}

func (r *customerRepo) Create(ctx context.Context, name string, phoneNumber string, address pgtype.Text) (*db.Customer, error) {
	c, err := r.query.CustomerCreate(ctx, db.CustomerCreateParams{
		Name:        name,
		PhoneNumber: phoneNumber,
		Address:     address,
	})
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *customerRepo) GetByID(ctx context.Context, id uuid.UUID) (*db.Customer, error) {
	c, err := r.query.CustomerGetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *customerRepo) GetByName(ctx context.Context, name string) (*db.Customer, error) {
	c, err := r.query.CustomerGetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *customerRepo) List(ctx context.Context) ([]db.Customer, error) {
	return r.query.CustomerList(ctx)
}

func (r *customerRepo) Update(ctx context.Context, id uuid.UUID, name string, phoneNumber string, address pgtype.Text) (*db.Customer, error) {
	c, err := r.query.CustomerUpdate(ctx, db.CustomerUpdateParams{
		ID:          id,
		Name:        name,
		PhoneNumber: phoneNumber,
		Address:     address,
	})
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *customerRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.query.CustomerDelete(ctx, id)
}
