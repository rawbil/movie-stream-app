package movies

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	repository "github.com/rawbil/movie-stream-app/internal/adapters/sqlc"
	"github.com/rawbil/movie-stream-app/internal/utils"
)

type Service interface {
	CreateGenre(ctx context.Context, genreName string) (string, error)
}

type Svc struct {
	repository repository.Queries
}

func NewService(repo repository.Queries) Service {
	return &Svc{
		repository: repo,
	}
}

// ! Create Genre
func (svc *Svc) CreateGenre(ctx context.Context, genreName string) (string, error) {
	//~ ensure generName is provided
	if genreName == "" {
		return "", utils.AllFieldsRequiredError
	}

	//~ ensure genre does not already exist
	_, err := svc.repository.GetGenre(ctx, genreName)
	if err == nil {
		return "", utils.DuplicateRecordError
	} else if !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	genreRecord, err := svc.repository.CreateGenre(ctx, strings.ToLower(genreName))
	if err != nil {
		return "", err
	}

	genre_id, err := genreRecord.LastInsertId()
	if err != nil {
		return "", utils.NoRecordError
	}

	genre, err := svc.repository.GetGenreByID(ctx, genre_id)
	if err != nil {
		return "", err
	}

	return genre.GenreName, err
}
