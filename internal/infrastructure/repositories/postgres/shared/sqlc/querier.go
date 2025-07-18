// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	ArchiveConversation(ctx context.Context, id pgtype.UUID) error
	CleanupExpiredMagicLinks(ctx context.Context) error
	CountMessagesByConversationID(ctx context.Context, conversationID pgtype.UUID) (int64, error)
	CreateArtifact(ctx context.Context, arg CreateArtifactParams) (Artifact, error)
	CreateConversation(ctx context.Context, arg CreateConversationParams) (Conversation, error)
	CreateMagicLink(ctx context.Context, arg CreateMagicLinkParams) (MagicLink, error)
	CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error)
	CreateModel(ctx context.Context, arg CreateModelParams) (Model, error)
	CreateProvider(ctx context.Context, arg CreateProviderParams) (Provider, error)
	CreateTool(ctx context.Context, arg CreateToolParams) (Tool, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateUserProviderSetting(ctx context.Context, arg CreateUserProviderSettingParams) (UserProviderSetting, error)
	DeactivateUser(ctx context.Context, id pgtype.UUID) error
	DeleteArtifact(ctx context.Context, id pgtype.UUID) error
	DeleteConversation(ctx context.Context, id pgtype.UUID) error
	DeleteMessage(ctx context.Context, id pgtype.UUID) error
	DeleteModel(ctx context.Context, id pgtype.UUID) error
	DeleteProvider(ctx context.Context, id pgtype.UUID) error
	DeleteTool(ctx context.Context, id pgtype.UUID) error
	DeleteUserProviderSetting(ctx context.Context, id pgtype.UUID) error
	GetActiveModelsByProviderID(ctx context.Context, providerID pgtype.UUID) ([]Model, error)
	GetActiveProviders(ctx context.Context) ([]Provider, error)
	GetAllProviders(ctx context.Context) ([]Provider, error)
	GetArtifactByID(ctx context.Context, id pgtype.UUID) (Artifact, error)
	GetArtifactsByMessageID(ctx context.Context, messageID pgtype.UUID) ([]Artifact, error)
	GetAvailableModelsForUser(ctx context.Context, userID pgtype.UUID) ([]GetAvailableModelsForUserRow, error)
	GetAvailableTools(ctx context.Context, providerID pgtype.UUID) ([]Tool, error)
	GetConversationByID(ctx context.Context, id pgtype.UUID) (Conversation, error)
	GetConversationsByUserID(ctx context.Context, arg GetConversationsByUserIDParams) ([]Conversation, error)
	GetMagicLinkByToken(ctx context.Context, token string) (GetMagicLinkByTokenRow, error)
	GetMessageByID(ctx context.Context, id pgtype.UUID) (Message, error)
	GetMessageThread(ctx context.Context, parentID pgtype.UUID) ([]Message, error)
	GetMessagesByConversationID(ctx context.Context, arg GetMessagesByConversationIDParams) ([]Message, error)
	GetMessagesByConversationIDWithCursor(ctx context.Context, arg GetMessagesByConversationIDWithCursorParams) ([]Message, error)
	GetModelByID(ctx context.Context, id pgtype.UUID) (Model, error)
	GetModelByName(ctx context.Context, arg GetModelByNameParams) (Model, error)
	GetModelsByProviderID(ctx context.Context, providerID pgtype.UUID) ([]Model, error)
	GetProviderByID(ctx context.Context, id pgtype.UUID) (Provider, error)
	GetProviderByName(ctx context.Context, name string) (Provider, error)
	GetProvidersWithModels(ctx context.Context) ([]GetProvidersWithModelsRow, error)
	GetPublicArtifacts(ctx context.Context, arg GetPublicArtifactsParams) ([]Artifact, error)
	GetToolByID(ctx context.Context, id pgtype.UUID) (Tool, error)
	GetToolByName(ctx context.Context, name string) (Tool, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id pgtype.UUID) (User, error)
	GetUserProviderSetting(ctx context.Context, arg GetUserProviderSettingParams) (UserProviderSetting, error)
	InvalidateUserMagicLinks(ctx context.Context, arg InvalidateUserMagicLinksParams) error
	ListUserProviderSettings(ctx context.Context, userID pgtype.UUID) ([]UserProviderSetting, error)
	LogToolUsage(ctx context.Context, arg LogToolUsageParams) (MessageTool, error)
	UpdateArtifact(ctx context.Context, arg UpdateArtifactParams) (Artifact, error)
	UpdateConversation(ctx context.Context, arg UpdateConversationParams) (Conversation, error)
	UpdateConversationLastMessageAt(ctx context.Context, arg UpdateConversationLastMessageAtParams) error
	UpdateConversationTitle(ctx context.Context, arg UpdateConversationTitleParams) error
	UpdateMessage(ctx context.Context, arg UpdateMessageParams) (Message, error)
	UpdateModel(ctx context.Context, arg UpdateModelParams) (Model, error)
	UpdateProvider(ctx context.Context, arg UpdateProviderParams) (Provider, error)
	UpdateTool(ctx context.Context, arg UpdateToolParams) (Tool, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpdateUserProviderSetting(ctx context.Context, arg UpdateUserProviderSettingParams) (UserProviderSetting, error)
	UseMagicLink(ctx context.Context, id pgtype.UUID) (MagicLink, error)
	VerifyUserEmail(ctx context.Context, id pgtype.UUID) (User, error)
}

var _ Querier = (*Queries)(nil)
