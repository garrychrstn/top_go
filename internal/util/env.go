package util

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/garrychrstn/top-go/internal/repository/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

// LoadDotEnv loads KEY=VALUE pairs from path into the environment without
// overriding variables that are already set. A missing file is not an error.
// Minimal subset of godotenv: skips blank lines and # comments, and strips a
// single surrounding pair of quotes from values.
func LoadDotEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("open env file: %w", err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return fmt.Errorf("%s:%d: expected KEY=VALUE", path, lineNo)
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if len(value) >= 2 && (value[0] == '"' && value[len(value)-1] == '"' ||
			value[0] == '\'' && value[len(value)-1] == '\'') {
			value = value[1 : len(value)-1]
		}

		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}
	return sc.Err()
}

// InitiateService ensures an initial user exists. First searches for any user,
// and if none found, creates user "admin" with password "admin" using HashPassword.
func InitiateService(ctx context.Context, pool *pgxpool.Pool) error {
	q := db.New(pool)

	users, err := q.LisstUser(ctx)
	if err != nil {
		return fmt.Errorf("list users: %w", err)
	}
	if len(users) > 0 {
		return nil
	}

	username := os.Getenv("INIT_USERNAME")
	if username == "" {
		username = "admin"
	}
	password := os.Getenv("INIT_PASSWORD")
	if password == "" {
		password = "admin"
	}

	encPassword, err := HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	_, err = q.CreateUser(ctx, db.CreateUserParams{
		Username: username,
		Password: encPassword,
	})
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	slog.Info("initialized user", "username", username)
	return nil
}
