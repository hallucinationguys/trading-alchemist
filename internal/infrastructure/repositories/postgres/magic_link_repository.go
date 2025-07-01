package postgres

import (
	"context"
	"fmt"
	"net"
	"net/netip"

	"trading-alchemist/internal/domain/entities"
	"trading-alchemist/internal/domain/repositories"
	"trading-alchemist/internal/infrastructure/repositories/postgres/sqlc"
	"trading-alchemist/pkg/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// MagicLinkRepository implements the domain's MagicLinkRepository interface using PostgreSQL.
type MagicLinkRepository struct {
	queries *sqlc.Queries
}

// NewMagicLinkRepository creates a new postgres magic link repository.
func NewMagicLinkRepository(db sqlc.DBTX) repositories.MagicLinkRepository {
	return &MagicLinkRepository{
		queries: sqlc.New(db),
	}
}

// Create creates a new magic link.
func (r *MagicLinkRepository) Create(ctx context.Context, magicLink *entities.MagicLink) (*entities.MagicLink, error) {
	params := sqlc.CreateMagicLinkParams{
		UserID:    pgtype.UUID{Bytes: magicLink.UserID, Valid: true},
		Token:     magicLink.Token,
		TokenHash: magicLink.TokenHash,
		ExpiresAt: pgtype.Timestamptz{Time: magicLink.ExpiresAt, Valid: true},
		Purpose:   string(magicLink.Purpose),
	}

	// Handle optional IP address
	if magicLink.IPAddress != nil {
		if ip := net.ParseIP(*magicLink.IPAddress); ip != nil {
			if netipAddr, ok := netip.AddrFromSlice(ip); ok {
				params.IpAddress = &netipAddr
			}
		}
	}

	// Handle optional user agent
	if magicLink.UserAgent != nil {
		params.UserAgent = pgtype.Text{String: *magicLink.UserAgent, Valid: true}
	}

	sqlcMagicLink, err := r.queries.CreateMagicLink(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create magic link: %w", err)
	}

	return r.sqlcMagicLinkToEntity(&sqlcMagicLink), nil
}

// GetByToken retrieves a magic link by its token (with user info).
func (r *MagicLinkRepository) GetByToken(ctx context.Context, token string) (*entities.MagicLink, *entities.User, error) {
	row, err := r.queries.GetMagicLinkByToken(ctx, token)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil, errors.ErrMagicLinkNotFound
		}
		return nil, nil, fmt.Errorf("failed to get magic link by token: %w", err)
	}

	// Convert magic link
	magicLink := &entities.MagicLink{
		Token:     row.Token,
		TokenHash: row.TokenHash,
		Purpose:   entities.MagicLinkPurpose(row.Purpose),
		CreatedAt: row.CreatedAt.Time,
	}

	if row.ID.Valid {
		magicLink.ID = row.ID.Bytes
	}
	if row.UserID.Valid {
		magicLink.UserID = row.UserID.Bytes
	}
	if row.ExpiresAt.Valid {
		magicLink.ExpiresAt = row.ExpiresAt.Time
	}
	if row.UsedAt.Valid {
		magicLink.UsedAt = &row.UsedAt.Time
	}
	if row.IpAddress != nil {
		ipStr := row.IpAddress.String()
		magicLink.IPAddress = &ipStr
	}
	if row.UserAgent.Valid {
		magicLink.UserAgent = &row.UserAgent.String
	}

	// Convert user
	user := &entities.User{
		ID:            row.UserID.Bytes,
		Email:         row.Email,
		EmailVerified: row.EmailVerified.Bool,
	}

	return magicLink, user, nil
}

// MarkAsUsed marks a magic link as used.
func (r *MagicLinkRepository) MarkAsUsed(ctx context.Context, linkID uuid.UUID) (*entities.MagicLink, error) {
	linkUUID := pgtype.UUID{Bytes: linkID, Valid: true}
	
	sqlcMagicLink, err := r.queries.UseMagicLink(ctx, linkUUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrMagicLinkNotFound
		}
		return nil, fmt.Errorf("failed to mark magic link as used: %w", err)
	}

	return r.sqlcMagicLinkToEntity(&sqlcMagicLink), nil
}

// InvalidateUserLinks invalidates all unused magic links for a user with a specific purpose.
func (r *MagicLinkRepository) InvalidateUserLinks(ctx context.Context, userID uuid.UUID, purpose entities.MagicLinkPurpose) error {
	params := sqlc.InvalidateUserMagicLinksParams{
		UserID:  pgtype.UUID{Bytes: userID, Valid: true},
		Purpose: string(purpose),
	}

	err := r.queries.InvalidateUserMagicLinks(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to invalidate user magic links: %w", err)
	}

	return nil
}

// CleanupExpired removes expired magic links.
func (r *MagicLinkRepository) CleanupExpired(ctx context.Context) error {
	err := r.queries.CleanupExpiredMagicLinks(ctx)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired magic links: %w", err)
	}

	return nil
}

// sqlcMagicLinkToEntity converts a SQLC MagicLink to a domain MagicLink entity.
func (r *MagicLinkRepository) sqlcMagicLinkToEntity(sqlcMagicLink *sqlc.MagicLink) *entities.MagicLink {
	magicLink := &entities.MagicLink{
		Token:     sqlcMagicLink.Token,
		TokenHash: sqlcMagicLink.TokenHash,
		Purpose:   entities.MagicLinkPurpose(sqlcMagicLink.Purpose),
		CreatedAt: sqlcMagicLink.CreatedAt.Time,
	}

	// Convert UUID
	if sqlcMagicLink.ID.Valid {
		magicLink.ID = sqlcMagicLink.ID.Bytes
	}
	if sqlcMagicLink.UserID.Valid {
		magicLink.UserID = sqlcMagicLink.UserID.Bytes
	}

	// Convert timestamps
	if sqlcMagicLink.ExpiresAt.Valid {
		magicLink.ExpiresAt = sqlcMagicLink.ExpiresAt.Time
	}
	if sqlcMagicLink.UsedAt.Valid {
		magicLink.UsedAt = &sqlcMagicLink.UsedAt.Time
	}

	// Convert optional fields
	if sqlcMagicLink.IpAddress != nil {
		ipStr := sqlcMagicLink.IpAddress.String()
		magicLink.IPAddress = &ipStr
	}
	if sqlcMagicLink.UserAgent.Valid {
		magicLink.UserAgent = &sqlcMagicLink.UserAgent.String
	}

	return magicLink
} 