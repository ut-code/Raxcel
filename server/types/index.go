package types

import (
	middleware "github.com/ut-code/Raxcel/server/middlewares"
	"github.com/ut-code/Raxcel/server/routes"
)

type AuthMiddlewareReturnType = middleware.AuthMiddlewareReturnType

// Auth requests and responses
type SignupRequest = routes.SignupRequest
type SignupResponse = routes.SignupResponse

type SigninRequest = routes.SigninRequest
type SigninResponse = routes.SigninResponse

// Message requests and responses
type ChatWithAIRequest = routes.ChatWithAIRequest

type ChatWithAIResponse struct {
	routes.ChatWithAIResponse
	*AuthMiddlewareReturnType
}

type LoadChatHistoryResponse struct {
	routes.LoadChatHistoryResponse
	*AuthMiddlewareReturnType
}

// User responses
type GetCurrentUserResponse struct {
	routes.GetCurrentUserResponse
	*AuthMiddlewareReturnType
}
