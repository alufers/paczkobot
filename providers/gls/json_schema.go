package gls

type Response struct {
	TuStatus []TuStatus `json:"tuStatus"`
}
type References struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}
type Owners struct {
	Type string `json:"type"`
	Code string `json:"code"`
}
type Infos struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}
type ArrivalTime struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type Address struct {
	CountryCode string `json:"countryCode"`
	City        string `json:"city"`
	CountryName string `json:"countryName"`
}
type History struct {
	Time    string  `json:"time"`
	Date    string  `json:"date"`
	Address Address `json:"address"`
	EvtDscr string  `json:"evtDscr"`
}
type StatusBar struct {
	Status      string `json:"status"`
	StatusText  string `json:"statusText"`
	ImageStatus string `json:"imageStatus"`
	ImageText   string `json:"imageText"`
}
type ProgressBar struct {
	Level       int         `json:"level"`
	StatusText  string      `json:"statusText"`
	StatusBar   []StatusBar `json:"statusBar"`
	RetourFlag  bool        `json:"retourFlag"`
	EvtNos      []string    `json:"evtNos"`
	ColourIndex int         `json:"colourIndex"`
	StatusInfo  string      `json:"statusInfo"`
}
type TuStatus struct {
	References             []References `json:"references"`
	Owners                 []Owners     `json:"owners"`
	Infos                  []Infos      `json:"infos"`
	DeliveryOwnerCode      string       `json:"deliveryOwnerCode"`
	ArrivalTime            ArrivalTime  `json:"arrivalTime"`
	TuNo                   string       `json:"tuNo"`
	History                []History    `json:"history"`
	ProgressBar            ProgressBar  `json:"progressBar"`
	ChangeDeliveryPossible bool         `json:"changeDeliveryPossible"`
}

