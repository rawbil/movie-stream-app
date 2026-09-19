package utils

type CreateMovieParams struct {
	ImdbID       string `json:"imdb_id" validate:"required"`
	Title        string `json:"title" validate:"required,min=3,max=500"`
	PosterPath   string `json:"poster_path" validate:"required,url"`
	YoutubeID    string `json:"youtube_id"`
	AdminReview  string `json:"admin_review"`
	RankingValue int    `json:"ranking_value"`
	RankingName  string `json:"ranking_name"`
}

type CreateUserParams struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password_format"`
}

type CreateGenreParams struct {
	GenreName string `json:"current_genre" validate:"required"`
}

type UpdateGenreParams struct {
	CurrentGenre string `json:"current_genre" validate:"required"`
	OldGenre     string `json:"old_genre" validate:"required"`
}
