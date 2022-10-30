package inpostextra

import "time"

type MultiCompartment struct {
	Collected       bool     `json:"collected,omitempty"`
	Presentation    bool     `json:"presentation,omitempty"`
	ShipmentNumbers []string `json:"shipmentNumbers,omitempty"`
	Uuid            string   `json:"uuid,omitempty"`
}

type AddressDetails struct {
	BuildingNumber string `json:"buildingNumber,omitempty"`
	City           string `json:"city,omitempty"`
	FlatNumber     string `json:"flatNumber,omitempty"`
	PostCode       string `json:"postCode,omitempty"`
	Province       string `json:"province,omitempty"`
	Street         string `json:"street,omitempty"`
}

type PointLocation struct {
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

type DeliveryPointData struct {
	AddressDetails      *AddressDetails   `json:"addressDetails,omitempty"`
	Location            *PointLocation    `json:"location,omitempty"`
	Location247         bool              `json:"location247,omitempty"`
	LocationDescription string            `json:"locationDescription,omitempty"`
	Name                string            `json:"name,omitempty"`
	OpeningHours        string            `json:"openingHours,omitempty"`
	PaymentType         map[string]string `json:"paymentType,omitempty"`
	Status              string            `json:"status,omitempty"`
	Type_               []string          `json:"type,omitempty"`
	Virtual             int32             `json:"virtual,omitempty"`
}

type CashOnDelivery struct {
	Paid    bool   `json:"paid,omitempty"`
	PayCode string `json:"payCode,omitempty"`
	Price   string `json:"price,omitempty"`
	Url     string `json:"url,omitempty"`
	// ddd
}

type ParcelHistory struct {
	Date   time.Time `json:"date,omitempty"`
	Status string    `json:"status,omitempty"`
}

type InpostParcel struct {
	CashOnDelivery          *CashOnDelivery    `json:"cashOnDelivery,omitempty"`
	EndOfWeekCollection     bool               `json:"endOfWeekCollection,omitempty"`
	ExpiryDate              time.Time          `json:"expiryDate,omitempty"`
	IsMobileCollectPossible bool               `json:"isMobileCollectPossible,omitempty"`
	IsObserved              bool               `json:"isObserved,omitempty"`
	MultiCompartment        *MultiCompartment  `json:"multiCompartment,omitempty"`
	OpenCode                string             `json:"openCode,omitempty"`
	PhoneNumber             string             `json:"phoneNumber,omitempty"`
	PickupDate              time.Time          `json:"pickupDate,omitempty"`
	PickupPoint             *DeliveryPointData `json:"pickupPoint,omitempty"`
	QrCode                  string             `json:"qrCode,omitempty"`
	ReturnedToSenderDate    time.Time          `json:"returnedToSenderDate,omitempty"`
	SenderName              string             `json:"senderName,omitempty"`
	ShipmentNumber          string             `json:"shipmentNumber,omitempty"`
	ShipmentType            string             `json:"shipmentType,omitempty"`
	Status                  string             `json:"status,omitempty"`
	StatusHistory           []ParcelHistory    `json:"statusHistory,omitempty"`
	StoredDate              time.Time          `json:"storedDate,omitempty"`
}

type ValidateCompartmentResponse struct {
	SessionExpirationTime time.Time `json:"sessionExpirationTime,omitempty"`
	SessionUUID           string    `json:"sessionUuid,omitempty"`
}
