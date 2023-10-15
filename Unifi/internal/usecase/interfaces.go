package usecase

import (
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
)

type (
	PolyInterface interface {
		GetEntityMap(int) (map[string]entity.PolyStruct, error)
		Survey(map[string]entity.PolyStruct) (map[string][]entity.PolyStruct, error)
		TicketsCreating(map[string][]entity.PolyStruct) (map[string]entity.PolyStruct, error)
	}

	PolyRepo interface {
		DownloadMapFromDBvcsErr(int) (map[string]entity.PolyStruct, error)
		UpdateMapsToDBerr(map[string]entity.PolyStruct) error
	}

	PolySoap interface {
		CreatePolyTicketErr(entity.PolyTicket) (entity.PolyTicket, error) //[]string, error)
		CheckTicketStatusErr(entity.PolyTicket) (entity.PolyTicket, error)
		ChangeStatusErr(entity.PolyTicket) error
		AddCommentErr(entity.PolyTicket) error
	}

	PolyWebApi interface {
		ApiLineInfoErr(entity.PolyStruct) (entity.PolyStruct, error) //string, error)
		ApiSafeRestart2(entity.PolyStruct) (string, error)
	}

	PolyNetDial interface {
		NetDialTmtErr(entity.PolyStruct) (entity.PolyStruct, error)
	}
)
