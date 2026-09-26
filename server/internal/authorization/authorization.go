package authorization

var (
	AdminRole = "ADMIN"
	UserRole  = "USER"
)

const (
	PermissionCreateRole = "role:create"
	PermissionAddGenre = "genre:add"
	PermissionUpdateGenre  = "genre:update"
	PermissionAddRankings = "rankings:add"
	PermissionAddReview = "reviews:add"
	PermissionCreateMovie = "movie:create"
)

var RolePermissions = map[string][]string{
	AdminRole: {
		PermissionCreateRole,
		PermissionAddGenre,
		PermissionUpdateGenre,
		PermissionAddRankings,
		PermissionAddReview,
		PermissionCreateMovie,
	},
	UserRole: {
		PermissionAddReview,
	},
}
