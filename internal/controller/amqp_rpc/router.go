package amqp_rpc

import (
	"github.com/deniskaponchik/GoSoft/internal/usecase"
	"github.com/deniskaponchik/GoSoft/pkg/rabbitmq/rmq_rpc/server"
)

// NewRouter -.
func NewRouter(uri usecase.UnifiRestIn) map[string]server.CallHandler {
	routes := make(map[string]server.CallHandler)
	{
		newUnifiRoutes(routes, uri)
	}

	return routes
}

//https://github.com/evrone/go-clean-template/blob/master/internal/controller/amqp_rpc/router.go
/*
// NewRouter -.
func NewRouter(t usecase.Translation) map[string]server.CallHandler {
	routes := make(map[string]server.CallHandler)
	{
		newTranslationRoutes(routes, t)
	}

	return routes
}*/
