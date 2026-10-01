package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/garrychrstn/top-go/internal/repository/db"
	"github.com/garrychrstn/top-go/internal/types"
)

// UserRepository is the persistence boundary used by handlers.
// It wraps the sqlc-generated queries and maps errors to sentinels.
type UserRepository interface {
	List(ctx context.Context) ([]types.UserResponse, error)
	GetByID(ctx context.Context, id int64) (types.UserResponse, error)
	Create(ctx context.Context, req types.CreateUserRequest) (types.UserResponse, error)
}

type userRepository struct {
	q *db.Queries
}

// NewUserRepository builds the UserRepository from a pgx pool
// (constructor injection).
func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{q: db.New(pool)}
}

func (r *userRepository) List(ctx context.Context) ([]types.UserResponse, error) {
	users, err := r.q.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return toUsers(users), nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (types.UserResponse, error) {
	user, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return types.UserResponse{}, ErrNotFound
		}
		return types.UserResponse{}, fmt.Errorf("get user by id: %w", err)
	}
	return toUser(user), nil
}

func (r *userRepository) Create(ctx context.Context, req types.CreateUserRequest) (types.UserResponse, error) {
	user, err := r.q.CreateUser(ctx, db.CreateUserParams{
		Name:  strings.TrimSpace(req.Name),
		Email: strings.ToLower(strings.TrimSpace(req.Email)),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return types.UserResponse{}, ErrEmailTaken
		}
		return types.UserResponse{}, fmt.Errorf("create user: %w", err)
	}
	return toUser(user), nil
}

func toUser(u db.User) types.UserResponse {
	return types.UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

func toUsers(users []db.User) []types.UserResponse {
	out := make([]types.UserResponse, len(users))
	for i, u := range users {
		out[i] = toUser(u)
	}
	return out
}
