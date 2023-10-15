package cainiao

type CainiaoResponse struct {
	Module  []Module `json:"module"`
	Success bool     `json:"success"`
}

type GetCityResponse struct {
	Module  string `json:"module"`
	Success bool   `json:"success"`
}

type ProgressPointList struct {
	PointName string `json:"pointName"`
	Light     bool   `json:"light,omitempty"`
	Reload    bool   `json:"reload,omitempty"`
}
type ProcessInfo struct {
	ProgressStatus    string              `json:"progressStatus"`
	ProgressRate      float64             `json:"progressRate"`
	Type              string              `json:"type"`
	ProgressPointList []ProgressPointList `json:"progressPointList"`
}
type GlobalEtaInfo struct {
	EtaDesc         string `json:"etaDesc"`
	DeliveryMinTime int64  `json:"deliveryMinTime"`
	DeliveryMaxTime int64  `json:"deliveryMaxTime"`
}
type Group struct {
	NodeCode       string `json:"nodeCode"`
	NodeDesc       string `json:"nodeDesc"`
	CurrentIconURL string `json:"currentIconUrl"`
	HistoryIconURL string `json:"historyIconUrl"`
}
type GlobalCombinedLogisticsTraceDTO struct {
	Time         int64  `json:"time"`
	TimeStr      string `json:"timeStr"`
	Desc         string `json:"desc"`
	StanderdDesc string `json:"standerdDesc"`
	DescTitle    string `json:"descTitle"`
	TimeZone     string `json:"timeZone"`
	ActionCode   string `json:"actionCode"`
	Group        Group  `json:"group"`
}
type LatestTrace struct {
	Time         int64  `json:"time"`
	TimeStr      string `json:"timeStr"`
	Desc         string `json:"desc"`
	StanderdDesc string `json:"standerdDesc"`
	DescTitle    string `json:"descTitle"`
	TimeZone     string `json:"timeZone"`
	ActionCode   string `json:"actionCode"`
	Group        Group  `json:"group"`
}
type DetailListItem struct {
	Time         int64  `json:"time"`
	TimeStr      string `json:"timeStr"`
	Desc         string `json:"desc"`
	StanderdDesc string `json:"standerdDesc"`
	DescTitle    string `json:"descTitle"`
	TimeZone     string `json:"timeZone"`
	ActionCode   string `json:"actionCode"`
	Group        Group  `json:"group,omitempty"`
}
type Module struct {
	MailNo                          string                          `json:"mailNo"`
	OriginCountry                   string                          `json:"originCountry"`
	DestCountry                     string                          `json:"destCountry"`
	MailType                        string                          `json:"mailType"`
	MailTypeDesc                    string                          `json:"mailTypeDesc"`
	Status                          string                          `json:"status"`
	StatusDesc                      string                          `json:"statusDesc"`
	MailNoSource                    string                          `json:"mailNoSource"`
	ProcessInfo                     ProcessInfo                     `json:"processInfo"`
	GlobalEtaInfo                   GlobalEtaInfo                   `json:"globalEtaInfo"`
	GlobalCombinedLogisticsTraceDTO GlobalCombinedLogisticsTraceDTO `json:"globalCombinedLogisticsTraceDTO"`
	LatestTrace                     LatestTrace                     `json:"latestTrace"`
	DetailList                      []DetailListItem                `json:"detailList"`
	DaysNumber                      string                          `json:"daysNumber"`
}
