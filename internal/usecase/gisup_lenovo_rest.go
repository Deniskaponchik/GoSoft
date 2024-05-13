package usecase

func (luc *LenovoUseCase) NewAlarm(zabbixAlarm, lenovoVcs) {
	//при поступлении нового аларма

	//1. найти устройство по Имени переговорной в мапе ВКС

	luc.mx.RLock()
	defer luc.mx.RUnlock()
	vcs, exisHost := luc.hostnameVcs[lenovoVcs.hostname]
	if exisHost {

	} else {

	}

}

func (luc *LenovoUseCase) CloseAlarm(zabbixAlarm, lenovoVcs) {
	//при закрытии старого аларма

}
