package movies

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	repository "github.com/rawbil/movie-stream-app/internal/adapters/sqlc"
	"github.com/rawbil/movie-stream-app/internal/utils"
)

type movieHandlerServiceStub struct{}

func (movieHandlerServiceStub) CreateGenre(context.Context, string) (string, error) {
	return "drama", nil
}
func (movieHandlerServiceStub) UpdateGenre(context.Context, utils.UpdateGenreParams) (string, error) {
	return "drama", nil
}
func (movieHandlerServiceStub) ListGenres(context.Context) ([]repository.Genre, error) {
	return []repository.Genre{}, nil
}
func (movieHandlerServiceStub) ListMovies(context.Context, repository.ListMoviesParams) ([]repository.ListMoviesRow, int64, error) {
	return []repository.ListMoviesRow{}, 0, nil
}
func (movieHandlerServiceStub) GetMovie(context.Context, uuid.UUID) (repository.GetMovieRow, error) {
	return repository.GetMovieRow{}, nil
}
func (movieHandlerServiceStub) CreateMovie(context.Context, utils.CreateMovieParams) error {
	return nil
}
func (movieHandlerServiceStub) UpdateMovie(context.Context, utils.UpdateMovieParams) error {
	return nil
}
func (movieHandlerServiceStub) DeleteMovie(context.Context, uuid.UUID) error {
	return nil
}
func (movieHandlerServiceStub) AddRankings(context.Context, utils.AddRankingsParams) (error, string) {
	return nil, ""
}
func (movieHandlerServiceStub) AddReview(context.Context, uuid.UUID, int64, utils.AddReviewParams) error {
	return nil
}
func (movieHandlerServiceStub) GetMovieRecommendations(context.Context, int64) ([]repository.GetRecommendedMoviesRow, error) {
	return []repository.GetRecommendedMoviesRow{}, nil
}
func (movieHandlerServiceStub) GetMovieReviews(context.Context, repository.GetMovieReviewsParams) ([]repository.GetMovieReviewsRow, int64, error) {
	return []repository.GetMovieReviewsRow{}, 0, nil
}

func performMovieHandler(h gin.HandlerFunc, method, path, body string, userID any) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	routePath := strings.SplitN(path, "?", 2)[0]
	if routePath != "/" {
		routePath = "/:public_id"
	}
	r.Handle(method, routePath, func(c *gin.Context) {
		if userID != nil {
			c.Set("user_id", userID)
		}
		h(c)
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(method, path, strings.NewReader(body)))
	return w
}

func TestMovieHandlers(t *testing.T) {
	h := NewHandler(movieHandlerServiceStub{})
	id := uuid.NewString()
	cases := []struct {
		name, method, path, body string
		handler                  gin.HandlerFunc
		userID                   any
		want                     int
	}{
		{"create genre", http.MethodPost, "/", `{"current_genre":"drama"}`, h.CreateGenre, nil, http.StatusOK},
		{"update genre", http.MethodPatch, "/", `{"old_genre":"comedy","current_genre":"drama"}`, h.UpdateGenre, nil, http.StatusOK},
		{"list genres", http.MethodGet, "/", "", h.ListGenres, nil, http.StatusOK},
		{"list movies", http.MethodGet, "/?page=2", "", h.ListMovies, nil, http.StatusOK},
		{"get movie", http.MethodGet, "/?public_id=" + id, "", h.GetMovie, nil, http.StatusOK},
		{"create movie", http.MethodPost, "/", `{"imdb_id":"tt123","title":"A film","poster_path":"https://example.com/poster.jpg","genres":["drama"]}`, h.CreateMovie, nil, http.StatusOK},
		{"update movie", http.MethodPut, "/", `{"public_id":"` + id + `","imdb_id":"tt123","title":"A film","poster_path":"https://example.com/poster.jpg","genres":["drama"]}`, h.UpdateMovie, nil, http.StatusOK},
		{"delete movie", http.MethodDelete, "/" + id, "", h.DeleteMovie, nil, http.StatusOK},
		{"add rankings", http.MethodPost, "/", `{"rankings":[{"ranking_name":"great","ranking_value":5}]}`, h.AddRankings, nil, http.StatusOK},
		{"add review", http.MethodPost, "/" + id, `{"review":"A thoughtful review"}`, h.AddMovieReview, int64(42), http.StatusOK},
		{"recommendations", http.MethodGet, "/", "", h.GetMovieRecommendations, int64(42), http.StatusOK},
		{"get reviews", http.MethodGet, "/" + id + "?page=2", "", h.GetMovieReviews, nil, http.StatusOK},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := performMovieHandler(tc.handler, tc.method, tc.path, tc.body, tc.userID)
			if w.Code != tc.want {
				t.Fatalf("status=%d want=%d body=%s", w.Code, tc.want, w.Body.String())
			}
		})
	}
}

func TestMovieHandlerRejectsInvalidInputs(t *testing.T) {
	h := NewHandler(movieHandlerServiceStub{})
	for _, tc := range []struct {
		name, method, path string
		handler            gin.HandlerFunc
		userID             any
	}{
		{"bad movie uuid", http.MethodGet, "/?public_id=bad", h.GetMovie, nil},
		{"bad update movie body", http.MethodPut, "/", h.UpdateMovie, nil},
		{"bad delete movie uuid", http.MethodDelete, "/bad", h.DeleteMovie, nil},
		{"bad movies page", http.MethodGet, "/?page=abc", h.ListMovies, nil},
		{"bad reviews page", http.MethodGet, "/bad?page=abc", h.GetMovieReviews, nil},
		{"review without user", http.MethodPost, "/" + uuid.NewString(), h.AddMovieReview, nil},
		{"recommendations without user", http.MethodGet, "/", h.GetMovieRecommendations, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := performMovieHandler(tc.handler, tc.method, tc.path, `{}`, tc.userID)
			if w.Code != http.StatusBadRequest && w.Code != http.StatusNotFound {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
		})
	}
}

type failingMovieHandlerStub struct{ movieHandlerServiceStub }

func (failingMovieHandlerStub) CreateGenre(context.Context, string) (string, error) {
	return "", utils.DuplicateRecordError
}
func (failingMovieHandlerStub) UpdateGenre(context.Context, utils.UpdateGenreParams) (string, error) {
	return "", utils.NoRecordError
}
func (failingMovieHandlerStub) ListGenres(context.Context) ([]repository.Genre, error) {
	return nil, errors.New("db")
}

func TestMovieHandlerServiceErrors(t *testing.T) {
	h := NewHandler(failingMovieHandlerStub{})
	for _, tc := range []struct {
		name, method string
		handler      gin.HandlerFunc
		body         string
		want         int
	}{
		{"genre conflict", http.MethodPost, h.CreateGenre, `{"current_genre":"drama"}`, http.StatusConflict},
		{"genre missing", http.MethodPatch, h.UpdateGenre, `{"old_genre":"old","current_genre":"new"}`, http.StatusBadRequest},
		{"genre list failure", http.MethodGet, h.ListGenres, "", http.StatusInternalServerError},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := performMovieHandler(tc.handler, tc.method, "/", tc.body, nil)
			if w.Code != tc.want {
				t.Fatalf("status=%d want=%d", w.Code, tc.want)
			}
		})
	}
}
