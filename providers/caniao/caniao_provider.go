package caniao

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/commonerrors"
)

type CaniaoProvider struct {
}

func (pp *CaniaoProvider) GetName() string {
	return "caniao"
}

func (pp *CaniaoProvider) MatchesNumber(trackingNumber string) bool {
	return true
}

func (pp *CaniaoProvider) Track(ctx context.Context, trackingNumber string) (*commondata.TrackingData, error) {

	requestData := url.Values{}
	requestData.Set("barcodes", trackingNumber)

	req, err := http.NewRequest(
		"GET",
		"https://global.cainiao.com/detail.htm?mailNoList="+url.QueryEscape(trackingNumber),
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
		return nil, fmt.Errorf("HTTP status code %v", httpResponse.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTML response from caniao: %w", err)
	}
	dataTextarea := doc.Find("textarea#waybill_list_val_box")
	if strings.TrimSpace(dataTextarea.Text()) == "" {
		return nil, fmt.Errorf("textarea #waybill_list_val_box not found in response")
	}
	var trackingData CaniaoJSONRoot
	if err := json.Unmarshal([]byte(dataTextarea.Text()), &trackingData); err != nil {
		log.Printf("caniao malformed JSON: %v", dataTextarea.Text())
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	if len(trackingData.Data) <= 0 {
		return nil, fmt.Errorf("len(trackingData.Data) <= 0")
	}
	if trackingData.Data[0].ErrorCode == "RESULT_EMPTY" || trackingData.Data[0].ErrorCode == "ORDER_NOT_FOUND" {
		return nil, commonerrors.NotFoundError
	}

	if !trackingData.Data[0].Success {
		return nil, fmt.Errorf("caniao error: trackingData.Data[0].ErrorCode = %v", trackingData.Data[0].ErrorCode)
	}
	td := &commondata.TrackingData{
		ShipmentNumber: trackingNumber,
		ProviderName:   pp.GetName(),
		Destination:    trackingData.Data[0].DestCountry,
		TrackingSteps:  []*commondata.TrackingStep{},
	}
	for _, d := range trackingData.Data[0].Section2.DetailList {
		t, _ := time.Parse("2006-01-02 15:04:05", d.Time)
		status := d.Status
		if status == "" {
			status = d.Desc
		}
		td.TrackingSteps = append(td.TrackingSteps, &commondata.TrackingStep{
			Datetime:   t,
			CommonType: commondata.CommonTrackingStepType_UNKNOWN,
			Message:    d.Desc,
		})
	}

	return td, nil
}
