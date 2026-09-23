package auth

import (
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
