package ubiq

import (
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"github.com/deniskaponchik/GoSoft/Unifi/pkg/unpoller"
)

type Ubiq struct {
	//username string
	//password string
	//*unifi.Config
	unp *unpoller.Unpoller
}

func NewUbiq(up *unpoller.Unpoller) *Ubiq {
	return &Ubiq{
		unp: up,
	}
}

func (ubq *Ubiq) GetApSlice() ([]entity.Ap, error) {
	//return &UnifiUnifi{}
	return nil, nil
}

func (ubq *Ubiq) GetAnomalySlice() ([]entity.Anomaly, error) {
	//return &UnifiUnifi{}
	return nil, nil
}

func (ubq *Ubiq) GetClientSlice() ([]entity.Client, error) {
	return nil, nil
}
