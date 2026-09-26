package movies

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	repository "github.com/rawbil/movie-stream-app/internal/adapters/sqlc"
	"github.com/rawbil/movie-stream-app/internal/utils"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type Service interface {
	CreateGenre(ctx context.Context, genreName string) (string, error)
	UpdateGenre(ctx context.Context, arg utils.UpdateGenreParams) (string, error)
	ListGenres(ctx context.Context) ([]repository.Genre, error)
	ListMovies(ctx context.Context) ([]repository.ListMoviesRow, error)
	GetMovie(ctx context.Context, publicID uuid.UUID) (repository.GetMovieRow, error)
	CreateMovie(ctx context.Context, arg utils.CreateMovieParams) error
	AddRankings(ctx context.Context, arg utils.AddRankingsParams) (error, string)
	AddReview(ctx context.Context, publicID uuid.UUID, arg utils.AddReviewParams) error
}

type Svc struct {
	repository repository.Queries
	db         *sql.DB
}

func NewService(repo repository.Queries, db *sql.DB) Service {
	return &Svc{
		repository: repo,
		db:         db,
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

// ! List All Movies
func (svc *Svc) ListMovies(ctx context.Context) ([]repository.ListMoviesRow, error) {
	return svc.repository.ListMovies(ctx)
}

// ! Get Movie
func (svc *Svc) GetMovie(ctx context.Context, publicID uuid.UUID) (repository.GetMovieRow, error) {

	movie, err := svc.repository.GetMovie(ctx, publicID[:])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.GetMovieRow{}, utils.NoRecordError
		}
		return repository.GetMovieRow{}, err
	}

	return movie, nil
}

// ! Create Movie
func (svc *Svc) CreateMovie(ctx context.Context, arg utils.CreateMovieParams) error {
	//~ Validate fields
	if err := utils.ValidateCreateMovie(repository.CreateMovieParams{
		ImdbID:     arg.ImdbID,
		Title:      arg.Title,
		PosterPath: arg.PosterPath,
	}); err != nil {
		if utils.ValidationErrors("required", err) {
			return utils.AllFieldsRequiredError
		}

		if utils.ValidationErrors("min", err) {
			return utils.MinTitleError
		}

		if utils.ValidationErrors("max", err) {
			return utils.MaxTitleError
		}

		if utils.ValidationErrors("url", err) {
			return utils.InvalidUrlError
		}

		return err
	}

	if len(arg.Genres) < 1 {
		return utils.GenreMissing
	}

	public_id := uuid.New()
	// fmt.Println(public_id)

	//? Start transaction for creating movie and genre
	tx, err := svc.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	qtx := svc.repository.WithTx(tx)

	//~ Ensure movie does not exist
	if _, err := qtx.GetUniqueMovie(ctx, arg.ImdbID); err == nil {
		return utils.MovieExistsError
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	//~ Create movie
	movie_result, err := qtx.CreateMovie(ctx, repository.CreateMovieParams{
		PublicID:   public_id[:], // gives the underlying 16-bytes value from the []byte
		ImdbID:     arg.ImdbID,
		Title:      arg.Title,
		PosterPath: arg.PosterPath,
		YoutubeID: sql.NullString{
			String: arg.YoutubeID,
			Valid:  true,
		},
	})
	if err != nil {
		return err
	}

	movie_id, err := movie_result.LastInsertId()
	if err != nil {
		return err
	}

	for _, genreName := range arg.Genres {
		//~ Get genre
		genreName = strings.ToLower(strings.TrimSpace(genreName))
		genre, err := qtx.GetGenre(ctx, genreName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				//~ Create genre if not found
				new_genre, err := qtx.CreateGenre(ctx, strings.ToLower(genreName))
				if err != nil {
					return err
				}

				//~ get id
				genre_id, err := new_genre.LastInsertId()
				if err != nil {
					return err
				}
				genre.GenreID = genre_id

			} else {
				return err
			}

		}

		//~ Create movie genre
		if _, err := qtx.CreateMovieGenre(ctx, repository.CreateMovieGenreParams{
			MovieID: movie_id,
			GenreID: genre.GenreID,
		}); err != nil {
			return err
		}
	}

	//~ Commit context
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}

// ! Add Rankings
func (svc *Svc) AddRankings(ctx context.Context, arg utils.AddRankingsParams) (error, string) {
	//~Validate fields
	if err := utils.ValidateAddRankings(arg); err != nil {
		if utils.ValidationErrors("required", err) {
			return utils.AllFieldsRequiredError, ""
		}

		if utils.ValidationErrors("min", err) {
			return utils.GenreMissing, ""
		}

		return err, ""
	}

	for _, ranking := range arg.Rankings {
		//~ Ensure ranking does not exist
		if _, err := svc.repository.GetRanking(ctx, ranking.RankingName); err == nil {
			return utils.MovieExistsError, ranking.RankingName
		} else if !errors.Is(err, sql.ErrNoRows) {
			return err, ""
		}

		//~ Add ranking
		if _, err := svc.repository.AddRanking(ctx, repository.AddRankingParams{
			RankingName:  ranking.RankingName,
			RankingValue: ranking.RankingValue,
		}); err != nil {
			return err, ""
		}
	}

	return nil, ""
}

// ! Add Review
// Add a review, and run all the reviews through AI and return the final ranking
func (svc *Svc) AddReview(ctx context.Context, publicID uuid.UUID, arg utils.AddReviewParams) error {
	//~ Validate field
	if err := utils.ValidateAddReview(arg); err != nil {
		if utils.ValidationErrors("required", err) {
			return utils.AllFieldsRequiredError
		}
		return err
	}

	//~ Ensure movie exists

	movie, err := svc.repository.GetMovieInternal(ctx, publicID[:])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NoRecordError
		}
		return err
	}

	//~ Get Rankings
	rankings, err := svc.repository.GetRankings(ctx)
	if err != nil {
		return err
	}

	if len(rankings) < 1 {
		return utils.NoRecordError
	}

	var ranking_names string

	for _, ranking := range rankings {
		ranking_names += ranking.RankingName + ","
	}

	//~ Trim off the last comma
	ranking_names = strings.Trim(ranking_names, ",")

	tx, err := svc.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	qtx := svc.repository.WithTx(tx)

	//~ Add review
	if _, err := qtx.AddMovieReview(ctx, repository.AddMovieReviewParams{
		MovieID: movie.MovieID,
		Review:  arg.Review,
	}); err != nil {
		return err
	}

	//~ Get Reviews
	reviews, err := qtx.GetMovieReviews(ctx, movie.MovieID)
	if err != nil {
		return err
	}

	var all_reviews string

	for _, r := range reviews[:min(5, len(reviews))] {
		all_reviews += r.Review + ","
	}

	all_reviews = strings.Trim(all_reviews, ",")

	PROMPT_MESSAGE := fmt.Sprintf("Return a response using one of these words: %s. The response should be a single word and should not contain any other text. The response should be based on the following reviews: %s", ranking_names, all_reviews)

	groq_api_key := utils.ServerConfigFunc().GroqApiKey
	if groq_api_key == "" {
		return errors.New("Groq API Key missing")
	}

	llm, err := openai.New(
		openai.WithModel("openai/gpt-oss-120b"),
		openai.WithBaseURL("https://api.groq.com/openai/v1"),
		openai.WithToken(groq_api_key),
	)

	if err != nil {
		if llms.IsAuthenticationError(err) {
			return utils.InvalidAPiKey
		}
		if llms.IsRateLimitError(err) {
			return utils.GroqApiLimit
		}
		return err
	}

	llm_response, err := llm.GenerateContent(ctx, []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, PROMPT_MESSAGE),
	})
	if err != nil {
		return err
	}

	var ranking_value int32

	for _, ranking := range rankings {
		if ranking.RankingName == llm_response.Choices[0].Content {
			ranking_value = (ranking.RankingValue)
		}
	}

	if _, err := qtx.UpdateMovieRankings(ctx, repository.UpdateMovieRankingsParams{
		PublicID: publicID[:],
		RankingName: sql.NullString{
			String: llm_response.Choices[0].Content,
			Valid:  true,
		},
		RankingValue: sql.NullInt32{
			Int32: ranking_value,
			Valid: true,
		},
	}); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
