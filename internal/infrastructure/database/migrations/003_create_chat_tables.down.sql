DROP TRIGGER IF EXISTS update_artifacts_updated_at ON artifacts;
DROP TRIGGER IF EXISTS update_messages_updated_at ON messages;
DROP TRIGGER IF EXISTS update_conversations_updated_at ON conversations;
DROP TRIGGER IF EXISTS update_tools_updated_at ON tools;
DROP TRIGGER IF EXISTS update_models_updated_at ON models;
DROP TRIGGER IF EXISTS update_user_provider_settings_updated_at ON user_provider_settings;
DROP TRIGGER IF EXISTS update_providers_updated_at ON providers;

DROP TABLE IF EXISTS message_tools;
DROP TABLE IF EXISTS artifacts;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS conversations;
DROP TABLE IF EXISTS tools;
DROP TABLE IF EXISTS models;
DROP TABLE IF EXISTS user_provider_settings;
DROP TABLE IF EXISTS providers; 