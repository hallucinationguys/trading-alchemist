package database

import (
	"context"
	"fmt"

	"trading-alchemist/internal/domain/repositories"
	"trading-alchemist/internal/infrastructure/repositories/postgres"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBTX is an interface that both *pgx.Tx and *pgxpool.Pool satisfy.
// This allows us to use either a transaction or a connection pool for our repositories.
type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

// RepositoryProvider defines the interface for accessing all repositories.
// This allows use cases to depend on this interface for transactional database operations.
type RepositoryProvider interface {
	User() repositories.UserRepository
	MagicLink() repositories.MagicLinkRepository
	Provider() repositories.ProviderRepository
	UserProviderSetting() repositories.UserProviderSettingRepository
	Conversation() repositories.ConversationRepository
	Message() repositories.MessageRepository
	Artifact() repositories.ArtifactRepository
	Tool() repositories.ToolRepository
	Model() repositories.ModelRepository
}

// transactionalRepositoryProvider provides repositories that are bound to a specific database transaction.
type transactionalRepositoryProvider struct {
	tx DBTX
}

// NewTransactionalRepositoryProvider creates a new provider for transactional repositories.
func NewTransactionalRepositoryProvider(tx DBTX) RepositoryProvider {
	return &transactionalRepositoryProvider{tx: tx}
}

func (p *transactionalRepositoryProvider) User() repositories.UserRepository {
	return postgres.NewUserRepository(p.tx)
}

func (p *transactionalRepositoryProvider) MagicLink() repositories.MagicLinkRepository {
	return postgres.NewMagicLinkRepository(p.tx)
}

func (p *transactionalRepositoryProvider) Conversation() repositories.ConversationRepository {
	return postgres.NewConversationRepository(p.tx)
}

func (p *transactionalRepositoryProvider) Message() repositories.MessageRepository {
	return postgres.NewMessageRepository(p.tx)
}

func (p *transactionalRepositoryProvider) Artifact() repositories.ArtifactRepository {
	return postgres.NewArtifactRepository(p.tx)
}

func (p *transactionalRepositoryProvider) Provider() repositories.ProviderRepository {
	return postgres.NewProviderRepository(p.tx)
}

func (p *transactionalRepositoryProvider) Tool() repositories.ToolRepository {
	return postgres.NewToolRepository(p.tx)
}

func (p *transactionalRepositoryProvider) UserProviderSetting() repositories.UserProviderSettingRepository {
	return postgres.NewUserProviderSettingRepository(p.tx)
}

func (trp *transactionalRepositoryProvider) Model() repositories.ModelRepository {
	return postgres.NewModelRepository(trp.tx)
}

// Service provides a high-level abstraction for database operations,
// including transaction management.
type Service struct {
	pool *pgxpool.Pool
}

// NewService creates a new database service.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// ExecuteInTx runs the provided function within a single database transaction.
// If the function returns an error, the transaction is automatically rolled back.
// Otherwise, the transaction is committed.
func (s *Service) ExecuteInTx(ctx context.Context, fn func(provider RepositoryProvider) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // Rollback is a no-op if the transaction has been committed.

	// Create a repository provider that uses the transaction.
	provider := NewTransactionalRepositoryProvider(tx)

	// Execute the core logic.
	if err := fn(provider); err != nil {
		return err // The defer will handle the rollback.
	}

	// Commit the transaction.
	return tx.Commit(ctx)
} 