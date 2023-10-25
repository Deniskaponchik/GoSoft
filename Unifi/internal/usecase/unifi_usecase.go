package usecase

import (
	"fmt"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
)

type UnifiUseCase struct {
	repo         UnifiRepo //interface
	soap         UnifiSoap //interface
	ui           Ui        //interface
	everyCodeMap map[int]bool
	//restartHour  int
}

// реализуем Инъекцию зависимостей DI. Используется в app
func NewUnifiUC(r UnifiRepo, s UnifiSoap, ui Ui, everyCode map[int]bool) *UnifiUseCase {
	return &UnifiUseCase{
		//Мы можем передать сюда ЛЮБОЙ репозиторий (pg, s3 и т.д.) НО КОД НЕ ПОМЕНЯЕТСЯ! В этом смысл DI
		repo:         r,
		soap:         s,
		ui:           ui,
		everyCodeMap: everyCode,
		//restartHour:  restartHour,
	}
}

// Переменные, которые используются во всех методах ниже
var apMap map[string]entity.Ap
var machineMap map[string]entity.Client  //client = machine
var anomalyMap map[string]entity.Anomaly
var region_unifiSlice map[string][]entity.
var err error

func (puc *UnifiUseCase) InfinityProcessingUnifi() (error) {
	//ubiq.Ui.
	err = puc.ui.GetSites()
	if err == nil {

	}else{
		fmt.Println("sites НЕ загрузились")
		fmt.Println(err.Error())
	}
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
