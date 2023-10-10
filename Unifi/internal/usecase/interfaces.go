package usecase

import (
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
)

type (
	Polycom interface {
	}

	PolyRepo interface {
		UpdateMapsToDBerr(string, []string)
		UploadMapsToDBerr
		DownloadMapFromDBvcsErr(string) map[string]entity.PolyStruct
	}

	PolyWebApi interface {
	}
)
