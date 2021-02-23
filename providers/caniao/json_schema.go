package caniao

type CaniaoJSONRoot struct {
	Data        []Data  `json:"data"`
	Success     bool    `json:"success"`
	TimeSeconds float64 `json:"timeSeconds"`
}
type LatestTrackingInfo struct {
	Desc     string `json:"desc"`
	Status   string `json:"status"`
	Time     string `json:"time"`
	TimeZone string `json:"timeZone"`
}
type Section1 struct {
	CountryName string        `json:"countryName"`
	DetailList  []interface{} `json:"detailList"`
}
type DetailList struct {
	Desc     string `json:"desc"`
	Status   string `json:"status"`
	Time     string `json:"time"`
	TimeZone string `json:"timeZone"`
}
type Section2 struct {
	CountryName string       `json:"countryName"`
	DetailList  []DetailList `json:"detailList"`
}
type Data struct {
	AllowRetry         bool               `json:"allowRetry"`
	BizType            string             `json:"bizType"`
	CachedTime         string             `json:"cachedTime"`
	DestCountry        string             `json:"destCountry"`
	DestCpList         []interface{}      `json:"destCpList"`
	HasRefreshBtn      bool               `json:"hasRefreshBtn"`
	LatestTrackingInfo LatestTrackingInfo `json:"latestTrackingInfo"`
	MailNo             string             `json:"mailNo"`
	OriginCountry      string             `json:"originCountry"`
	OriginCpList       []interface{}      `json:"originCpList"`
	Section1           Section1           `json:"section1"`
	Section2           Section2           `json:"section2"`
	ShippingTime       float64            `json:"shippingTime"`
	ShowEstimateTime   bool               `json:"showEstimateTime"`
	Status             string             `json:"status"`
	StatusDesc         string             `json:"statusDesc"`
	Success            bool               `json:"success"`
	ErrorCode          string             `json:"errorCode""`
	SyncQuery          bool               `json:"syncQuery"`
}
