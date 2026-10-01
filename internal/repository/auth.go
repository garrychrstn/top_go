package repository

import (
	"context"
	"errors"

	"github.com/garrychrstn/top-go/internal/repository/db"
	"github.com/garrychrstn/top-go/internal/util"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	GetUser(ctx context.Context, username string) (*db.User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, password string) error
}

type repo struct {
	query *db.Queries
}

func AuthInitRepo(pool *pgxpool.Pool) AuthRepository {
	return &repo{query: db.New(pool)}
}

func (r *repo) GetUser(ctx context.Context, username string) (*db.User, error) {
	user, err := r.query.UserGetByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("User not found")
	}
	return &user, nil
}
func (r *repo) UpdatePassword(ctx context.Context, id uuid.UUID, password string) error {
	hash, err := util.HashPassword(password)
	if err != nil {
		return err
	}

	_, error := r.query.UserPasswordReset(ctx, db.UserPasswordResetParams{
		Password: hash,
		ID:       id,
	})
	if error != nil {
		return error
	}
	return nil
}
