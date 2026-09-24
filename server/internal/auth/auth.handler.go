package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rawbil/movie-stream-app/internal/utils"
)

type Handler struct {
	Service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) RegisterUser(c *gin.Context) {
	var params utils.CreateUserParams

	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "error decoding json body", err)
		return
	}

	if err := h.Service.RegisterUser(c.Request.Context(), params); err != nil {
		if err == utils.AllFieldsRequiredError || err == utils.InvalidEmailFormat || err == utils.InvalidPasswordFormat {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
			return
		}

		if err == utils.MinTitleError {
			utils.ErrorResponse(c, http.StatusBadRequest, "username should have at least 3 characters", err)
			return
		}

		if err == utils.MaxTitleError {
			utils.ErrorResponse(c, http.StatusBadRequest, "username should not have more than 20 characters", err)
			return
		}

		if err == utils.GenreMissing {
			utils.ErrorResponse(c, http.StatusBadRequest, "choose at least one genre", err)
			return
		}

		if err == utils.NoEmptyGenre {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
			return
		}

		if err == utils.DuplicateRecordError {
			utils.ErrorResponse(c, http.StatusConflict, "user with email already exists", err)
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful!",
	})
}

func (h *Handler) LoginUser(c *gin.Context) {
	var params utils.LoginParams

	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "error decoding json body", err)
		return
	}

	user, access_token, refresh_token, err := h.Service.LoginUser(c.Request.Context(), params)
	if err != nil {

		if err == utils.AllFieldsRequiredError || err == utils.InvalidEmailFormat {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
			return
		}

		if err == utils.NoRecordError || err == utils.IncorrectPassword {
			utils.ErrorResponse(c, http.StatusBadRequest, "invalid credentials... try again", err)
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	// rt_cookie := http.Cookie{
	// 	Name:     "__refresh_token__",
	// 	Value:    refresh_token,
	// 	MaxAge:   7 * 24 * 60 * 60, // 7 days in seconds
	// 	Path:     "/",
	// 	HttpOnly: true,                                     // Prevents client-side JS access
	// 	Secure:   utils.ServerConfigFunc().AppEnv != "dev", // true in prod (SET APP_ENV=prod)
	// }

	// http.SetCookie(w, &rt_cookie)
	maxAge := 7 * 24 * 60 * 60

	// Determine if we are in production
	isProd := utils.ServerConfigFunc().AppEnv != "dev"
	cookieName := "__refresh_token"
	if isProd {
		cookieName = "__Host-refresh_token" //__HOST makes the cookie secure and ensures it is only transmitted via https
	}

	c.SetCookie(
		cookieName,    // Cookie Name
		refresh_token, // Value
		maxAge,        // MaxAge in seconds (Fixed)
		"/",           // Path
		"",            // Domain (Leaving this empty is usually best)
		isProd,        // Secure (true means HTTPS only)
		true,          // HttpOnly (Prevents XSS access)
	)

	updated_user := map[string]any{
		"public_id":  user.PublicID,
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Login successful",
		"access_token": access_token,
		"user":         updated_user,
	})
}

func (h *Handler) CreateRole(c *gin.Context) {
	var params utils.CreateRoleParams

	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "error decoding json body", err)
		return
	}

	if _, err := h.Service.CreateRole(c.Request.Context(), params.Role); err != nil {
		if err == utils.AllFieldsRequiredError {
			utils.ErrorResponse(c, http.StatusBadRequest, "role field required", err)
			return
		}

		if err == utils.DuplicateRecordError {
			utils.ErrorResponse(c, http.StatusConflict, "role already exists", err)
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "role added",
	})
}

func (h *Handler) Logout(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "user not found in context", errors.New("user missing in context"))
		return
	}

	user_id, ok := userIDValue.(int64)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid user id in context", errors.New("invalid user id in context"))
		return
	}

	if err := h.Service.Logout(c.Request.Context(), user_id); err != nil {
		if err == utils.NoRecordError {
			utils.ErrorResponse(c, http.StatusUnauthorized, "user not found", err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

func (h *Handler) RefreshTokens(c *gin.Context) {
	isProd := utils.ServerConfigFunc().AppEnv != "dev"

	cookieName := "__refresh_token"
	if isProd {
		cookieName = "__Host-refresh_token"	
	}

	rt_cookie, err := c.Request.Cookie(cookieName)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "refresh token missing", err)
		return
	}

	rt := rt_cookie.Value

	access_token, refresh_token, err := h.Service.RefreshTokens(c.Request.Context(), rt)
	if err != nil {
		if err == utils.AllFieldsRequiredError {
			utils.ErrorResponse(c, http.StatusBadRequest, "refresh token missing", err)
			return
		}

		if err == utils.NoRecordError {
			utils.ErrorResponse(c, http.StatusNotFound, "token missing in records", err)
			return
		}

		if err == utils.TokenRevoked || err == utils.InvalidToken {
			utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	maxAge := 7 * 24 * 60 * 60

	c.SetCookie(
		"__Host-refresh_token", // Cookie Name (Consider changing to "__Host-refresh_token" for maximum security)
		refresh_token,          // Value
		maxAge,                 // MaxAge in seconds (Fixed)
		"/",                    // Path
		"",                     // Domain (Leaving this empty is usually best)
		isProd,                 // Secure (true means HTTPS only)
		true,                   // HttpOnly (Prevents XSS access)
	)

	c.JSON(http.StatusOK, gin.H{
		"message":      "success",
		"access_token": access_token,
	})
}
