package inpost

import (
	"context"
	"fmt"
	"regexp"

	"github.com/alufers/paczkobot/commondata"

	"github.com/alufers/paczkobot/providers/inpost/inposttrackingapi"
)

var statusMappings = map[string]commondata.CommonTrackingStepType{
	"confirmed":             commondata.CommonTrackingStepType_INFORMATION_PREPARED,
	"collected_from_sender": commondata.CommonTrackingStepType_SENT,
	"dispatched_by_sender":  commondata.CommonTrackingStepType_SENT,

	"taken_by_courier":         commondata.CommonTrackingStepType_IN_TRANSIT,
	"adopted_at_source_branch": commondata.CommonTrackingStepType_IN_TRANSIT,
	"sent_from_source_branch":  commondata.CommonTrackingStepType_IN_TRANSIT,

	"out_for_delivery":     commondata.CommonTrackingStepType_OUT_FOR_DELIVERY,
	"ready_to_pickup":      commondata.CommonTrackingStepType_READY_FOR_PICKUP,
	"stack_in_box_machine": commondata.CommonTrackingStepType_READY_FOR_PICKUP,
	"delivered":            commondata.CommonTrackingStepType_DELIVERED,
}

type InpostProvider struct {
}

var inpostNumberRegex = regexp.MustCompile(`^\d{24}$`)

func (ip *InpostProvider) GetName() string {
	return "inpost"
}

func (ip *InpostProvider) MatchesNumber(trackingNumber string) bool {
	return inpostNumberRegex.MatchString(trackingNumber)
}

func (ip *InpostProvider) Track(ctx context.Context, trackingNumber string) (*commondata.TrackingData, error) {
	data, err := inposttrackingapi.GetTrackingData(ctx, trackingNumber)
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
		var commonStep commondata.CommonTrackingStepType
		if s, ok := statusMappings[d.Status]; ok {
			commonStep = s
		} else {
			commonStep = commondata.CommonTrackingStepType_UNKNOWN
		}
		td.TrackingSteps = append(td.TrackingSteps, &commondata.TrackingStep{
			Datetime:   d.Datetime,
			CommonType: commonStep,
			Message:    d.Status,
		})
	}
	return td, nil
}
