package auth

import (
	"context"
	"database/sql"
	"errors"

	repository "github.com/rawbil/movie-stream-app/internal/adapters/sqlc"
	"github.com/rawbil/movie-stream-app/internal/auth/authutils"
	"github.com/rawbil/movie-stream-app/internal/authorization"
	"github.com/rawbil/movie-stream-app/internal/utils"
)

type Service interface {
	CreateRole(ctx context.Context, role string) (sql.Result, error)
}

type Svc struct {
	repository repository.Queries
	db         *sql.DB
}

func NewService(repository repository.Queries, db *sql.DB) Service {
	return &Svc{
		repository: repository,
		db:         db,
	}
}

// ! Register
func (svc *Svc) RegisterUser(ctx context.Context, arg utils.CreateUserParams) (error) {
	//~ Validate fields
	if err := utils.ValidateCreateUser(repository.CreateUserParams{
		Username: arg.Username,
		Email:    arg.Email,
		Password: arg.Password,
	}); err != nil {
		if utils.ValidationErrors("required", err) {
			return utils.AllFieldsRequiredError
		}

		if utils.ValidationErrors("email", err) {
			return utils.InvalidEmailFormat
		}

		if utils.ValidationErrors("password_format", err) {
			return utils.InvalidPasswordFormat
		}

		if utils.ValidationErrors("min", err) {
			return utils.MinTitleError
		}

		if utils.ValidationErrors("max", err) {
			return utils.MaxTitleError
		}
		return err
	}

	//~ Ensure user does not already exist
	if _, err := svc.repository.GetUserByEmail(ctx, arg.Email); err == nil {
		return utils.DuplicateRecordError
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	//~ Hash password
	hashed_password, err := authutils.HashPassword(arg.Password)
	if err != nil {
		return err
	}

	//? TRANSACTION FOR CREATING USER, UPDATING USER_ROLES TABLE AND FAV_GENRES TABLE
	tx, err := svc.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	qtx := svc.repository.WithTx(tx)

	//~ Save user
	create_user_result, err := qtx.CreateUser(ctx, repository.CreateUserParams{
		Username: arg.Username,
		Email:    arg.Email,
		Password: hashed_password,
	})
	if err != nil {
		return err
	}

	user_id, err := create_user_result.LastInsertId()
	if err != nil {
		return err
	}

	//~ Get Default Role
	role, err := qtx.GetRole(ctx, authorization.UserRole)
	if err != nil {
		// create role if not available
		if errors.Is(err, sql.ErrNoRows) {
			if _, err := qtx.CreateRole(ctx, authorization.UserRole); err != nil {
				return err
			}
		}
		return err
	}

	//~ Update User Role
	if _, err := qtx.CreateUserRoles(ctx, repository.CreateUserRolesParams{
		UserID: user_id,
		RoleID: role.RoleID,
	}); err != nil {
		return err
	}

	//~ Create user favourite genres
	for _, genre := range arg.FavGenres {
		// Get genre details
		genre_record, err := qtx.GetGenre(ctx, genre)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// create genre if it is missing
				created_genre, err := qtx.CreateGenre(ctx, genre)
				if err != nil {
					return err
				}

				created_genre_id, err := created_genre.LastInsertId()
				if err != nil {
					return err
				}

				genre_record.GenreID = created_genre_id

			}
			return err
		}

		if _, err := qtx.CreateUserGenre(ctx, repository.CreateUserGenreParams{
			UserID:  user_id,
			GenreID: genre_record.GenreID,
		}); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// ! CreateRole
func (svc *Svc) CreateRole(ctx context.Context, role string) (sql.Result, error) {
	//~ Ensure role is provided
	if role == "" {
		return nil, utils.AllFieldsRequiredError
	}

	//~ Ensure role is not already available
	if _, err := svc.repository.GetRole(ctx, role); err != nil {
		return nil, utils.DuplicateRecordError
	}

	return svc.repository.CreateRole(ctx, role)
}
