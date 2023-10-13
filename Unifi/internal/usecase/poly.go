package usecase

import (
	//Not have package imports from the outer layer.
	"fmt"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
)

type PolyUseCase struct {
	repo    PolyRepo    //interface
	webAPI  PolyWebApi  //interface
	netDial PolyNetDial //interface
	soap    PolySoap    //interface
}

// реализуем Инъекцию зависимостей DI. Используется в app
func New(r PolyRepo, a PolyWebApi, n PolyNetDial, s PolySoap) *PolyUseCase {
	return &PolyUseCase{
		repo:    r,
		webAPI:  a,
		netDial: n,
		soap:    s,
	}
}

// Получение списка устройств
func (puc *PolyUseCase) GetEntityMap() (map[string]entity.PolyStruct, error) {

}

// Опрос устройств
func (puc *PolyUseCase) Survey(polyMap map[string]entity.PolyStruct) error {

	srStatusCodesForNewTicket := map[string]bool{
		"Отменено":     true, //Cancel  6e5f4218-f46b-1410-fe9a-0050ba5d6c38
		"Решено":       true, //Resolve  ae7f411e-f46b-1410-009b-0050ba5d6c38
		"Закрыто":      true, //Closed  3e7f420c-f46b-1410-fc9a-0050ba5d6c38
		"На уточнении": true, //Clarification 81e6a1ee-16c1-4661-953e-dde140624fb
		"Тикет введён не корректно": true,
		//"": true,
	}
	srStatusCodesForCancelTicket := map[string]bool{
		"Визирование":  true,
		"Назначено":    true,
		"На уточнении": true, //Clarification 81e6a1ee-16c1-4661-953e-dde140624fb
	}

	//Make calls to the outer layer through the interface (!).
	polyMap, err := uc.repo.DownloadMapFromDBvcsErr
	if err != nil {
		return nil, fmt.Errorf("TranslationUseCase - History - s.repo.GetHistory: %w", err)
	}
	return nil
}

// Создание заявок
func (puc *PolyUseCase) Ticketing() (polyTicket entity.PolyTicket, err error)

//Перезагрузка устройств
