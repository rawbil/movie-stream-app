package middlewares

import (
	"errors"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	repository "github.com/rawbil/movie-stream-app/internal/adapters/sqlc"
	"github.com/rawbil/movie-stream-app/internal/auth/authutils"
	"github.com/rawbil/movie-stream-app/internal/utils"
)

func RoleMiddleware(repository repository.Queries, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		//~ Get user id from context
		user_id, err := authutils.GetUserIDFromContext(c)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "user_id not found in context", err)
			c.Abort()
			return
		}

		user_permissions, err := repository.ListUserPermissions(c.Request.Context(), user_id)
		if err != nil {
			utils.ErrorResponse(c, http.StatusForbidden, "error getting user permissions", err)
			c.Abort()
			return
		}

		if slices.Contains(user_permissions, permission) {
			c.Next()
			return
		}

		utils.ErrorResponse(c, http.StatusForbidden, "forbidden resource", errors.New("forbidden resource"))
		c.Abort()

	}
}
