package authutils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
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

//! AuthMiddleware

// func AuthMiddleware(repository repository.Queries) utils.Middleware {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			//~ Get the token from auth headers
// 			auth := r.Header.Get("Authorization")

// 			if auth == "" {
// 				utils.ErrorResponse(w, "Authorization header missing", errors.New("Authorization header missing"), http.StatusUnauthorized)
// 				return
// 			}

// 			bearer_token := strings.Split(auth, " ")
// 			if len(bearer_token) != 2 || bearer_token[0] != "Bearer" {
// 				utils.ErrorResponse(w, "token malformed", errors.New("bearer token malformed"), http.StatusUnauthorized)
// 				return
// 			}

// 			auth_token := bearer_token[1]

// 			//~ Validate token
// 			claims, err := ValidateAccessToken(auth_token)
// 			if err != nil {
// 				utils.ErrorResponse(w, err.Error(), err, http.StatusUnauthorized)
// 				return
// 			}

// 			//~ Find user from the claims user
// 			user, err := repository.GetUserByID(r.Context(), claims.UserID)
// 			if err != nil {
// 				if errors.Is(err, sql.ErrNoRows) {
// 					utils.ErrorResponse(w, "decoded uuser not found. Login again", err, http.StatusNotFound)
// 					return
// 				}
// 				utils.ErrorResponse(w, err.Error(), err, http.StatusUnauthorized)
// 				return
// 			}

// 			//~ Add user to context
// 			ctx := context.WithValue(r.Context(), userIDContextKey, user.UserID)

// 			r = r.WithContext(ctx)

// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

// func ValidateAccessToken(tokenString string) (*Claims, error) {
// 	claims := &Claims{}
// 	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
// 		if token.Method != jwt.SigningMethodHS256 {
// 			return nil, errors.New("invalid token signing method")
// 		}

// 		return []byte(config.ServerConfigFunc().JwtSecret), nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	if !token.Valid {
// 		return nil, errors.New("invalid token")
// 	}

// 	if claims.TokenUse != "access" {
// 		return nil, errors.New("invalid access token")
// 	}

// 	return claims, nil
// }

// func ValidateRefreshToken(tokenString string) (*Claims, error) {
// 	claims := &Claims{}
// 	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
// 		if token.Method != jwt.SigningMethodHS256 {
// 			return nil, errors.New("invalid token signing method")
// 		}

// 		return []byte(config.ServerConfigFunc().JwtSecret), nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	if !token.Valid {
// 		return nil, errors.New("invalid token")
// 	}

// 	if claims.TokenUse != "refresh" {
// 		return nil, errors.New("invalid refresh token")
// 	}

// 	return claims, nil
// }

// func GetUserIDFromContext(ctx context.Context) (int64, bool) {
// 	user_id, ok := ctx.Value(userIDContextKey).(int64)

// 	return user_id, ok
// }
