// Package entity defines main entities for business logic (services), database mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

type Ap struct {
	Mac            string     `json:"mac"        example:"a0-b1-c2-d3-e4-f5"`
	SiteName       string     `json:"region"     example:"Волгоград"`
	SiteID         string     `example:"5e74aaa6a1a76964e770815c"` //уточнить, нужен ли
	Name           string     `json:"name"       example:"XXX-OPENSPACE"`
	UserLogin      string     `json:"login"      example:"vasya.pupkin"`
	SrID           string     `json:"srid"       example:"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"`
	Exception      int        `json:"exception"  example:"1"` //Исключение для аномалий клиентов, а не для отключений точек
	Controller     int        `json:"controller" example:"1"` //
	CommentCount   int        `example:"1"`                   //0 - нет комментариев, 1 - комментарий "точка появилась в сети", 2 - Попытка закрыть обращение. commentForUpdate
	StateInt       int        `example:"0"`                   // 0 - available
	CountAttempts  int        `example:"0"`                   // Число заходов на создание заявок. на втором заходе создаём тикет. не берётся из БД
	SliceAnomalies []*Anomaly //Аномалии точки за 30 дней
	Date30count    int        `example:"27"` //Используется в DownloadClientsWithAnomalySlice
	//SliceClients   []*Client           //Аномалии точки за 30 дней
	Ticket       *Ticket
	CountAnomaly int //кол-во аномалий за последние 30 дней
}

type Client struct {
	Mac      string `json:"mac_client" example:"a0:b1:c2:d3:e4:f5"`
	Hostname string `json:"hostname"   example:"XXXX-PUPKIN"`
	//SiteName нужен только на этапе создания заявок по клиентам. Поэтому при обработке каждого клиента его не получаю.
	SiteName   string `json:"sitename"   example:"Москва"`
	SrID       string `json:"srid"       example:"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"`
	Controller int    `json:"controller" example:"1"`
	Exception  int    `json:"exception"  example:"1"`
	ApName     string `json:"ap_name"    example:"XXX-OPENSPACE"` //отключаю, чтобы не было неразберихи. не заполянется этот параметр на контроллере
	ApMac      string `json:"ap_mac"     example:"a0:b1:c2:d3:e4:f5"`
	Modified   string `json:"modified"   example:"2023-10-28"`
	UserLogin  string `example:"vasya.pupkin"` //не однозначная характеристика. нигде не используется

	//DateTicketCreateAttempt time.Time `example: "2023-10-28"` //До первого захода либо nil, либо прошлая дата, после 1 захода - сегодняшняя дата
	DateTicketCreateAttempt int                 `example:"28"` //Используется в DownloadMacClientsWithAnomalies
	Date30count             int                 `example:"27"` //Используется в DownloadClientsWithAnomalySlice
	SliceAnomalies          []*Anomaly          //Аномалии клиента за 30 дней
	Date_Anomaly            map[string]*Anomaly //Аномалии клиента за 30 дней. Вроде, не должна больше использоваться
	CountAnomaly            int                 //кол-во аномалий за последние 30 дней

	//TODO: hardware from c3po
	//serial number
	//disk0
	//etc
}

// Структура должна обнуляться каждый час и при выгрузке раз в сутки
type Anomaly struct {
	ClientMac  string `json:"mac_client" example:"a0:b1:c2:d3:e4:f5"`
	ClientName string
	SiteName   string `json:"sitename"   example:"Москва"`
	Controller int    `json:"controller" example:"1"`

	//при обработке каждой аномалии подключаюсь к мапе Клиентов.
	//А при обработке каждого клиента подключаюсь к мапе точек, чтобы была актальная инфа по Exception
	//Поэтому каждый раз получаю актуальный: имя точки, мак точки, сумму исключений точки и клиента
	ApName    string `json:"name_ap"     example:"XXX-FL1-01-OPENSPACE"`
	ApMac     string `json:"mac_ap"      example:"68:d7:9a:1c:f2:b9"`
	Exception int    `json:"exception"   example:"1"` //берётся от Client. 2 = exception from Ap and Client

	//AnomalySlice []string  `json:"anomalies"  example:"USER_HIGH_TCP_LATENCY;USER_LOW_PHY_RATE;USER_SLEEPY_CLIENT;USER_HIGH_TCP_PACKET_LOSS;USER_HIGH_WIFI_RETRIES;USER_SIGNAL_STRENGTH_FAILURES;USER_DNS_TIMEOUT;USER_HIGH_WIFI_LATENCY;USER_POOR_STREAM_EFF;USER_HIGH_DNS_LATENCY"`
	AnomStr      string `json:"anomalies"   example:"USER_HIGH_TCP_LATENCY;USER_LOW_PHY_RATE;USER_SLEEPY_CLIENT;USER_HIGH_TCP_PACKET_LOSS;"`
	SliceAnomStr []string
	//TimeStr_sliceAnomStr map[string][]string //day - 2023-09-01, hour - 2023-09-01 12:00:00

	DateHour string `json:"date_hour"  example:"2023-09-01 12:00:00"`
	//DateHour     time.Time `json:"date_hour"  example:"2023-09-01 12:00:00"`
}

type User struct {
	Login      string
	Password   string
	Sid        string
	FIO        string `example:"Иванов Иван Иванович"`
	GivenName  string `example:"Иван"`
	MiddleName string `example:"Иванович"`
	SurName    string `example:"Иванов"`
	//PCs   []*Client
}
