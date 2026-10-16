package auth

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	repository "github.com/rawbil/movie-stream-app/internal/adapters/sqlc"
	"github.com/rawbil/movie-stream-app/internal/utils"
)

type handlerServiceStub struct {
	register func(context.Context, utils.CreateUserParams) error
	login    func(context.Context, utils.LoginParams) (repository.User, string, string, error)
	role     func(context.Context, string) (sql.Result, error)
	logout   func(context.Context, int64) error
	refresh  func(context.Context, string) (string, string, repository.User, error)
}

func (s handlerServiceStub) RegisterUser(ctx context.Context, p utils.CreateUserParams) error {
	return s.register(ctx, p)
}
func (s handlerServiceStub) LoginUser(ctx context.Context, p utils.LoginParams) (repository.User, string, string, error) {
	return s.login(ctx, p)
}
func (s handlerServiceStub) CreateRole(ctx context.Context, role string) (sql.Result, error) {
	return s.role(ctx, role)
}
func (s handlerServiceStub) Logout(ctx context.Context, id int64) error { return s.logout(ctx, id) }
func (s handlerServiceStub) RefreshTokens(ctx context.Context, token string) (string, string, repository.User, error) {
	return s.refresh(ctx, token)
}

func performAuthHandler(h gin.HandlerFunc, method, path, body string, userID any, cookies ...*http.Cookie) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Handle(method, path, func(c *gin.Context) {
		if userID != nil {
			c.Set("user_id", userID)
		}
		h(c)
	})
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestAuthHandlers(t *testing.T) {
	t.Run("register success and invalid body", func(t *testing.T) {
		called := false
		h := NewHandler(handlerServiceStub{register: func(context.Context, utils.CreateUserParams) error { called = true; return nil }})
		if got := performAuthHandler(h.RegisterUser, http.MethodPost, "/", `{"username":"sam","email":"sam@example.com","password":"pass","fav_genres":["drama"]}`, nil).Code; got != http.StatusOK || !called {
			t.Fatalf("status=%d called=%v", got, called)
		}
		if got := performAuthHandler(h.RegisterUser, http.MethodPost, "/", `{`, nil).Code; got != http.StatusBadRequest {
			t.Fatalf("invalid body status=%d", got)
		}
	})

	t.Run("login success and service error", func(t *testing.T) {
		h := NewHandler(handlerServiceStub{login: func(context.Context, utils.LoginParams) (repository.User, string, string, error) {
			return repository.User{Email: "sam@example.com"}, "access", "refresh", nil
		}})
		w := performAuthHandler(h.LoginUser, http.MethodPost, "/", `{"email":"sam@example.com","password":"pass"}`, nil)
		if w.Code != http.StatusOK || len(w.Result().Cookies()) != 1 {
			t.Fatalf("status=%d cookies=%v", w.Code, w.Result().Cookies())
		}
		h = NewHandler(handlerServiceStub{login: func(context.Context, utils.LoginParams) (repository.User, string, string, error) {
			return repository.User{}, "", "", utils.NoRecordError
		}})
		if got := performAuthHandler(h.LoginUser, http.MethodPost, "/", `{"email":"sam@example.com","password":"pass"}`, nil).Code; got != http.StatusBadRequest {
			t.Fatalf("service error status=%d", got)
		}
	})

	t.Run("create role success and conflict", func(t *testing.T) {
		h := NewHandler(handlerServiceStub{role: func(context.Context, string) (sql.Result, error) { return nil, nil }})
		if got := performAuthHandler(h.CreateRole, http.MethodPost, "/", `{"role":"admin"}`, nil).Code; got != http.StatusOK {
			t.Fatalf("success status=%d", got)
		}
		h = NewHandler(handlerServiceStub{role: func(context.Context, string) (sql.Result, error) { return nil, utils.DuplicateRecordError }})
		if got := performAuthHandler(h.CreateRole, http.MethodPost, "/", `{"role":"admin"}`, nil).Code; got != http.StatusConflict {
			t.Fatalf("duplicate status=%d", got)
		}
	})

	t.Run("logout context and success", func(t *testing.T) {
		h := NewHandler(handlerServiceStub{logout: func(context.Context, int64) error { return nil }})
		if got := performAuthHandler(h.Logout, http.MethodPost, "/", "", nil).Code; got != http.StatusNotFound {
			t.Fatalf("missing context status=%d", got)
		}
		w := performAuthHandler(h.Logout, http.MethodPost, "/", "", int64(7))
		if w.Code != http.StatusOK || len(w.Result().Cookies()) != 1 || w.Result().Cookies()[0].MaxAge >= 0 {
			t.Fatalf("logout response status=%d cookies=%v", w.Code, w.Result().Cookies())
		}
	})

	t.Run("refresh token missing and success", func(t *testing.T) {
		h := NewHandler(handlerServiceStub{refresh: func(context.Context, string) (string, string, repository.User, error) {
			return "access", "new-refresh", repository.User{}, nil
		}})
		if got := performAuthHandler(h.RefreshTokens, http.MethodPost, "/", "", nil).Code; got != http.StatusUnauthorized {
			t.Fatalf("missing cookie status=%d", got)
		}
		settings := getRefreshCookieSettings()
		cookie := &http.Cookie{Name: settings.name, Value: "refresh"}
		w := performAuthHandler(h.RefreshTokens, http.MethodPost, "/", "", nil, cookie)
		if w.Code != http.StatusOK || len(w.Result().Cookies()) != 1 {
			t.Fatalf("refresh status=%d cookies=%v", w.Code, w.Result().Cookies())
		}
		h = NewHandler(handlerServiceStub{refresh: func(context.Context, string) (string, string, repository.User, error) {
			return "", "", repository.User{}, errors.New("failure")
		}})
		if got := performAuthHandler(h.RefreshTokens, http.MethodPost, "/", "", nil, cookie).Code; got != http.StatusUnauthorized {
			t.Fatalf("service error status=%d", got)
		}
	})
}
