package usecase

import (
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
)

type UnifiUseCase struct {
	repo  UnifiRepo //interface
	soap  UnifiSoap //interface
	unifi UnifiUnifi

	everyCodeMap map[int]bool
	//restartHour  int
}

// реализуем Инъекцию зависимостей DI. Используется в app
func NewUnifi(r UnifiRepo, u UnifiUnifi, s UnifiSoap, everyCode map[int]bool) *UnifiUseCase {
	return &UnifiUseCase{
		//Мы можем передать сюда ЛЮБОЙ репозиторий (pg, s3 и т.д.) НО КОД НЕ ПОМЕНЯЕТСЯ! В этом смысл DI
		repo:         r,
		unifi:        u,
		soap:         s,
		everyCodeMap: everyCode,
		//restartHour:  restartHour,
	}
}

// Переменные, которые используются во всех методах ниже
var aaaaaMap map[string]entity.PolyStruct
var region_unifiSlice map[string][]entity.PolyStruct
var erru error

func (puc *UnifiUseCase) InfinityProcessingUnifi() error {

	return nil
}

// Опрос устройств
func (puc *UnifiUseCase) Survey() error {

	return nil
}

// Создание заявок
func (puc *UnifiUseCase) TicketsCreating() error {

	return nil
}

//Перезагрузка устройств
