package usecase

import (
	//Not have package imports from the outer layer.
	"fmt"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
)

type PolyUseCase struct {
	repo   PolyRepo    //interface
	webAPI PolyWebApi  //interface
	http   PolyNetDial //interface
	soap   PolySoap    //interface
}

// реализуем Инъекцию зависимостей DI. Используется в app
func New(r PolyRepo, a PolyWebApi, h PolyNetDial, s PolySoap) *PolyUseCase {
	return &PolyUseCase{
		repo:   r,
		webAPI: a,
		http:   h,
		soap:   s,
	}
}

//Создать метод Survey

func (psuc *PolyUseCase) GetPolyStructMapFromDB() (map[string]entity.PolyStruct, error) {
	polyMap := map[string]entity.PolyStruct{}
	//Make calls to the outer layer through the interface (!).
	polyMap, err := uc.repo.DownloadMapFromDBvcsErr
	if err != nil {
		return nil, fmt.Errorf("TranslationUseCase - History - s.repo.GetHistory: %w", err)
	}
	return polyMap, nil
}

// Создать метод опроса Codec

//Создать метод опроса Visual
