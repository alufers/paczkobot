package dpd_eu

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/commonerrors"
)

type DPDEuProvider struct{}

func (p *DPDEuProvider) GetName() string {
	return "dpd-eu"
}

func (p *DPDEuProvider) MatchesNumber(trackingNumber string) bool {
	return true
}

func (p *DPDEuProvider) Track(ctx context.Context, trackingNumber string) (*commondata.TrackingData, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("https://tracking.dpd.de/rest/plc/en_US/%s", trackingNumber),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, commonerrors.NewNetworkError(p.GetName(), req)
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return nil, commonerrors.NotFoundError
	}
	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		log.Printf("DPD ERROR BODY: %v", string(body))
		return nil, fmt.Errorf("HTTP status code %v", res.StatusCode)
	}

	decoder := json.NewDecoder(res.Body)
	decodedBody := &DpdEuResponse{}
	if err := decoder.Decode(decodedBody); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}
	if decodedBody.ParcellifecycleResponse == nil {
		return nil, commonerrors.NotFoundError
	}

	trackingData := &commondata.TrackingData{
		ProviderName:  p.GetName(),
		TrackingSteps: []*commondata.TrackingStep{},
	}

	for _, event := range decodedBody.ParcellifecycleResponse.ParcelLifeCycleData.ScanInfo.Scan {
		datetime, _ := time.Parse("2006-01-02T15:04:05", event.Date)
		location := event.ScanData.Location
		dscr := strings.Join(event.ScanDescription.Content, " ")

		trackingData.TrackingSteps = append(trackingData.TrackingSteps, &commondata.TrackingStep{
			Datetime:   datetime,
			CommonType: commondata.CommonTrackingStepType_UNKNOWN,
			Message:    dscr,
			Location:   location,
		})
	}

	return trackingData, nil
}
