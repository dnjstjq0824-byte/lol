package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"bn-riftstats/internal/domain"
)

type PostgresRepo struct {
	db *DB
}

func NewPostgresRepo(db *DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// GetUserByEmail finds a user profile and their password hash
func (r *PostgresRepo) GetUserByEmail(ctx context.Context, email string) (*domain.UserProfileDto, string, error) {
	if r.db == nil || r.db.pool == nil {
		return nil, "", fmt.Errorf("postgres pool unavailable")
	}

	query := `
		SELECT u.user_id, u.email, u.password_hash, u.role, p.display_name
		FROM users u
		LEFT JOIN user_profiles p ON u.user_id = p.user_id
		WHERE LOWER(u.email) = LOWER($1) AND u.is_active = TRUE
		LIMIT 1;
	`
	var (
		userID       uuid.UUID
		userEmail    string
		passwordHash string
		role         string
		displayName  *string
	)

	err := r.db.pool.QueryRow(ctx, query, email).Scan(&userID, &userEmail, &passwordHash, &role, &displayName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, "", fmt.Errorf("user not found")
		}
		return nil, "", err
	}

	name := userEmail
	if displayName != nil && *displayName != "" {
		name = *displayName
	}

	profile := &domain.UserProfileDto{
		UserID:      userID,
		Email:       userEmail,
		DisplayName: name,
		Role:        role,
	}
	return profile, passwordHash, nil
}

// UpsertUser inserts or updates a user
func (r *PostgresRepo) UpsertUser(ctx context.Context, email, passwordHash, displayName, role string) (*domain.UserProfileDto, error) {
	if r.db == nil || r.db.pool == nil {
		return nil, fmt.Errorf("postgres pool unavailable")
	}

	tx, err := r.db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var userID uuid.UUID
	userQuery := `
		INSERT INTO users (email, password_hash, role, is_active, updated_at)
		VALUES (LOWER($1), $2, $3, TRUE, NOW())
		ON CONFLICT (email) DO UPDATE
		SET password_hash = EXCLUDED.password_hash,
		    role = EXCLUDED.role,
		    updated_at = NOW()
		RETURNING user_id;
	`
	err = tx.QueryRow(ctx, userQuery, email, passwordHash, role).Scan(&userID)
	if err != nil {
		return nil, err
	}

	profileQuery := `
		INSERT INTO user_profiles (user_id, display_name, preferred_locale, preferred_region, updated_at)
		VALUES ($1, $2, 'ko_KR', 'KR', NOW())
		ON CONFLICT (user_id) DO UPDATE
		SET display_name = EXCLUDED.display_name,
		    updated_at = NOW();
	`
	_, err = tx.Exec(ctx, profileQuery, userID, displayName)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &domain.UserProfileDto{
		UserID:      userID,
		Email:       email,
		DisplayName: displayName,
		Role:        role,
	}, nil
}

// UpsertRiotAccount saves account data
func (r *PostgresRepo) UpsertRiotAccount(ctx context.Context, acc *domain.RiotAccount) error {
	if r.db == nil || r.db.pool == nil {
		return fmt.Errorf("postgres pool unavailable")
	}

	query := `
		INSERT INTO riot_accounts (puuid, game_name, tag_line, game_name_normalized, tag_line_normalized, primary_platform_region, updated_at)
		VALUES ($1, $2, $3, LOWER($2), LOWER($3), $4, NOW())
		ON CONFLICT (puuid) DO UPDATE
		SET game_name = EXCLUDED.game_name,
		    tag_line = EXCLUDED.tag_line,
		    game_name_normalized = EXCLUDED.game_name_normalized,
		    tag_line_normalized = EXCLUDED.tag_line_normalized,
		    primary_platform_region = EXCLUDED.primary_platform_region,
		    updated_at = NOW();
	`
	_, err := r.db.pool.Exec(ctx, query, acc.Puuid, acc.GameName, acc.TagLine, acc.Region)
	return err
}

// GetFavorites returns a user's bookmarked summoners or champions
func (r *PostgresRepo) GetFavorites(ctx context.Context, userID uuid.UUID) ([]domain.FavoriteItem, error) {
	if r.db == nil || r.db.pool == nil {
		return nil, fmt.Errorf("postgres pool unavailable")
	}

	query := `
		SELECT favorite_id, target_type, COALESCE(target_puuid, ''), COALESCE(target_champion_id, 0), platform_region, created_at
		FROM favorites
		WHERE user_id = $1
		ORDER BY created_at DESC;
	`
	rows, err := r.db.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.FavoriteItem
	for rows.Next() {
		var (
			favID      uuid.UUID
			targetType string
			puuid      string
			champID    int
			region     string
			createdAt  time.Time
		)
		if err := rows.Scan(&favID, &targetType, &puuid, &champID, &region, &createdAt); err != nil {
			return nil, err
		}
		list = append(list, domain.FavoriteItem{
			FavoriteID:       favID,
			TargetType:       targetType,
			PlatformRegion:   region,
			TargetPuuid:      puuid,
			TargetChampionID: champID,
			CreatedAt:        createdAt,
		})
	}
	return list, nil
}

// AddFavorite inserts a bookmark
func (r *PostgresRepo) AddFavorite(ctx context.Context, userID uuid.UUID, fav domain.FavoriteItem) error {
	if r.db == nil || r.db.pool == nil {
		return fmt.Errorf("postgres pool unavailable")
	}

	var puuidPtr *string
	if fav.TargetPuuid != "" {
		puuidPtr = &fav.TargetPuuid
	}
	var champIDPtr *int
	if fav.TargetChampionID > 0 {
		champIDPtr = &fav.TargetChampionID
	}

	query := `
		INSERT INTO favorites (user_id, target_type, target_puuid, target_champion_id, platform_region, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (user_id, target_type, platform_region, target_puuid, target_champion_id) DO NOTHING;
	`
	_, err := r.db.pool.Exec(ctx, query, userID, strings.ToUpper(fav.TargetType), puuidPtr, champIDPtr, fav.PlatformRegion)
	return err
}
