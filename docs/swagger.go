// Package docs contains the Swagger documentation for Trading Alchemist API
// @title Trading Alchemist API
// @version 1.0.0
// @description A modern Go-based authentication system with magic link functionality using Clean Architecture principles.
// @description This API provides secure authentication mechanisms including magic link authentication, JWT token management, and user profile management.
// @termsOfService http://swagger.io/terms/
// @contact.name Trading Alchemist Team
// @contact.email team@tradingalchemist.dev
// @license.name MIT
// @license.url http://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT token. Usage: "Bearer {token}"
package docs

import (
	_ "trading-alchemist/internal/application/auth"
	_ "trading-alchemist/internal/application/chat"
	_ "trading-alchemist/internal/presentation/http/handlers"
	_ "trading-alchemist/internal/presentation/responses"
) 