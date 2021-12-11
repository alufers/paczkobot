package dpdcompl

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/commonerrors"

	"github.com/PuerkitoBio/goquery"
)

var descriptionMappings = map[string]commondata.CommonTrackingStepType{

	"Przyjęcie przesyłki w oddziale DPD": commondata.CommonTrackingStepType_IN_TRANSIT,
	"Przesyłka odebrana przez Kuriera":   commondata.CommonTrackingStepType_IN_TRANSIT,
	"Przekazanie przesyłki kurierowi":    commondata.CommonTrackingStepType_IN_TRANSIT,

	"Przesyłka doręczona":                    commondata.CommonTrackingStepType_DELIVERED,
	"przesyłka oczekuje na odbiór w DHL POP": commondata.CommonTrackingStepType_READY_FOR_PICKUP,
	"Wydanie przesyłki do doręczenia":        commondata.CommonTrackingStepType_OUT_FOR_DELIVERY,
	"Nadanie przesyłki w punkcie Pickup":     commondata.CommonTrackingStepType_SENT,
	"Zarejestrowano dane przesyłki":          commondata.CommonTrackingStepType_INFORMATION_PREPARED,
}

type DpdComPlProvider struct {
}

func (dp *DpdComPlProvider) GetName() string {
	return "dpd-com-pl"
}

func (dp *DpdComPlProvider) MatchesNumber(trackingNumber string) bool {
	return true
}

func (dp *DpdComPlProvider) Track(ctx context.Context, trackingNumber string) (*commondata.TrackingData, error) {

	requestData := url.Values{}
	requestData.Set("q", trackingNumber)
	requestData.Set("typ", "1")
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://tracktrace.dpd.com.pl/findPackage",
		strings.NewReader(requestData.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}
	commondata.SetCommonHTTPHeaders(&req.Header)
	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	httpResponse, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, commonerrors.NewNetworkError(dp.GetName(), req)
	}

	if httpResponse.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP status code %v", httpResponse.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTML response from DPD: %w", err)
	}
	datatables := doc.Find("table.table-track")
	if datatables.Length() <= 0 {
		return nil, commonerrors.NotFoundError
	}
	td := &commondata.TrackingData{
		ShipmentNumber: trackingNumber,
		ProviderName:   dp.GetName(),
		TrackingSteps:  []*commondata.TrackingStep{},
	}
	datatables.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		date := row.Find("td:nth-child(1)").Text()
		timeText := row.Find("td:nth-child(2)").Text()
		description := row.Find("td:nth-child(3)").Text()
		location := row.Find("td:nth-child(4)").Text()

		t, err := time.Parse("2006-01-02 15:04:05", strings.TrimSpace(date)+" "+strings.TrimSpace(timeText))
		if err != nil {
			log.Printf("error while parsing date from DPD: %v", err)
		}
		var commonType commondata.CommonTrackingStepType
		for k, v := range descriptionMappings {
			if strings.Contains(description, k) {
				commonType = v
				break
			}
		}
		td.TrackingSteps = append(td.TrackingSteps, &commondata.TrackingStep{
			Datetime:   t,
			CommonType: commonType,
			Message:    strings.TrimSpace(description),
			Location:   location,
		})
	})

	return td, nil
}
