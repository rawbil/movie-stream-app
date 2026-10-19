package auth

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	repository "github.com/rawbil/movie-stream-app/internal/adapters/sqlc"
	"github.com/rawbil/movie-stream-app/internal/utils"
)

type MockService struct {
	CreateUserFunc      func(ctx context.Context, arg repository.CreateUserParams) (sql.Result, error)
	GetRoleFunc         func(ctx context.Context, role string) (repository.Role, error)
	CreateRoleFunc      func(ctx context.Context, role string) (sql.Result, error)
	CreateUserRolesFunc func(ctx context.Context, arg repository.CreateUserRolesParams) (sql.Result, error)
	GetGenreFunc        func(ctx context.Context, genreName string) (repository.Genre, error)
	CreateGenreFunc     func(ctx context.Context, genreName string) (sql.Result, error)
	CreateUserGenreFunc func(ctx context.Context, arg repository.CreateUserGenreParams) (sql.Result, error)
	GetUserByEmailFunc  func(ctx context.Context, email string) (repository.GetUserByEmailRow, error)

	GetRefreshTokenFunc    func(ctx context.Context, userID int64) (repository.RefreshToken, error)
	UpdateRefreshTokenFunc func(ctx context.Context, arg repository.UpdateRefreshTokenParams) (sql.Result, error)
	CreateRefreshTokenFunc func(ctx context.Context, arg repository.CreateRefreshTokenParams) (sql.Result, error)
	GetUserByIDFunc        func(ctx context.Context, userID int64) (repository.User, error)
	RevokeRefreshTokenFunc func(ctx context.Context, arg repository.RevokeRefreshTokenParams) (sql.Result, error)
}

func (m *MockService) CreateUser(ctx context.Context, arg repository.CreateUserParams) (sql.Result, error) {
	return m.CreateUserFunc(ctx, arg)
}

func (m *MockService) GetRole(ctx context.Context, role string) (repository.Role, error) {
	return m.GetRoleFunc(ctx, role)
}

func (m *MockService) CreateRole(ctx context.Context, role string) (sql.Result, error) {
	return m.CreateRole(ctx, role)
}

func (m *MockService) CreateUserRoles(ctx context.Context, arg repository.CreateUserRolesParams) (sql.Result, error) {
	return m.CreateUserRolesFunc(ctx, arg)
}

func (m *MockService) CreateGenre(ctx context.Context, genreName string) (sql.Result, error) {
	return m.CreateGenreFunc(ctx, genreName)
}

func (m *MockService) CreateUserGenre(ctx context.Context, arg repository.CreateUserGenreParams) (sql.Result, error) {
	return m.CreateUserGenreFunc(ctx, arg)
}

func (m *MockService) GetUserByEmail(ctx context.Context, email string) (repository.GetUserByEmailRow, error) {
	return m.GetUserByEmailFunc(ctx, email)
}

func (m *MockService) GetRefreshToken(ctx context.Context, userID int64) (repository.RefreshToken, error) {
	return m.GetRefreshTokenFunc(ctx, userID)
}

func (m *MockService) UpdateRefreshToken(ctx context.Context, arg repository.UpdateRefreshTokenParams) (sql.Result, error) {
	return m.UpdateRefreshTokenFunc(ctx, arg)
}

func (m *MockService) CreateRefreshToken(ctx context.Context, arg repository.CreateRefreshTokenParams) (sql.Result, error) {
	return m.CreateRefreshTokenFunc(ctx, arg)
}

func (m *MockService) GetUserByID(ctx context.Context, userID int64) (repository.User, error) {
	return m.GetUserByIDFunc(ctx, userID)
}

func (m *MockService) RevokeRefreshToken(ctx context.Context, arg repository.RevokeRefreshTokenParams) (sql.Result, error) {
	return m.RevokeRefreshTokenFunc(ctx, arg)
}

func (m *MockService) GetGenre(ctx context.Context, genreName string) (repository.Genre, error) {
	return m.GetGenreFunc(ctx, genreName)
}

func TestRegisterValidation(t *testing.T) {

	type Body struct {
		Username  string
		Email     string
		Password  string
		FavGenres []string
	}

	Validation := []struct {
		name     string
		errorMsg error
		body     Body
	}{
		{
			name:     "Required Error",
			errorMsg: utils.AllFieldsRequiredError,
			body: Body{
				Username:  "",
				Email:     "bildadsimiyu6@gmail.com",
				Password:  "@Admin123",
				FavGenres: []string{"fantasy", "horror"},
			},
		},

		{
			name:     "Invalid Email",
			errorMsg: utils.InvalidEmailFormat,
			body: Body{
				Username:  "Rawbil",
				Email:     "bildadsimiyu6@gmail",
				Password:  "@Admin123",
				FavGenres: []string{"fantasy", "horror"},
			},
		},

		{
			name:     "Invalid Password Format",
			errorMsg: utils.InvalidPasswordFormat,
			body: Body{
				Username:  "Rawbil",
				Email:     "bildadsimiyu6@gmail.com",
				Password:  "@Admin",
				FavGenres: []string{"fantasy", "horror"},
			},
		},
		{
			name:     "Min Username",
			errorMsg: utils.MinTitleError,
			body: Body{
				Username:  "Ra",
				Email:     "bildadsimiyu6@gmail.com",
				Password:  "@Admin123",
				FavGenres: []string{"fantasy", "horror"},
			},
		},

		{
			name:     "Max Username",
			errorMsg: utils.MaxTitleError,
			body: Body{
				Username:  "hewhomustnotbenamed254001LordVoldemort",
				Email:     "bildadsimiyu6@gmail.com",
				Password:  "@Admin123",
				FavGenres: []string{"fantasy", "horror"},
			},
		},

		{
			name:     "Genre Required",
			errorMsg: utils.GenreMissing,
			body: Body{
				Username:  "Rawbil",
				Email:     "bildadsimiyu6@gmail.com",
				Password:  "@Admin123",
				FavGenres: []string{},
			},
		},

		{
			name:     "No Empty Genre",
			errorMsg: utils.NoEmptyGenre,
			body: Body{
				Username:  "Rawbil",
				Email:     "bildadsimiyu6@gmail.com",
				Password:  "@Admin123",
				FavGenres: []string{"fantasy", "horror", ""},
			},
		},
	}

	for _, test := range Validation {
		t.Run(test.name, func(t *testing.T) {

			mock := &MockService{}

			svc := NewService(mock, &sql.DB{})

			err := svc.RegisterUser(context.Background(), utils.CreateUserParams{
				Username:  test.body.Username,
				Email:     test.body.Email,
				Password:  test.body.Password,
				FavGenres: test.body.FavGenres,
			})

			if err != test.errorMsg {
				t.Fatal(err)
			}
		})
	}
}

func TestUserExistsError(t *testing.T) {
	mock := &MockService{
		GetUserByEmailFunc: func(ctx context.Context, email string) (repository.GetUserByEmailRow, error) {
			return repository.GetUserByEmailRow{}, utils.DuplicateRecordError
		},
	}

	svc := NewService(mock, &sql.DB{})

	err := svc.RegisterUser(context.Background(), utils.CreateUserParams{
		Username:  "Rawbil",
		Email:     "testuser@gmail.com",
		Password:  "@Admin123",
		FavGenres: []string{"drama", "fiction"},
	})

	if err != utils.DuplicateRecordError {
		t.Fatal(err)
	}

}

func TestSuccessfulRegistration(t *testing.T) {
	mock := &MockService{
		GetUserByEmailFunc: func(ctx context.Context, email string) (repository.GetUserByEmailRow, error) {
			return repository.GetUserByEmailRow{}, sql.ErrNoRows
		},

		CreateUserFunc: func(ctx context.Context, arg repository.CreateUserParams) (sql.Result, error) {
			return nil, nil
		},

		GetRoleFunc: func(ctx context.Context, role string) (repository.Role, error) {
			return repository.Role{}, nil
		},

		CreateUserRolesFunc: func(ctx context.Context, arg repository.CreateUserRolesParams) (sql.Result, error) {
			return nil, nil
		},

		GetGenreFunc: func(ctx context.Context, genreName string) (repository.Genre, error) {
			return repository.Genre{}, nil
		},

		CreateGenreFunc: func(ctx context.Context, genreName string) (sql.Result, error) {
			return nil, nil
		},

		CreateUserGenreFunc: func(ctx context.Context, arg repository.CreateUserGenreParams) (sql.Result, error) {
			return nil, nil
		},
	}

	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_DSN to run the MySQL registration integration test")
	}
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { db.Close() })

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	svc := NewService(mock, db)

	if err := svc.RegisterUser(context.Background(), utils.CreateUserParams{
		Username:  "Rawbil",
		Email:     fmt.Sprintf("registration-%d@example.test", time.Now().UnixNano()),
		Password:  "@Admin123",
		FavGenres: []string{"fantasy", "horror"},
	}); err != nil {
		t.Fatal(err)
	}
}
