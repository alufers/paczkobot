package inposttrackingapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/alufers/paczkobot/commonerrors"
	"github.com/alufers/paczkobot/httphelpers"
)

func GetTrackingData(ctx context.Context, parcelNumber string) (*TrackingAPISchema, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://api-shipx-pl.easypack24.net/v1/tracking/%s", url.QueryEscape(parcelNumber)), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request to get tracking data: %w", err)
	}

	client := httphelpers.NewClientWithLogger()
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request to get tracking data (url: %v): %w", request.URL.String(), err)
	}
	defer resp.Body.Close()
	if resp.StatusCode > 399 {
		data, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == 404 {
			return nil, commonerrors.NotFoundError
		}
		return nil, fmt.Errorf("received HTTP error code from tracking endpoint %v, HTTP body: %v", resp.StatusCode, string(data))
	}
	decoder := json.NewDecoder(resp.Body)
	trackingData := &TrackingAPISchema{}
	err = decoder.Decode(trackingData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}
	return trackingData, nil
}
