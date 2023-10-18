package usecase

import (
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
)

type (
	PolyInterface interface {
		InfinityPolyProcessing() error
		Survey() error          //map[string]entity.PolyStruct, map[string][]entity.PolyStruct, error)
		TicketsCreating() error //map[string]entity.PolyStruct, map[string][]entity.PolyStruct
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
		ApiSafeRestart(entity.PolyStruct) error
	}

	PolyNetDial interface {
		NetDialTmtErr(entity.PolyStruct) (entity.PolyStruct, error)
	}
)

type (
	UnifiInterface interface {
		InfinityUnifiProcessing() error
		ApSurvey() error
		ApTicketCtreating() error
		ClientSurvey() error
		ClientTicketCreating() error
	}
	UnifiRepo interface {
		//DownloadMapFromDB
	}
	UnifiSoap interface {
	}
	UnifiUnifi interface {
	}
)
