package seed

import (
	"context"
	"database/sql"

	repository "github.com/rawbil/movie-stream-app/internal/adapters/sqlc"
	"github.com/rawbil/movie-stream-app/internal/authorization"
)

func SeedPermissions(db *sql.DB) error {
	permissions := []string{
		authorization.PermissionCreateRole,
		authorization.PermissionAddGenre,
		authorization.PermissionUpdateGenre,
		authorization.PermissionAddRankings,
		authorization.PermissionAddReview,
		authorization.PermissionCreateMovie,
	}

	for _, permission := range permissions {
		query := `
		INSERT IGNORE INTO user_permissions(permission)
		VALUES (?)
		`

		_, err := db.Exec(query, permission)
		if err != nil {
			return err
		}
	}

	return nil
}

func SeedRoles(db *sql.DB) error {
	roles := []string{
		authorization.AdminRole,
		authorization.UserRole,
	}

	for _, role := range roles {
		query := `
		INSERT IGNORE INTO roles(role)
		VALUES (?)
		`

		if _, err := db.Exec(query, role); err != nil {
			return err
		}
	}

	return nil
}

func SeedRolePermissions(ctx context.Context, repository repository.Queries, db *sql.DB) error {
	for roleName, permissions := range authorization.RolePermissions {
		role_id, err := repository.GetRoleID(ctx, roleName)
		if err != nil {
			return err
		}

		for _, permission := range permissions {
			permission_id, err := repository.GetPermissionID(ctx, permission)
			if err != nil {
				return err
			}

			query := `
			INSERT IGNORE INTO role_permissions(role_id, permission_id)
			VALUES (?, ?)
			`

			if _, err := db.Exec(query, role_id, permission_id); err != nil {
				return err
			}
		}
	}

	return nil
}

func SeedMovieRankings(db *sql.DB) error {
	rankings := map[int]string{
		1: "Excellent",
		2: "Good",
		3: "Okay",
		4: "Bad",
		5: "Terrible",
	}

	for index, value := range rankings {
		query := `
		INSERT IGNORE INTO rankings(ranking_name, ranking_value)
        VALUES (?, ?);
		`

		if _, err := db.Exec(query, value, index); err != nil {
			return err
		}
	}

	return nil
}
