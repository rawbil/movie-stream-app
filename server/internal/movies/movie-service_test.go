package movies

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"testing"

	"github.com/google/uuid"
	repository "github.com/rawbil/movie-stream-app/internal/adapters/sqlc"
	"github.com/rawbil/movie-stream-app/internal/utils"
)

var errMovieServiceTestDB = errors.New("mock database failure")

type movieServiceTestDriver struct{}
type movieServiceTestConn struct{}
type movieServiceTestTx struct{}

func (movieServiceTestDriver) Open(string) (driver.Conn, error) { return movieServiceTestConn{}, nil }
func (movieServiceTestConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepare is unsupported")
}
func (movieServiceTestConn) Close() error              { return nil }
func (movieServiceTestConn) Begin() (driver.Tx, error) { return movieServiceTestTx{}, nil }
func (movieServiceTestConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	return movieServiceTestTx{}, nil
}
func (movieServiceTestConn) ExecContext(context.Context, string, []driver.NamedValue) (driver.Result, error) {
	return nil, errMovieServiceTestDB
}
func (movieServiceTestConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	return nil, errMovieServiceTestDB
}
func (movieServiceTestTx) Commit() error   { return nil }
func (movieServiceTestTx) Rollback() error { return nil }

type movieServiceTestRows struct{}

func (movieServiceTestRows) Columns() []string         { return nil }
func (movieServiceTestRows) Close() error              { return nil }
func (movieServiceTestRows) Next([]driver.Value) error { return io.EOF }

func TestMovieServiceMethodsReturnRepositoryErrors(t *testing.T) {
	sql.Register("movie-service-test", movieServiceTestDriver{})
	db, err := sql.Open("movie-service-test", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	queries := repository.New(db)
	svc := NewService(*queries, db)
	ctx := context.Background()
	publicID := uuid.New()

	cases := []struct {
		name string
		call func() error
	}{
		{"create genre", func() error { _, err := svc.CreateGenre(ctx, "drama"); return err }},
		{"update genre", func() error {
			_, err := svc.UpdateGenre(ctx, utils.UpdateGenreParams{OldGenre: "comedy", CurrentGenre: "drama"})
			return err
		}},
		{"list genres", func() error { _, err := svc.ListGenres(ctx); return err }},
		{"list movies", func() error { _, _, err := svc.ListMovies(ctx, repository.ListMoviesParams{Limit: 10}); return err }},
		{"get movie", func() error { _, err := svc.GetMovie(ctx, publicID); return err }},
		{"create movie", func() error {
			return svc.CreateMovie(ctx, utils.CreateMovieParams{ImdbID: "tt123", Title: "A film", PosterPath: "https://example.test/poster.jpg", Genres: []string{"drama"}})
		}},
		{"add rankings", func() error {
			err, _ := svc.AddRankings(ctx, utils.AddRankingsParams{Rankings: []utils.RankingParams{{RankingName: "good", RankingValue: 5}}})
			return err
		}},
		{"add review", func() error { return svc.AddReview(ctx, publicID, 1, utils.AddReviewParams{Review: "A good film"}) }},
		{"recommendations", func() error { _, err := svc.GetMovieRecommendations(ctx, 1); return err }},
		{"get reviews", func() error {
			_, _, err := svc.GetMovieReviews(ctx, repository.GetMovieReviewsParams{PublicID: publicID[:]})
			return err
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.call(); !errors.Is(err, errMovieServiceTestDB) {
				t.Fatalf("err=%v, want repository error", err)
			}
		})
	}
}
