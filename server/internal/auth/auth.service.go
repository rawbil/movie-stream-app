package auth

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	repository "github.com/rawbil/movie-stream-app/internal/adapters/sqlc"
	"github.com/rawbil/movie-stream-app/internal/auth/authutils"
	"github.com/rawbil/movie-stream-app/internal/authorization"
	"github.com/rawbil/movie-stream-app/internal/utils"
)

type Service interface {
	RegisterUser(ctx context.Context, arg utils.CreateUserParams) error
	CreateRole(ctx context.Context, role string) (sql.Result, error)
	LoginUser(ctx context.Context, arg utils.LoginParams) (repository.User, string, string, error)
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
func (svc *Svc) RegisterUser(ctx context.Context, arg utils.CreateUserParams) error {
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

	//~ Ensure fav_genres has at least one item
	if len(arg.FavGenres) < 1 {
		return utils.GenreMissing
	}

	for _, g := range arg.FavGenres {
		if strings.TrimSpace(g) == "" {
			return utils.NoEmptyGenre
		}
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
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		// create role if not available
		createdRole, err := qtx.CreateRole(ctx, authorization.UserRole)
		if err != nil {
			return err
		}

		roleID, err := createdRole.LastInsertId()
		if err != nil {
			return err
		}

		role.RoleID = roleID
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
		genre = strings.ToLower(strings.TrimSpace(genre))
		genre_record, err := qtx.GetGenre(ctx, genre)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}

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

// ! Login
func (svc *Svc) LoginUser(ctx context.Context, arg utils.LoginParams) (repository.User, string, string, error) {
	//~ Validate fields
	if err := utils.ValidateLoginUser(arg); err != nil {
		if utils.ValidationErrors("required", err) {
			return repository.User{}, "", "", utils.AllFieldsRequiredError
		}

		if utils.ValidationErrors("email", err) {
			return repository.User{}, "", "", utils.InvalidEmailFormat
		}

		return repository.User{}, "", "", err
	}

	//~ Get user with email
	user, err := svc.repository.GetUserByEmail(ctx, arg.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.User{}, "", "", utils.NoRecordError
		}
		return repository.User{}, "", "", err
	}

	//~ Compare passwords
	if err := authutils.ComparePasswords(user.Password, arg.Password); err != nil {
		return repository.User{}, "", "", utils.IncorrectPassword
	}

	//~ Generate JWT Tokens
	jwt_secret := utils.ServerConfigFunc().JwtSecret
	access_token, refresh_token, _, _, err := authutils.GenerateAuthToken(user.UserID, jwt_secret)
	if err != nil {
		return repository.User{}, "", "", err
	}

	//~ Hash refresh token
	hashed_rt := authutils.RefreshTokenHash(refresh_token)

	//~ Create refresh token record if not exist, or update if user already present
	if _, err := svc.repository.GetRefreshToken(ctx, user.UserID); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return repository.User{}, "", "", err
		}
		// if record does not exist, create it(First-time login)
		if _, err := svc.repository.CreateRefreshToken(ctx, repository.CreateRefreshTokenParams{
			UserID:      user.UserID,
			HashedToken: hashed_rt,
		}); err != nil {
			return repository.User{}, "", "", err
		}
	} else {
		// if record exists, update it
		if _, err := svc.repository.UpdateRefreshToken(ctx, repository.UpdateRefreshTokenParams{
			UserID:      user.UserID,
			HashedToken: hashed_rt,
			Revoked:     false,
		}); err != nil {
			return repository.User{}, "", "", err
		}
	}

	return user, access_token, refresh_token, nil
}

// ! CreateRole
func (svc *Svc) CreateRole(ctx context.Context, role string) (sql.Result, error) {

	role = strings.ToUpper(role)
	//~ Ensure role is provided
	if role == "" {
		return nil, utils.AllFieldsRequiredError
	}

	//~ Ensure role is not already available
	if _, err := svc.repository.GetRole(ctx, role); err == nil {
		return nil, utils.DuplicateRecordError
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return svc.repository.CreateRole(ctx, role)
}
