package usecase

import (
	"errors"
	"github.com/deniskaponchik/GoSoft/internal/entity"
	"sort"
	"strconv"
)

func (uuc *UnifiUseCase) OfficeSapcnChange(oldSapcn string, newSapcn string) error {
	return nil
}

func (uuc *UnifiUseCase) OfficeLoginChange(sapcn string, newLogin string) error {
	office := siteApCutName_Office[sapcn]
	if office.UserLogin != newLogin {
		err = uuc.repo.UpdateOfficeLogin(sapcn, newLogin)
		if err == nil {
			office.UserLogin = newLogin
			return nil
		} else {
			return err
		}
	} else {
		return errors.New("новый логин соответствует старому")
	}
}

func (uuc *UnifiUseCase) OfficeTimeZoneChange(sapcn string, newTimeZone string) error {
	return nil
}

func (uuc *UnifiUseCase) OfficeExceptionChange(sapcn string, newException string) error {
	office := siteApCutName_Office[sapcn]

	newExceptionInt, errStrConv := strconv.Atoi(newException)
	if errStrConv != nil {

		if office.Exception != newExceptionInt {
			err = uuc.repo.UpdateOfficeException(sapcn, newException)
			if err == nil {
				office.Exception = newExceptionInt
				return nil
			} else {
				return err
			}
		} else {
			return errors.New("новый Exception соответствует старому")
		}
	} else {
		return errStrConv
	}
}

func (uuc *UnifiUseCase) OfficeNew(newOffice *entity.Office) error { //newSapcn string, login string
	_, exist := siteApCutName_Office[newOffice.Site_ApCutName]
	if exist == false {
		err = uuc.repo.InsertOffice(newOffice)
		if err == nil {
			siteApCutName_Office[newOffice.Site_ApCutName] = newOffice
			return nil
		} else {
			return err
		}
	} else {
		return errors.New("Офис уже существует")
	}
}

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
