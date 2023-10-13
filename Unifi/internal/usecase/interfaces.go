package usecase

import (
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
)

type (
	PolyInterface interface {
		Survey(map[string]entity.PolyStruct) error
		Ticketing
	}

	PolyRepo interface {
		UpdateMapsToDBerr(map[string]entity.PolyStruct) error
		DownloadMapFromDBvcsErr(int) (map[string]entity.PolyStruct, error)
	}

	PolySoap interface {
		CreatePolyTicketErr(entity.PolyTicket) (entity.PolyTicket, error) //[]string, error)
		CheckTicketStatusErr(entity.PolyTicket) (entity.PolyTicket, error)
		ChangeStatusErr(entity.PolyTicket) (entity.PolyTicket, error)
		AddCommentErr(entity.PolyTicket) error
	}

	PolyWebApi interface {
		ApiLineInfo(entity.PolyStruct) (string, error)
		ApiSafeRestart2(entity.PolyStruct) (string, error)
	}

	PolyNetDial interface {
		NetDialTmtErr(entity.PolyStruct) (string, error)
	}
)
