package types

import (
	middleware "github.com/ut-code/Raxcel/server/middlewares"
	"github.com/ut-code/Raxcel/server/routes"
)

type AuthMiddlewareReturn = middleware.AuthMiddlewareReturn

// Auth requests and responses
type SignupRequest = routes.SignupRequest
type SignupResponse = routes.SignupResponse

type SigninRequest = routes.SigninRequest
type SigninResponse = routes.SigninResponse

// Message requests and responses
type ChatWithAIRequest = routes.ChatWithAIRequest

type ChatWithAIResponse struct {
	routes.ChatWithAIResponse
	*AuthMiddlewareReturn
}

type LoadChatHistoryResponse struct {
	routes.LoadChatHistoryResponse
	*AuthMiddlewareReturn
}

// User responses
type GetCurrentUserResponse struct {
	routes.GetCurrentUserResponse
	*AuthMiddlewareReturn
}
