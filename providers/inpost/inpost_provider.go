package inpost

import (
	"fmt"
	"github.com/alufers/paczkobot/commondata"
	"regexp"

	"github.com/alufers/paczkobot/providers/inpost/inposttrackingapi"
)

type InpostProvider struct {
}

var inpostNumberRegex = regexp.MustCompile(`^\d{24}$`)

func (ip *InpostProvider) GetName() string {
	return "inpost"
}

func (ip *InpostProvider) MatchesNumber(trackingNumber string) bool {
	return inpostNumberRegex.MatchString(trackingNumber)
}

func (ip *InpostProvider) Track(trackingNumber string) (*commondata.TrackingData, error) {
	data, err := inposttrackingapi.GetTrackingData(trackingNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from inpost API: %w", err)
	}
	td := &commondata.TrackingData{
		ShipmentNumber: trackingNumber,
		ProviderName:   "inpost",
		TrackingSteps:  []*commondata.TrackingStep{},
	}
	for i := len(data.TrackingDetails) - 1; i >= 0; i-- {
		d := data.TrackingDetails[i]
		td.TrackingSteps = append(td.TrackingSteps, &commondata.TrackingStep{
			Datetime:   d.Datetime,
			CommonType: d.Status,
			Message:    d.Status,
		})
	}
	return td, nil
}
