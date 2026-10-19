package auth

import (
	"context"
	"database/sql"
	"testing"

	repository "github.com/rawbil/movie-stream-app/internal/adapters/sqlc"
	"github.com/rawbil/movie-stream-app/internal/utils"
)

func TestAuthServiceValidationAndNotFoundPaths(t *testing.T) {
	t.Run("login rejects invalid input before repository access", func(t *testing.T) {
		svc := NewService(&MockService{}, nil)
		_, _, _, err := svc.LoginUser(context.Background(), utils.LoginParams{})
		if err != utils.AllFieldsRequiredError {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("login maps missing user", func(t *testing.T) {
		mock := &MockService{GetUserByEmailFunc: func(context.Context, string) (repository.GetUserByEmailRow, error) { return repository.GetUserByEmailRow{}, sql.ErrNoRows }}
		svc := NewService(mock, nil)
		_, _, _, err := svc.LoginUser(context.Background(), utils.LoginParams{Email: "user@example.com", Password: "password"})
		if err != utils.NoRecordError {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("create role requires a value", func(t *testing.T) {
		svc := NewService(&MockService{}, nil)
		if _, err := svc.CreateRole(context.Background(), ""); err != utils.AllFieldsRequiredError {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("create role maps duplicate", func(t *testing.T) {
		mock := &MockService{GetRoleFunc: func(context.Context, string) (repository.Role, error) { return repository.Role{Role: "ADMIN"}, nil }}
		svc := NewService(mock, nil)
		if _, err := svc.CreateRole(context.Background(), "admin"); err != utils.DuplicateRecordError {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("logout maps missing user", func(t *testing.T) {
		mock := &MockService{GetUserByIDFunc: func(context.Context, int64) (repository.User, error) { return repository.User{}, sql.ErrNoRows }}
		svc := NewService(mock, nil)
		if err := svc.Logout(context.Background(), 9); err != utils.NoRecordError {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("refresh rejects empty and malformed tokens", func(t *testing.T) {
		svc := NewService(&MockService{}, nil)
		if _, _, _, err := svc.RefreshTokens(context.Background(), ""); err != utils.AllFieldsRequiredError {
			t.Fatalf("empty token err=%v", err)
		}
		if _, _, _, err := svc.RefreshTokens(context.Background(), "not-a-jwt"); err == nil {
			t.Fatal("expected malformed token error")
		}
	})
}
