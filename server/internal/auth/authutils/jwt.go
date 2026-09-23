package authutils

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	repository "github.com/rawbil/movie-stream-app/internal/adapters/sqlc"
	"github.com/rawbil/movie-stream-app/internal/utils"
)

type Claims struct {
	UserID   int64  `json:"user_id"`
	TokenUse string `json:"token_use"`

	jwt.RegisteredClaims
}

type ContextKey string

var userIDContextKey ContextKey

func GenerateAuthToken(user_id int64, secret string) (string, string, time.Time, time.Time, error) {
	issued_at := time.Now()
	at_expires_at := issued_at.Add(time.Second * time.Duration(3600)) // 1 hour
	rt_expires_at := issued_at.Add(time.Hour * time.Duration(24) * time.Duration(7))

	at_claims := &Claims{
		UserID:   user_id,
		TokenUse: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issued_at),
			ExpiresAt: jwt.NewNumericDate(at_expires_at),
		},
	}

	rt_claims := &Claims{
		UserID:   user_id,
		TokenUse: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issued_at),
			ExpiresAt: jwt.NewNumericDate(rt_expires_at),
		},
	}
	// each claim is a new variable, meaning they will point to different memory locations.

	access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, at_claims)
	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, rt_claims)

	access_token_string, err := access_token.SignedString([]byte(secret))
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	refresh_token_string, err := refresh_token.SignedString([]byte(secret))
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	return access_token_string, refresh_token_string, issued_at, rt_expires_at, nil
}

// ! AuthMiddleware
func AuthMiddleware(repository repository.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "auth header missing", errors.New("auth header missing"))
			c.Abort()
			return
		}

		bearer_token := strings.SplitN(auth, " ", 2)
		if len(bearer_token) != 2 || bearer_token[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "token misconfigured", errors.New("token misconfigured"))
			c.Abort()
			return
		}

		token := bearer_token[1]

		//~ Validate token
		claims, err := ValidateAccessToken(token)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), err)
			c.Abort()
			return
		}

		//~ Find user using the stored claims
		user, err := repository.GetUserByID(c.Request.Context(), claims.UserID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				utils.ErrorResponse(c, http.StatusUnauthorized, "user from claim not found", errors.New("user from claim not found"))
				c.Abort()
				return
			}
			utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), err)
			c.Abort()
			return
		}

		// ctx := context.WithValue(c.Request.Context(), userIDContextKey, user.UserID)
		// c.Request = c.Request.WithContext(ctx)
		c.Set("user_id", user.UserID)

		c.Next()
	}
}

func ValidateAccessToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid token signing method")
		}

		return []byte(utils.ServerConfigFunc().JwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.TokenUse != "access" {
		return nil, errors.New("invalid access token")
	}

	return claims, nil
}

func ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid token signing method")
		}

		return []byte(utils.ServerConfigFunc().JwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.TokenUse != "refresh" {
		return nil, errors.New("invalid refresh token")
	}

	return claims, nil
}

func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	user_id, ok := ctx.Value(userIDContextKey).(int64)

	return user_id, ok
}
