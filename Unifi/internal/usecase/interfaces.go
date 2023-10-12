package usecase

import (
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
)

type (
	PolyInterface interface {
	}

	PolyRepo interface {
		UpdateMapsToDBerr(string, []string)
		UploadMapsToDBerr(map[string]entity.PolyStruct) error
		DownloadMapFromDBvcsErr(string) (map[string]entity.PolyStruct, error)
	}

	PolySoap interface {
		CreatePolyTicketErr(entity.PolyStruct) ([]string, error)
		CheckTicketStatusErr(entity.PolyStruct) (string, error)
		ChangeStatusErr(entity.PolyStruct) (string, error)
		AddCommentErr(entity.PolyStruct) (string, error)
	}

	PolyWebApi interface {
		ApiLineInfo(entity.PolyStruct) (string, error)
		ApiSafeRestart2(entity.PolyStruct) (string, error)
	}

	PolyNetDial interface {
		NetDialTmtErr(entity.PolyStruct) (string, error)
	}
)
