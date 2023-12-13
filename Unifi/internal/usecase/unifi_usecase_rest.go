package usecase

import (
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"sort"
)

func (uuc *UnifiUseCase) GetSapcnSortSliceForAdminkaPage() []string {
	lenSapcnMap := len(siteApCutName_Office)
	sortKeys := make([]string, lenSapcnMap)
	i := 0
	for k := range siteApCutName_Office {
		sortKeys[i] = k
		i++
	}
	sort.Strings(sortKeys)
	return sortKeys
}

func (uuc *UnifiUseCase) GetClientForRest(hostName string) *entity.Client { //c context.Context
	uuc.mx.RLock()
	defer uuc.mx.RUnlock()
	client, exisHost := uuc.hostnameClient[hostName]
	if exisHost {
		return client
	} else {
		return nil
	}
}

func (uuc *UnifiUseCase) GetApForRest(hostName string) *entity.Ap { //c context.Context
	uuc.mx.RLock()
	defer uuc.mx.RUnlock()
	ap, exisHost := uuc.hostnameAp[hostName]
	if exisHost {
		return ap
	} else {
		return nil
	}
}
