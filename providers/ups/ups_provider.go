package ups

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/commonerrors"
	"github.com/alufers/paczkobot/httphelpers"
)

type UPSProvider struct{}

func (pp *UPSProvider) GetName() string {
	return "ups"
}

func (pp *UPSProvider) MatchesNumber(trackingNumber string) bool {
	return true
}

func (pp *UPSProvider) Track(ctx context.Context, trackingNumber string) (*commondata.TrackingData, error) {
	client := httphelpers.NewClientWithLogger()
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}
	client.Jar = jar
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://www.ups.com/track?loc=en_US&requester=ST/trackdetails",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request to tracking page: %w", err)
	}
	commondata.SetCommonHTTPHeaders(&req.Header)
	httpResponse, err := client.Do(req)
	if err != nil {
		return nil, commonerrors.NewNetworkError(pp.GetName(), req)
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP home page status code %v", httpResponse.StatusCode)
	}
	u, _ := url.Parse("https://www.ups.com/")
	xsrfToken := ""
	for _, c := range jar.Cookies(u) {
		if c.Name == "X-XSRF-TOKEN-ST" {
			xsrfToken = c.Value
		}
	}
	reqBody, err := json.Marshal(map[string]interface{}{
		"Locale":         "en_US",
		"Requester:":     "wems_1z/trackdetails",
		"TrackingNumber": []interface{}{trackingNumber},
		"returnToValue":  "",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	req, err = http.NewRequestWithContext(
		ctx,
		"POST",
		"https://www.ups.com/track/api/Track/GetStatus?loc=en_US",
		bytes.NewBuffer(reqBody),
	)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("X-XSRF-Token", xsrfToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request to tracking page: %w", err)
	}
	commondata.SetCommonHTTPHeaders(&req.Header)
	httpResponse, err = client.Do(req)

	if err != nil {
		return nil, commonerrors.NewNetworkError(pp.GetName(), req)
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP status code %v", httpResponse.StatusCode)
	}

	decoder := json.NewDecoder(httpResponse.Body)
	decodedBody := &UPSJsonSchema{}
	if err := decoder.Decode(decodedBody); err != nil {
		return nil, fmt.Errorf("failed to parse tracking response JSON: %w", err)
	}
	if decodedBody.TrackDetails == nil || len(decodedBody.TrackDetails) == 0 {
		return nil, commonerrors.NotFoundError
	}
	details := decodedBody.TrackDetails[0]

	if details.ErrorCode == "504" || details.ErrorCode == "Tracking number not found in database" {
		return nil, commonerrors.NotFoundError
	}
	destinationParts := []string{}
	destinationPartsAll := []string{
		details.ShipToAddress.AttentionName,
		details.ShipToAddress.CompanyName,
		details.ShipToAddress.StreetAddress1,
		details.ShipToAddress.StreetAddress2,
		details.ShipToAddress.StreetAddress3,
		details.ShipToAddress.ZipCode,
		details.ShipToAddress.City,
		details.ShipToAddress.Province,
		details.ShipToAddress.Country,
	}

	for _, p := range destinationPartsAll {
		if strings.TrimSpace(p) != "" {
			destinationParts = append(destinationParts, p)
		}
	}
	steps := []*commondata.TrackingStep{}
	for _, d := range details.ShipmentProgressActivities {
		actScan, _ := d.ActivityScan.(string)
		t, err := time.Parse("01/02/2006", strings.TrimSpace(d.Date))
		if err != nil {
			log.Printf("error while parsing date from ups: %v", err)
			if actScan == "" {
				continue
			}
		}
		ref, _ := time.Parse("03:04 PM", "12:00 AM")
		t2, err := time.Parse("3:04 PM", strings.TrimSpace(strings.ReplaceAll(d.Time, ".", "")))
		if err != nil {
			log.Printf("failed to parse time from ups: %v", err)
		}
		t = t.Add(t2.Sub(ref))
		steps = append(steps, &commondata.TrackingStep{
			Datetime:   t,
			Message:    strings.TrimSpace(actScan),
			CommonType: commondata.CommonTrackingStepType_UNKNOWN,
			Location:   d.Location,
		})
	}
	trackingData := &commondata.TrackingData{
		ShipmentNumber: trackingNumber,
		ProviderName:   pp.GetName(),
		Destination:    strings.Join(destinationParts, ", "),
		TrackingSteps:  steps,
	}

	return trackingData, nil
}
