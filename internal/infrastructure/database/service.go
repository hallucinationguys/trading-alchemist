package database

import (
	"context"
	"fmt"

	"trading-alchemist/internal/domain/auth"
	"trading-alchemist/internal/domain/chat"
	authRepo "trading-alchemist/internal/infrastructure/repositories/postgres/auth"
	chatRepo "trading-alchemist/internal/infrastructure/repositories/postgres/chat"

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
	User() auth.UserRepository
	MagicLink() auth.MagicLinkRepository
	Provider() chat.ProviderRepository
	UserProviderSetting() chat.UserProviderSettingRepository
	Conversation() chat.ConversationRepository
	Message() chat.MessageRepository
	Artifact() chat.ArtifactRepository
	Tool() chat.ToolRepository
	Model() chat.ModelRepository
}

// transactionalRepositoryProvider provides repositories that are bound to a specific database transaction.
type transactionalRepositoryProvider struct {
	tx DBTX
}

// NewTransactionalRepositoryProvider creates a new provider for transactional repositories.
func NewTransactionalRepositoryProvider(tx DBTX) RepositoryProvider {
	return &transactionalRepositoryProvider{tx: tx}
}

func (p *transactionalRepositoryProvider) User() auth.UserRepository {
	return authRepo.NewUserRepository(p.tx)
}

func (p *transactionalRepositoryProvider) MagicLink() auth.MagicLinkRepository {
	return authRepo.NewMagicLinkRepository(p.tx)
}

func (p *transactionalRepositoryProvider) Conversation() chat.ConversationRepository {
	return chatRepo.NewConversationRepository(p.tx)
}

func (p *transactionalRepositoryProvider) Message() chat.MessageRepository {
	return chatRepo.NewMessageRepository(p.tx)
}

func (p *transactionalRepositoryProvider) Artifact() chat.ArtifactRepository {
	return chatRepo.NewArtifactRepository(p.tx)
}

func (p *transactionalRepositoryProvider) Provider() chat.ProviderRepository {
	return chatRepo.NewProviderRepository(p.tx)
}

func (p *transactionalRepositoryProvider) Tool() chat.ToolRepository {
	return chatRepo.NewToolRepository(p.tx)
}

func (p *transactionalRepositoryProvider) UserProviderSetting() chat.UserProviderSettingRepository {
	return chatRepo.NewUserProviderSettingRepository(p.tx)
}

func (trp *transactionalRepositoryProvider) Model() chat.ModelRepository {
	return chatRepo.NewModelRepository(trp.tx)
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