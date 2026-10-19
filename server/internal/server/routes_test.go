package server

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type routeHealthDriver struct{}
type routeHealthConn struct{}

func (routeHealthDriver) Open(string) (driver.Conn, error) { return routeHealthConn{}, nil }
func (routeHealthConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepare unsupported")
}
func (routeHealthConn) Close() error               { return nil }
func (routeHealthConn) Begin() (driver.Tx, error)  { return nil, errors.New("transaction unsupported") }
func (routeHealthConn) Ping(context.Context) error { return nil }

func TestHealthRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sql.Register("route-health-test", routeHealthDriver{})
	db, err := sql.Open("route-health-test", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	handler := (&Api{DB: db}).Mount()
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/health", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d want=%d", w.Code, http.StatusOK)
	}
}

func TestMountedRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	api := &Api{DB: &sql.DB{}}
	got := api.Mount().(*gin.Engine).Routes()
	registered := make(map[string]bool, len(got))
	for _, route := range got {
		registered[route.Method+" "+route.Path] = true
	}

	want := []string{
		"GET /api/v1/health",
		"POST /api/v1/auth/register",
		"POST /api/v1/auth/login",
		"POST /api/v1/auth/logout",
		"POST /api/v1/auth/create-role",
		"POST /api/v1/auth/refresh_tokens",
		"POST /api/v1/movies/genres/add",
		"PATCH /api/v1/movies/genres/update",
		"GET /api/v1/movies/genres/list",
		"GET /api/v1/movies/all",
		"POST /api/v1/movies/add-rankings",
		"POST /api/v1/movies/add-review/:public_id",
		"GET /api/v1/movies/reviews/:public_id",
		"GET /api/v1/movies/recommended",
		"POST /api/v1/movies/create",
		"PUT /api/v1/movies/update",
		"DELETE /api/v1/movies/delete/:public_id",
		"GET /api/v1/movies/one",
	}
	if len(got) != len(want) {
		t.Fatalf("mounted %d routes, want %d", len(got), len(want))
	}
	for _, route := range want {
		if !registered[route] {
			t.Errorf("route %q is not mounted", route)
		}
	}
}

func TestProtectedRoutesRequireAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := (&Api{DB: &sql.DB{}}).Mount()
	cases := []struct{ method, path string }{
		{http.MethodPost, "/api/v1/auth/logout"},
		{http.MethodPost, "/api/v1/auth/create-role"},
		{http.MethodPost, "/api/v1/movies/genres/add"},
		{http.MethodPatch, "/api/v1/movies/genres/update"},
		{http.MethodPost, "/api/v1/movies/add-rankings"},
		{http.MethodPut, "/api/v1/movies/update"},
		{http.MethodDelete, "/api/v1/movies/delete/550e8400-e29b-41d4-a716-446655440000"},
		{http.MethodPost, "/api/v1/movies/add-review/550e8400-e29b-41d4-a716-446655440000"},
		{http.MethodGet, "/api/v1/movies/reviews/550e8400-e29b-41d4-a716-446655440000"},
		{http.MethodGet, "/api/v1/movies/recommended"},
	}
	for _, tc := range cases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, tc.path, nil)
			handler.ServeHTTP(w, req)
			if w.Code != http.StatusUnauthorized {
				t.Fatalf("status=%d want=%d body=%s", w.Code, http.StatusUnauthorized, w.Body.String())
			}
		})
	}
}
