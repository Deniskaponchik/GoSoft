package main

func main() {
	anomalies := []string{
		"anomal1",
		"anomaly2",
		"anomaly3",
	}

	//SoapCreateTicket(usrLogin, v.clientName, v.corpAnomalies, v.apName, v.site)
	SoapCreateTicket("dmirty.pushkarev", "NBKH-PUSHKAREV", anomalies, "IRK-CO-1FL", "Иркутск")
}
