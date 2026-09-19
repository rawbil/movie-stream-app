package movies

import (
	"fmt"
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

// ! Create Genre
func (h *Handler) CreateGenre(c *gin.Context) {
	var params utils.CreateGenreParams
	// if err := utils.DecodeJson(r, &params); err != nil {
	// 	utils.ErrorResponse(w, "error reading json body", err, http.StatusBadRequest)
	// 	return
	// }
	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "error reading json body", err)
		return
	}

	genreName, err := h.Service.CreateGenre(c.Request.Context(), params.GenreName)

	if err != nil {
		if err == utils.AllFieldsRequiredError {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
			return
		}
		if err == utils.NoRecordError {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), err)
			return
		}
		if err == utils.DuplicateRecordError {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	// utils.SuccessResponse(c, utils.Response{
	// 	Message: fmt.Sprintf("Genre '%s' added", genreName),
	// })

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Genre '%s' added", genreName),
	})
}

// ! UpdateGenre
func (h *Handler) UpdateGenre(c *gin.Context) {
	var params utils.UpdateGenreParams

	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "error reading json body", err)
		return
	}

	new_genre, err := h.Service.UpdateGenre(c.Request.Context(), params)
	if err != nil {
		if err == utils.AllFieldsRequiredError {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
			return
		}

		if err == utils.NoRecordError {
			utils.ErrorResponse(c, http.StatusBadRequest, "genre not found in records", err)
			return
		}

		if err == utils.DuplicateRecordError {
			utils.ErrorResponse(c, http.StatusBadRequest, "genre already exists", err)
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Genre updated to '%s'", new_genre),
	})
}

// ! List Genres
func (h *Handler) ListGenres(c *gin.Context) {
	genres, err := h.Service.ListGenres(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"genres": genres,
		},
	})
}
