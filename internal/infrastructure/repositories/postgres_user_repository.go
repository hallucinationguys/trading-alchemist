package repositories

import (
	"context"
	"fmt"

	"trading-alchemist/internal/domain/entities"
	"trading-alchemist/internal/domain/repositories"
	"trading-alchemist/internal/infrastructure/repositories/sqlc"
	"trading-alchemist/pkg/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// PostgresUserRepository implements UserRepository interface using PostgreSQL
type PostgresUserRepository struct {
	queries *sqlc.Queries
}

// NewPostgresUserRepository creates a new postgres user repository
func NewPostgresUserRepository(db sqlc.DBTX) repositories.UserRepository {
	return &PostgresUserRepository{
		queries: sqlc.New(db),
	}
}

// Create creates a new user
func (r *PostgresUserRepository) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	params := sqlc.CreateUserParams{
		Email:         user.Email,
		EmailVerified: pgtype.Bool{Bool: user.EmailVerified, Valid: true},
	}

	if user.FirstName != nil {
		params.FirstName = pgtype.Text{String: *user.FirstName, Valid: true}
	}
	if user.LastName != nil {
		params.LastName = pgtype.Text{String: *user.LastName, Valid: true}
	}

	sqlcUser, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return r.sqlcUserToEntity(&sqlcUser), nil
}

// GetByID retrieves a user by their ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	userUUID := pgtype.UUID{Bytes: id, Valid: true}
	
	sqlcUser, err := r.queries.GetUserByID(ctx, userUUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return r.sqlcUserToEntity(&sqlcUser), nil
}

// GetByEmail retrieves a user by their email
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	sqlcUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return r.sqlcUserToEntity(&sqlcUser), nil
}

// Update updates user information
func (r *PostgresUserRepository) Update(ctx context.Context, user *entities.User) (*entities.User, error) {
	userUUID := pgtype.UUID{Bytes: user.ID, Valid: true}
	
	params := sqlc.UpdateUserParams{
		ID: userUUID,
	}

	if user.FirstName != nil {
		params.FirstName = pgtype.Text{String: *user.FirstName, Valid: true}
	}
	if user.LastName != nil {
		params.LastName = pgtype.Text{String: *user.LastName, Valid: true}
	}

	sqlcUser, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return r.sqlcUserToEntity(&sqlcUser), nil
}

// VerifyEmail marks a user's email as verified
func (r *PostgresUserRepository) VerifyEmail(ctx context.Context, userID uuid.UUID) (*entities.User, error) {
	userUUID := pgtype.UUID{Bytes: userID, Valid: true}
	
	sqlcUser, err := r.queries.VerifyUserEmail(ctx, userUUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to verify user email: %w", err)
	}

	return r.sqlcUserToEntity(&sqlcUser), nil
}

// Deactivate deactivates a user account
func (r *PostgresUserRepository) Deactivate(ctx context.Context, userID uuid.UUID) error {
	userUUID := pgtype.UUID{Bytes: userID, Valid: true}
	
	err := r.queries.DeactivateUser(ctx, userUUID)
	if err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	return nil
}

// sqlcUserToEntity converts SQLC User to domain entity
func (r *PostgresUserRepository) sqlcUserToEntity(sqlcUser *sqlc.User) *entities.User {
	user := &entities.User{
		Email:     sqlcUser.Email,
		IsActive:  sqlcUser.IsActive.Bool,
		CreatedAt: sqlcUser.CreatedAt.Time,
		UpdatedAt: sqlcUser.UpdatedAt.Time,
	}

	// Convert UUID
	if sqlcUser.ID.Valid {
		user.ID = sqlcUser.ID.Bytes
	}

	// Convert optional fields
	if sqlcUser.EmailVerified.Valid {
		user.EmailVerified = sqlcUser.EmailVerified.Bool
	}

	if sqlcUser.FirstName.Valid {
		user.FirstName = &sqlcUser.FirstName.String
	}

	if sqlcUser.LastName.Valid {
		user.LastName = &sqlcUser.LastName.String
	}

	if sqlcUser.AvatarUrl.Valid {
		user.AvatarURL = &sqlcUser.AvatarUrl.String
	}

	return user
} 