package authorization

var (
	AdminRole = "ADMIN"
	UserRole  = "USER"
)

const (
	PermissionCreateRole = "role:create"
	PermissionAddGenre    = "genre:add"
	PermissionUpdateGenre  = "genre:update"
	PermissionAddRankings  = "rankings:add"
	PermissionAddReview    = "reviews:add"
	PermissionGetReviews = "reviews:get"
	PermissionCreateMovie = "movie:create"
	PermissionUpdateMovie = "movie:update"
	PermissionDeleteMovie = "movie:delete"
)

var RolePermissions = map[string][]string{
	AdminRole: {
		PermissionCreateRole,
		PermissionAddGenre,
		PermissionUpdateGenre,
		PermissionAddRankings,
		PermissionAddReview,
		PermissionGetReviews,
		PermissionCreateMovie,
		PermissionUpdateMovie,
		PermissionDeleteMovie,
	},
	UserRole: {
		PermissionAddReview,
		PermissionGetReviews,
	},
}
