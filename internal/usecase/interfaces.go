package usecase

import (
	"github.com/deniskaponchik/GoSoft/internal/entity"
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

	PolyWebApi interface {
		ApiLineInfoErr(*entity.PolyStruct) error //(entity.PolyStruct, error) //string, error)
		ApiSafeRestart(entity.PolyStruct) error
	}

	PolyNetDial interface {
		NetDialTmtErr(*entity.PolyStruct) error //entity.PolyStruct, error)
	}
)

type (
	RepoInterface interface {
		//UnifiRepo
		//PolyRepo
	}
	GisupRepo interface {
		InsertOffice(*entity.Office) error
		UpdateOfficeLogin(string, string) error
		UpdateOfficeException(string, string) error

		UpdateDbAnomaly(map[string]*entity.Anomaly) error
		UpdateDbClient(map[string]*entity.Client) error
		UpdateDbAp(map[string]*entity.Ap) error
		UpdateDbOffice(map[string]*entity.Office) error

		DownloadMacMapsClientApWithAnomaly(map[string]*entity.Client, map[string]*entity.Ap, string, time.Time) error
		Download2MapFromDBclient() (map[string]*entity.Client, map[string]*entity.Client, error)
		Download2MapFromDBaps() (map[string]*entity.Ap, map[string]*entity.Ap, error)
		DownloadMapOffice() (map[string]*entity.Office, error)

		GetLoginPCerr(*entity.Client) (err error)

		DownloadMapFromDBvcsErr(int) (map[string]entity.PolyStruct, error)
		UpdateMapsToDBerr(map[string]entity.PolyStruct) error
	}
	GisupRestIn interface {
		CheckToken(string) (string, error)
		GetToken(*entity.User) (string, error)
		CheckUser(*entity.User) error //string, string) error

		OfficeSapcnChange(string, string) error
		OfficeLoginChange(string, string) error
		OfficeTimeZoneChange(string, string) error
		OfficeExceptionChange(string, string) error
		OfficeNew(*entity.Office) error
		GetSapcnSortSliceForAdminkaPage() []string

		GetClientForRest(string) *entity.Client
		GetApForRest(string) *entity.Ap
	}

	//Repo interface { //UnifiRepo	}

	GisupSoap interface {
		CreateTicketSmacWifi(ticket *entity.Ticket) (err error)
		CreateTicketSmacVcs(ticket *entity.Ticket) (err error)
		//CreatePolyTicketErr(*entity.Ticket) error
		CheckTicketStatusErr(ticket *entity.Ticket) (err error)
		ChangeStatusErr(ticket *entity.Ticket) (err error)
		AddCommentErr(ticket *entity.Ticket) (err error)
	}
	Authorization interface {
		GenerateToken(*entity.User) (string, error)
		ParseToken(string) (string, error)
	}
	Authentication interface {
		AuthSecur(user *entity.User) error
	}
	GisupRmq interface {
		Publish(message, queueName string) error
	}
	//исходящие rest запросы
	GisupC3po interface { //UnifiRestOut
		GetUserLogin(*entity.Client) error
		//getPc(*entity.Client) error
	}
)

type (
	//implement usecase methods to web
	UnifiRestIn interface {
		/*
			CheckToken(string) (string, error)
			GetToken(*entity.User) (string, error)
			CheckUser(*entity.User) error //string, string) error

			OfficeSapcnChange(string, string) error
			OfficeLoginChange(string, string) error
			OfficeTimeZoneChange(string, string) error
			OfficeExceptionChange(string, string) error
			OfficeNew(*entity.Office) error
			GetSapcnSortSliceForAdminkaPage() []string
		*/
		GetClientForRest(string) *entity.Client
		GetApForRest(string) *entity.Ap
	}

	Ui interface {
		GetSites() error
		AddAps2Maps(map[string]*entity.Ap, map[string]*entity.Ap) error
		//AddAps(map[string]*entity.Ap) error
		UpdateClients2MapWithoutApMap(map[string]*entity.Client, map[string]*entity.Client, string) error
		GetHourAnomaliesAddSlice(string, map[string]*entity.Client, map[string]*entity.Ap) (map[string]*entity.Anomaly, error)
		//GetHourAnomaliesAddSlice(map[string]*entity.Client, map[string]*entity.Ap) (map[string]*entity.Anomaly, error)
		//GetHourAnomalies(map[string]*entity.Client, map[string]*entity.Ap) (map[string]*entity.Anomaly, error)
	}
)
