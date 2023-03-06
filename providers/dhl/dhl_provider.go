package dhl

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
	"github.com/spf13/viper"
)

var statusCodeMappings = map[string]commondata.CommonTrackingStepType{
	"delivered": commondata.CommonTrackingStepType_DELIVERED,
	"transit":   commondata.CommonTrackingStepType_IN_TRANSIT,
	"failure":   commondata.CommonTrackingStepType_FAILURE,
}

var descriptionMappings = map[string]commondata.CommonTrackingStepType{
	"przesyłka oczekuje na odbiór w DHL POP":       commondata.CommonTrackingStepType_READY_FOR_PICKUP,
	"przesyłka w drodze do Punktu DHL POP":         commondata.CommonTrackingStepType_OUT_FOR_DELIVERY,
	"przesyłka przyjęta w terminalu nadawczym DHL": commondata.CommonTrackingStepType_SENT,
	"DHL otrzymał dane elektroniczne przesyłki.":   commondata.CommonTrackingStepType_INFORMATION_PREPARED,
}

type DHLProvider struct{}

func (pp *DHLProvider) GetName() string {
	return "dhl"
}

func (pp *DHLProvider) MatchesNumber(trackingNumber string) bool {
	return viper.GetBool("tracking.providers.dhl.enable")
}

func (pp *DHLProvider) Track(ctx context.Context, trackingNumber string) (*commondata.TrackingData, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("https://api-eu.dhl.com/track/shipments?trackingNumber=%v", trackingNumber),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request to tracking page: %w", err)
	}
	commondata.SetCommonHTTPHeaders(&req.Header)
	req.Header.Set("DHL-API-Key", viper.GetString("tracking.providers.dhl.api_key"))
	httpResponse, err := client.Do(req)
	if err != nil {
		return nil, commonerrors.NewNetworkError(pp.GetName(), req)
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusNotFound {
		return nil, commonerrors.NotFoundError
	}
	if httpResponse.StatusCode != 200 {
		body, _ := io.ReadAll(httpResponse.Body)
		log.Printf("DHL ERROR BODY: %v", string(body))
		return nil, fmt.Errorf("HTTP status code %v", httpResponse.StatusCode)
	}

	decoder := json.NewDecoder(httpResponse.Body)
	decodedBody := &DHLResponse{}
	if err := decoder.Decode(decodedBody); err != nil {
		return nil, fmt.Errorf("failed to parse tracking response JSON: %w", err)
	}
	if len(decodedBody.Shipments) <= 0 {
		return nil, commonerrors.NotFoundError
	}
	trackingData := &commondata.TrackingData{
		ShipmentNumber: trackingNumber,
		ProviderName:   pp.GetName(),
		TrackingSteps:  []*commondata.TrackingStep{},
	}

	if decodedBody.Shipments[0].Destination != nil {
		trackingData.Destination = decodedBody.Shipments[0].Destination.String()
	}

	for _, ev := range decodedBody.Shipments[0].Events {
		loc := ""
		if ev.Location != nil {
			loc = ev.Location.String()
		}
		datetime, _ := time.Parse("2006-01-02T15:04:05", ev.Timestamp)

		var commonType commondata.CommonTrackingStepType
		for k, v := range descriptionMappings {
			if strings.Contains(ev.Description, k) {
				commonType = v
				break
			}
		}
		if commonType == commondata.CommonTrackingStepType_UNKNOWN {
			for k, v := range statusCodeMappings {
				if ev.StatusCode == k {
					commonType = v
					break
				}
			}
		}

		trackingData.TrackingSteps = append(trackingData.TrackingSteps, &commondata.TrackingStep{
			Datetime:   datetime,
			Location:   loc,
			Message:    ev.Description,
			CommonType: commonType,
		})
	}

	return trackingData, nil
}
