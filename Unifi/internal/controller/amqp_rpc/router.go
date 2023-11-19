package amqp_rpc

import (
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase"
	"github.com/deniskaponchik/GoSoft/Unifi/pkg/rabbitmq/rmq_rpc/server"
)

// NewRouter -.
func NewRouter(t usecase.Translation) map[string]server.CallHandler {
	routes := make(map[string]server.CallHandler)
	{
		newTranslationRoutes(routes, t)
	}

	return routes
}
