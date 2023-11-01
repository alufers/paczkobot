package dpd_eu

type DpdEuResponse struct {
	ParcellifecycleResponse *ParcellifecycleResponse `json:"parcellifecycleResponse"`
}
type ServiceElements struct {
	Label   string   `json:"label"`
	Content []string `json:"content"`
}
type AdditionalProperties struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type ShipmentInfo struct {
	ParcelLabelNumber       string                 `json:"parcelLabelNumber"`
	ServiceElements         []ServiceElements      `json:"serviceElements"`
	SortingCode             string                 `json:"sortingCode"`
	ProductName             string                 `json:"productName"`
	CodInformationAvailable bool                   `json:"codInformationAvailable"`
	Documents               []any                  `json:"documents"`
	AdditionalProperties    []AdditionalProperties `json:"additionalProperties"`
}
type Description struct {
	Content []string `json:"content"`
}
type Depot struct {
	BusinessUnit string `json:"businessUnit"`
	Number       string `json:"number"`
}
type StatusInfo struct {
	Status               string      `json:"status"`
	Label                string      `json:"label"`
	Description          Description `json:"description"`
	StatusHasBeenReached bool        `json:"statusHasBeenReached"`
	IsCurrentStatus      bool        `json:"isCurrentStatus"`
	Location             string      `json:"location,omitempty"`
	Depot                Depot       `json:"depot,omitempty"`
	Date                 string      `json:"date,omitempty"`
	NormalItems          []any       `json:"normalItems"`
	ImportantItems       []any       `json:"importantItems"`
	ErrorItems           []any       `json:"errorItems"`
}
type ScanDepot struct {
	BusinessUnit string `json:"businessUnit"`
	Number       string `json:"number"`
}
type ScanType struct {
	Name       string `json:"name"`
	Code       string `json:"code"`
	DetailMode string `json:"detailMode"`
}
type AdditionalCodes struct {
	AdditionalCode []any `json:"additionalCode"`
}

type ScanDescription struct {
	Label   string   `json:"label"`
	Content []string `json:"content"`
}
type DestinationDepot struct {
	BusinessUnit string `json:"businessUnit"`
	Number       string `json:"number"`
}
type ServiceElements2 struct {
	Code       string `json:"code"`
	ShortName  string `json:"shortName"`
	DetailMode string `json:"detailMode"`
}
type ScanData struct {
	ScanDate         string             `json:"scanDate"`
	ScanTime         string             `json:"scanTime"`
	ScanDepot        ScanDepot          `json:"scanDepot"`
	DestinationDepot DestinationDepot   `json:"destinationDepot"`
	Location         string             `json:"location"`
	ScanType         ScanType           `json:"scanType"`
	Route            string             `json:"route"`
	AdditionalCodes  AdditionalCodes    `json:"additionalCodes"`
	Country          string             `json:"country"`
	ServiceElements  []ServiceElements2 `json:"serviceElements"`
	SortingCode      string             `json:"sortingCode"`
	Tour             string             `json:"tour"`
	DetailMode       string             `json:"detailMode"`
	InsertTimestamp  string             `json:"insertTimestamp"`
}
type Scan struct {
	Date            string          `json:"date"`
	IntegrationDate string          `json:"integrationDate"`
	ScanData        ScanData        `json:"scanData,omitempty"`
	AdditionalCodes []any           `json:"additionalCodes"`
	ScanDescription ScanDescription `json:"scanDescription"`
	Links           []any           `json:"links"`
	ProductName     string          `json:"productName,omitempty"`
}
type ScanInfo struct {
	Scan []Scan `json:"scan"`
}
type ParcelLifeCycleData struct {
	ShipmentInfo ShipmentInfo `json:"shipmentInfo"`
	StatusInfo   []StatusInfo `json:"statusInfo"`
	ContactInfo  []any        `json:"contactInfo"`
	ScanInfo     ScanInfo     `json:"scanInfo"`
}
type ParcellifecycleResponse struct {
	ParcelLifeCycleData ParcelLifeCycleData `json:"parcelLifeCycleData"`
}
