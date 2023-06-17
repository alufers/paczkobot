package inpostextra

import "time"

type MultiCompartment struct {
	Collected       bool     `json:"collected,omitempty"`
	Presentation    bool     `json:"presentation,omitempty"`
	ShipmentNumbers []string `json:"shipmentNumbers,omitempty"`
	Uuid            string   `json:"uuid,omitempty"`
}

type GetTrackedParcelsResponse struct {
	UpdatedUntil time.Time      `json:"updatedUntil"`
	More         bool           `json:"more"`
	Parcels      []InpostParcel `json:"parcels"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type AddressDetails struct {
	PostCode       string `json:"postCode"`
	City           string `json:"city"`
	Province       string `json:"province"`
	Street         string `json:"street"`
	BuildingNumber string `json:"buildingNumber"`
}
type PickUpPoint struct {
	Name                string          `json:"name"`
	Location            *Location       `json:"location"`
	LocationDescription string          `json:"locationDescription"`
	OpeningHours        string          `json:"openingHours"`
	AddressDetails      *AddressDetails `json:"addressDetails"`
	Virtual             int             `json:"virtual"`
	PointType           string          `json:"pointType"`
	Type                []string        `json:"type"`
	Location247         bool            `json:"location247"`
	Doubled             bool            `json:"doubled"`
	ImageURL            string          `json:"imageUrl"`
	EasyAccessZone      bool            `json:"easyAccessZone"`
	AirSensor           bool            `json:"airSensor"`
}

type CashOnDelivery struct {
	Paid    bool   `json:"paid,omitempty"`
	PayCode string `json:"payCode,omitempty"`
	Price   string `json:"price,omitempty"`
	Url     string `json:"url,omitempty"`
	// ddd
}

type Receiver struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
}
type Sender struct {
	Name string `json:"name"`
}

type Event struct {
	Type string    `json:"type"`
	Name string    `json:"name"`
	Date time.Time `json:"date"`
}

type Operations struct {
	ManualArchive         bool      `json:"manualArchive"`
	Delete                bool      `json:"delete"`
	Collect               bool      `json:"collect"`
	ExpandAvizo           bool      `json:"expandAvizo"`
	Highlight             bool      `json:"highlight"`
	RefreshUntil          time.Time `json:"refreshUntil"`
	RequestEasyAccessZone string    `json:"requestEasyAccessZone"`
	Voicebot              bool      `json:"voicebot"`
	CanShareToObserve     bool      `json:"canShareToObserve"`
	CanShareOpenCode      bool      `json:"canShareOpenCode"`
	CanShareParcel        bool      `json:"canShareParcel"`
}

type InpostParcel struct {
	ShipmentNumber         string       `json:"shipmentNumber"`
	ShipmentType           string       `json:"shipmentType"`
	OpenCode               string       `json:"openCode,omitempty"`
	QrCode                 string       `json:"qrCode,omitempty"`
	ExpiryDate             time.Time    `json:"expiryDate,omitempty"`
	StoredDate             time.Time    `json:"storedDate,omitempty"`
	ParcelSize             string       `json:"parcelSize"`
	Receiver               *Receiver    `json:"receiver"`
	Sender                 *Sender      `json:"sender"`
	PickUpPoint            *PickUpPoint `json:"pickUpPoint"`
	EndOfWeekCollection    bool         `json:"endOfWeekCollection"`
	Operations             Operations   `json:"operations"`
	Status                 string       `json:"status"`
	EventLog               []Event      `json:"eventLog"`
	AvizoTransactionStatus string       `json:"avizoTransactionStatus"`
	SharedTo               []any        `json:"sharedTo"`
	OwnershipStatus        string       `json:"ownershipStatus"`
	EconomyParcel          bool         `json:"economyParcel"`
}

type ValidateCompartmentResponse struct {
	SessionExpirationTime any    `json:"sessionExpirationTime,omitempty"`
	SessionUUID           string `json:"sessionUuid,omitempty"`
}
