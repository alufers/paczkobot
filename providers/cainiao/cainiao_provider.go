package cainiao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/commonerrors"
)

type CainiaoProvider struct{}

func (pp *CainiaoProvider) GetName() string {
	return "cainiao"
}

func (pp *CainiaoProvider) MatchesNumber(trackingNumber string) bool {
	return true
}

func (pp *CainiaoProvider) Track(ctx context.Context, trackingNumber string) (*commondata.TrackingData, error) {

	req, err := http.NewRequest(
		"GET",
		"https://global.cainiao.com/global/detail.json?mailNos="+url.QueryEscape(trackingNumber)+"&lang=en-US&language=en-US",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}
	commondata.SetCommonHTTPHeaders(&req.Header)
	httpResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, commonerrors.NewNetworkError(pp.GetName(), req)
	}

	if httpResponse.StatusCode != 200 {
		return nil, commonerrors.NotFoundError
	}

	var cainiaoResponse CainiaoResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&cainiaoResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	if len(cainiaoResponse.Module) == 0 {
		return nil, commonerrors.NotFoundError
	}
	module := cainiaoResponse.Module[0]
	if len(module.DetailList) == 0 {
		return nil, commonerrors.NotFoundError
	}

	td := &commondata.TrackingData{
		ProviderName:   pp.GetName(),
		ShipmentNumber: trackingNumber,
		TrackingSteps:  make([]*commondata.TrackingStep, 0),
		Destination:    module.DestCountry,
		SentFrom:       module.OriginCountry,
	}

	for _, detail := range module.DetailList {

		date, _ := time.Parse("2006-01-02 15:04:05", detail.TimeStr)
		step := &commondata.TrackingStep{
			Datetime:   date,
			Message:    detail.Desc,
			Location:   detail.Group.NodeDesc,
			CommonType: commondata.CommonTrackingStepType_UNKNOWN,
		}
		td.TrackingSteps = append(td.TrackingSteps, step)
	}

	// attempt to augement the destination with the city
	// requires a separate request, if it fails, it's not a big deal
	cityReq, err := http.NewRequest(
		"GET",
		"https://global.cainiao.com/global/getCity.json?mailNo="+url.QueryEscape(trackingNumber)+"&lang=en-US&language=en-US",
		nil,
	)
	if err != nil {
		return td, nil
	}
	commondata.SetCommonHTTPHeaders(&cityReq.Header)
	cityResp, err := http.DefaultClient.Do(cityReq)
	if err != nil {
		return td, nil
	}
	if cityResp.StatusCode != 200 {
		return td, nil
	}
	var cityResponse GetCityResponse
	err = json.NewDecoder(cityResp.Body).Decode(&cityResponse)
	if err != nil || !cityResponse.Success {
		return td, nil
	}

	td.Destination = cityResponse.Module + ", " + td.Destination

	return td, nil
}
