package usecase

import (
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
)

type PolySurveyUseCase struct {
	repo   PolyRepo   //interface
	webAPI PolyWebAPI //interface
	http   PolyHTTP   //interface
	//ping  PolyPing  //interface
}

func New(r PolyRepo, w PolyWebApi, h PolyHTTP) {
	return &PolySurveyUseCase{
		repo:   r,
		webAPI: w,
		http:   h,
	}
}

//Создать метод Survey

func (psuc *PolySurveyUseCase) GetPolyStructMapFromDB() (map[string]entity.PolyStruct{}, error) {
	polyMap := map[string]entity.PolyStruct{}
	polyMap, err := uc.repo.DownloadMapFromDBvcsErr
	if err != nil {
		return nil, fmt.Errorf("TranslationUseCase - History - s.repo.GetHistory: %w", err)
	}
	return polyMap, nil
}

// Создать метод опроса Codec


//Создать метод опроса Visual

