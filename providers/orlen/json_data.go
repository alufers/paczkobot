package orlen

// OrlenResponse is the root object in the JSONP response from Orlen
type OrlenResponse struct {
	Status      string         `json:"status"`
	Number      string         `json:"number"`
	Full        bool           `json:"full"`
	HistoryHTML string         `json:"historyHtml"`
	History     []HistoryEntry `json:"history"`
	Label       string         `json:"label"`
	Return      bool           `json:"return"`
	TruckNo     string         `json:"truckNo"`
	ReturnTruck string         `json:"returnTruck"`
}

type HistoryEntry struct {
	Date       string `json:"date"`
	Code       string `json:"code"`
	Label      string `json:"label"`
	LabelShort string `json:"labelShort"`
}
