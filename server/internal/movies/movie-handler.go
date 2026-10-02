package movies

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	repository "github.com/rawbil/movie-stream-app/internal/adapters/sqlc"
	"github.com/rawbil/movie-stream-app/internal/auth/authutils"
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

// ! List Movies
func (h *Handler) ListMovies(c *gin.Context) {
	// limitString := c.DefaultQuery("limit", "10")
	pageString := c.DefaultQuery("page", "1")

	// limit, err := strconv.Atoi(limitString)
	// if err != nil {
	// 	utils.ErrorResponse(c, http.StatusBadRequest, "limit must be a number", err)
	// 	return
	// }

	page, err := strconv.Atoi(pageString)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "page must be a number", err)
		return
	}

	// if limit < 1 {
	// 	utils.ErrorResponse(c, http.StatusBadRequest, "limit must be greater than zero", errors.New("invalid limit"))
	// 	return
	// }
	if page < 1 {
		utils.ErrorResponse(c, http.StatusBadRequest, "page must be greater than zero", errors.New("invalid page"))
		return
	}
	limit := 10

	offset := (page - 1) * limit

	movies, total_movies, err := h.Service.ListMovies(c.Request.Context(), repository.ListMoviesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	total_pages := math.Ceil(float64(total_movies) / float64(limit))

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"movies": movies,
			"page":   page,
			// "limit":        limit,
			"skip":        offset,
			"total_pages": total_pages,
		},
	})
}

// ! Get Movie
func (h *Handler) GetMovie(c *gin.Context) {
	public_id_from_url := c.Query("public_id")
	publicID, err := uuid.Parse(public_id_from_url)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid public_id", err)
		return
	}

	movie, err := h.Service.GetMovie(c.Request.Context(), publicID)
	if err != nil {
		if err == utils.NoRecordError {
			utils.ErrorResponse(c, http.StatusNotFound, "movie not found", err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"movie":   movie,
	})
}

// ! Create Movie
func (h *Handler) CreateMovie(c *gin.Context) {
	var params utils.CreateMovieParams

	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "error decoding json body", err)
		return
	}

	if err := h.Service.CreateMovie(c.Request.Context(), params); err != nil {
		if err == utils.AllFieldsRequiredError || err == utils.MinTitleError || err == utils.MaxTitleError || err == utils.InvalidUrlError || err == utils.GenreMissing {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
			return
		}

		if err == utils.MovieExistsError {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), err)
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Movie '%s' created", params.Title),
	})
}

// ! Add Rankings
func (h *Handler) AddRankings(c *gin.Context) {
	var params utils.AddRankingsParams

	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "error decoding json body", err)
		return
	}

	if err, ranking_name := h.Service.AddRankings(c.Request.Context(), params); err != nil {
		if err == utils.AllFieldsRequiredError {
			utils.ErrorResponse(c, http.StatusBadRequest, "missing value", err)
			return
		}

		if err == utils.GenreMissing {
			utils.ErrorResponse(c, http.StatusBadRequest, "at least one ranking needed", err)
			return
		}

		if err == utils.MovieExistsError {
			utils.ErrorResponse(c, http.StatusConflict, fmt.Sprintf("ranking '%s' already exists", ranking_name), err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "rankings added",
	})
}

// ! AddMovieReview
func (h *Handler) AddMovieReview(c *gin.Context) {
	var params utils.AddReviewParams

	public_id_from_url := c.Params.ByName("public_id")
	if public_id_from_url == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "public_id not found in params", errors.New("public_id not found in params"))
		return
	}

	publicID, err := uuid.Parse(public_id_from_url)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid public_id", err)
		return
	}

	user_id, err := authutils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), err)
		return
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "error decoding json body", err)
		return
	}

	if err := h.Service.AddReview(c.Request.Context(), publicID, user_id, params); err != nil {
		if err == utils.AllFieldsRequiredError {
			utils.ErrorResponse(c, http.StatusBadRequest, "review field required", err)
			return
		}

		if err == utils.NoRecordError {
			utils.ErrorResponse(c, http.StatusNotFound, "movie not found", err)
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

// ! MovieRecommendations
func (h *Handler) GetMovieRecommendations(c *gin.Context) {
	user_id, err := authutils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), err)
		return
	}

	movies, err := h.Service.GetMovieRecommendations(c.Request.Context(), user_id)
	if err != nil {
		if err == utils.NoRecordError {
			utils.ErrorResponse(c, http.StatusNotFound, "your favourite genres have no movies yet", err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movies": movies,
	})

}

// ! Get Movie Reviews
func (h *Handler) GetMovieReviews(c *gin.Context) {
	public_id_from_url := c.Params.ByName("public_id")
	if public_id_from_url == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "public_id not found in params", errors.New("public_id not found in params"))
		return
	}

	publicID, err := uuid.Parse(public_id_from_url)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid public_id", err)
		return
	}

	pageString := c.DefaultQuery("page", "1")

	page, err := strconv.Atoi(pageString)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "page must be a number", err)
		return
	}

	if page < 1 {
		utils.ErrorResponse(c, http.StatusBadRequest, "page must be greater than zero", errors.New("invalid page"))
		return
	}

	limit := 10

	offset := (page - 1) * limit

	reviews, totalReviews, err := h.Service.GetMovieReviews(c.Request.Context(), repository.GetMovieReviewsParams{
		PublicID: publicID[:],
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	totalPages := math.Ceil(float64(totalReviews) / float64(limit))

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"reviews":       reviews,
			"page":          page,
			"total_pages":   totalPages,
			"total_reviews": totalReviews,
		},
	})
}
