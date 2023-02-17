package commondata

import "time"

type TrackingData struct {
	ShipmentNumber string          `json:"shipmentNumber"`
	ProviderName   string          `json:"providerName"`
	Destination    string          `json:"destination"`
	SentFrom       string          `json:"sentFrom"`
	Weight         float64         `json:"weight"` // in kg
	TrackingSteps  []*TrackingStep `json:"trackingSteps"`
}

func (td *TrackingData) ApplyCommonTypeMappings(mappings map[string]CommonTrackingStepType) {
	for _, step := range td.TrackingSteps {
		if step.CommonType == CommonTrackingStepType_UNKNOWN {
			if newType, ok := mappings[step.Message]; ok {
				step.CommonType = newType
			}
		}
	}
}

type TrackingStep struct {
	Datetime time.Time `json:"datetime"`
	// CommonType denotes a type of this step that is well-known within the app
	CommonType CommonTrackingStepType `json:"commonType"`

	// Message is the message returned by the provider
	Message string `json:"message"`

	// Location is the place where the step happened, can be empty dependign on provider
	Location string `json:"location"`
}
