package gls

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/commonerrors"
)

type GLSProvider struct {
}

func (p *GLSProvider) GetName() string {
	return "gls"
}

func (p *GLSProvider) MatchesNumber(trackingNumber string) bool {
	// ymmv
	return len(trackingNumber) == 8 || len(trackingNumber) == 11
}

func (p *GLSProvider) Track(ctx context.Context, trackingNumber string) (*commondata.TrackingData, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("https://gls-group.eu/app/service/open/rest/PL/en/rstt001?match=%s", trackingNumber),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, commonerrors.NewNetworkError(p.GetName(), req)
	}
	if res.StatusCode == http.StatusNotFound {
		return nil, commonerrors.NotFoundError
	}
	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		log.Printf("DHL ERROR BODY: %v", string(body))
		return nil, fmt.Errorf("HTTP status code %v", res.StatusCode)
	}

	decoder := json.NewDecoder(res.Body)
	decodedBody := &Response{}
	if err := decoder.Decode(decodedBody); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}
	if len(decodedBody.TuStatus) == 0 {
		return nil, commonerrors.NotFoundError
	}

	parcel := decodedBody.TuStatus[0]
	trackingData := &commondata.TrackingData{
		ProviderName:  p.GetName(),
		TrackingSteps: []*commondata.TrackingStep{},
	}

	for _, ref := range parcel.References {
		if ref.Type == "UNITNO" {
			trackingData.ShipmentNumber = ref.Value
		}
	}

	for _, owner := range parcel.Owners {
		if owner.Type == "DELIVERY" {
			trackingData.Destination = owner.Code
		}
	}

	for _, event := range parcel.History {
		datetime, _ := time.Parse("2006-01-02T15:04:05", event.Date+"T"+event.Time)
		location := event.Address.CountryName
		if event.Address.City != "" {
			location = event.Address.City + ", " + location
		}

		trackingData.TrackingSteps = append(trackingData.TrackingSteps, &commondata.TrackingStep{
			Datetime:   datetime,
			CommonType: commondata.CommonTrackingStepType_UNKNOWN,
			Message:    event.EvtDscr,
			Location:   location,
		})
	}

	return trackingData, nil
}
