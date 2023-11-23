package usecase

import (
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"time"
	//"context"
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
		CreatePolyTicketErr(*entity.Ticket) error  //[]string, error)
		CheckTicketStatusErr(*entity.Ticket) error //(entity.Ticket, error)
		ChangeStatusErr(*entity.Ticket) error
		AddCommentErr(*entity.Ticket) error
	}

	PolyWebApi interface {
		ApiLineInfoErr(*entity.PolyStruct) error //(entity.PolyStruct, error) //string, error)
		ApiSafeRestart(entity.PolyStruct) error
	}

	PolyNetDial interface {
		NetDialTmtErr(*entity.PolyStruct) error //entity.PolyStruct, error)
	}
)

type (
	//UnifiRest interface {	}
	UnifiInterface interface {
		GetClientForRest(string) *entity.Client //, error) //context.Context
		InfinityProcessingUnifi()               //error
		HandlingAps() (map[string][]*entity.Ap, error)
		TicketsCreatingAps(map[string][]*entity.Ap) error
		TicketsCreatingMacClients(map[string]*entity.Client) error
		//TicketsCreatingAnomalies(map[string][]*entity.Client) error
	}
	UnifiRepo interface {
		UpdateDbAnomaly(map[string]*entity.Anomaly) error
		UpdateDbClient(map[string]*entity.Client) error
		UpdateDbAp(map[string]*entity.Ap) error
		UploadMapsToDBerr(string) error
		DownloadMacClientsWithAnomalies(map[string]*entity.Client, string, time.Time) error
		Download2MapFromDBclient() (map[string]*entity.Client, map[string]*entity.Client, error)
		DownloadMapFromDBmachinesErr() (map[string]*entity.Client, error)
		DownloadMapFromDBapsErr() (map[string]*entity.Ap, error)
		DownloadMapOffice() (map[string]*entity.Office, error)
		//DownloadMapFromDBerr() (map[string]string, error)
		GetLoginPCerr(*entity.Client) (err error)
	}
	UnifiSoap interface {
		CreateTicketSmacWifi(ticket *entity.Ticket) (err error)
		CreateTicketSmacVcs(ticket *entity.Ticket) (err error)
		CheckTicketStatusErr(ticket *entity.Ticket) (err error)
		ChangeStatusErr(ticket *entity.Ticket) (err error)
		AddCommentErr(ticket *entity.Ticket) (err error)
	}
	Ui interface {
		GetSites() error
		AddAps(map[string]*entity.Ap) error
		Update2MapClientsWithoutApMap(map[string]*entity.Client, map[string]*entity.Client, string) error
		//UpdateClientsWithoutApMap(map[string]*entity.Client, string) error
		GetHourAnomalies(map[string]*entity.Client, map[string]*entity.Ap) (map[string]*entity.Anomaly, error)
	}
)
