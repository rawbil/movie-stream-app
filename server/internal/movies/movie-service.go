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
	UpdateGenre(ctx context.Context, arg utils.UpdateGenreParams) (string, error)
	ListGenres(ctx context.Context) ([]repository.Genre, error)
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

// ! Update Genre
func (svc *Svc) UpdateGenre(ctx context.Context, arg utils.UpdateGenreParams) (string, error) {
	//~ Validate fields
	if err := utils.ValidateUpdateGenre(arg); err != nil {
		if utils.ValidationErrors("required", err) {
			return "", utils.AllFieldsRequiredError
		}
		return "", err
	}

	//~ Ensure genre exists
	genre, err := svc.repository.GetGenre(ctx, arg.OldGenre)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", utils.NoRecordError
		}
		return "", err
	}

	//~ Ensure new genre does not exist
	_, Geterr := svc.repository.GetGenre(ctx, arg.CurrentGenre)
	if Geterr == nil {
		return "", utils.DuplicateRecordError
	} else if !errors.Is(Geterr, sql.ErrNoRows) {
		return "", Geterr
	}

	//~ Update genre record
	if _, err := svc.repository.UpdateGenre(ctx, repository.UpdateGenreParams{
		GenreName: arg.CurrentGenre,
		GenreID:   genre.GenreID,
	}); err != nil {
		return "", err
	}

	return arg.CurrentGenre, nil
}

// ! List All Genres
func (svc *Svc) ListGenres(ctx context.Context) ([]repository.Genre, error) {
	return svc.repository.ListGenres(ctx)
}