package commondata

import "time"

type TrackingData struct {
	ShipmentNumber string
	ProviderName   string
	TrackingSteps  []*TrackingStep
}

type TrackingStep struct {
	Datetime time.Time `json:"datetime"`
	// CommonType denotes a type of this step that is well-known within the app
	CommonType string

	// Message is the message returned by the provider
	Message string

	// Location is the place where the step happened, can be empty dependign on provider
	Location string
}
