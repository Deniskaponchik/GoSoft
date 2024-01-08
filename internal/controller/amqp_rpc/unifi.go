package amqp_rpc

import (
	"github.com/deniskaponchik/GoSoft/internal/entity"
	"github.com/deniskaponchik/GoSoft/internal/usecase"
	"github.com/deniskaponchik/GoSoft/pkg/rabbitmq/rmq_rpc/server"
	"github.com/streadway/amqp"
)

type unifiRoutes struct {
	unifiUseCase usecase.UnifiRestIn //Interface
}

func newUnifiRoutes(routes map[string]server.CallHandler, uri usecase.UnifiRestIn) {
	r := &unifiRoutes{uri}
	{
		routes["getClient"] = r.getClient()
	}
}

type clientResponse struct {
	Client []entity.Client `json:"client"`
}

type apResponse struct {
	Ap []entity.Ap `json:"ap"`
}

func (r *unifiRoutes) getClient() server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		/*TODO: переписать под моё
		client, err := r.unifiUseCase.GetClientForRest(context.Background().Value())
		if err != nil {
			return nil, fmt.Errorf("amqp_rpc - translationRoutes - getHistory - r.translationUseCase.History: %w", err)
		}

		response := clientResponse{client}

		return response, nil*/
		return nil, nil
	}
}

/*
type translationRoutes struct {
	translationUseCase usecase.Translation
}

func newTranslationRoutes(routes map[string]server.CallHandler, t usecase.Translation) {
	r := &translationRoutes{t}
	{
		routes["getHistory"] = r.getHistory()
	}
}

type historyResponse struct {
	History []entity.Translation `json:"history"`
}

func (r *translationRoutes) getHistory() server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		translations, err := r.translationUseCase.History(context.Background())
		if err != nil {
			return nil, fmt.Errorf("amqp_rpc - translationRoutes - getHistory - r.translationUseCase.History: %w", err)
		}

		response := historyResponse{translations}

		return response, nil
	}
}*/
