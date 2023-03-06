package geis_pl

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
	"Delivered":                             commondata.CommonTrackingStepType_DELIVERED,
	"Acceptance at depot":                   commondata.CommonTrackingStepType_IN_TRANSIT,
	"Dostarczenie":                          commondata.CommonTrackingStepType_DELIVERED,
	"Załadunek na rozwóz":                   commondata.CommonTrackingStepType_OUT_FOR_DELIVERY,
	"Loading for distribution to consignee": commondata.CommonTrackingStepType_OUT_FOR_DELIVERY,
}

type GeisPlProvider struct{}

func (gp *GeisPlProvider) GetName() string {
	return "geis-pl"
}

func (gp *GeisPlProvider) MatchesNumber(trackingNumber string) bool {
	return true
}

func (gp *GeisPlProvider) Track(ctx context.Context, trackingNumber string) (*commondata.TrackingData, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://www.geis.pl/en/detail-of-cargo?packNumber="+url.QueryEscape(trackingNumber),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}
	commondata.SetCommonHTTPHeaders(&req.Header)

	httpResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, commonerrors.NewNetworkError(gp.GetName(), req)
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP status code %v", httpResponse.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTML response from DPD: %w", err)
	}
	datatables := doc.Find("table.table-tracking").First()
	if datatables.Length() <= 0 {
		return nil, commonerrors.NotFoundError
	}
	td := &commondata.TrackingData{
		ShipmentNumber: trackingNumber,
		ProviderName:   gp.GetName(),
		TrackingSteps:  []*commondata.TrackingStep{},
	}
	datatables.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		date := strings.ReplaceAll(row.Find("td:nth-child(1)").Text(), " ", "")
		timeHourMinute := strings.ReplaceAll(row.Find("td:nth-child(2)").Text(), " ", "")
		date = date + " " + timeHourMinute
		log.Printf("cols: %v, %v", row.Find("td:nth-child(1)").Text(), row.Find("td:nth-child(1)").Text())
		log.Printf("date: %v", date)
		location := row.Find("td:nth-child(4)").Text()
		description := row.Find("td:nth-child(5)").Text()

		// 02.06.2022 10:39
		t, err := time.Parse("02.01.2006 15:04", date)
		if err != nil {
			log.Printf("error while parsing date from geis PL: %v", err)
		}

		td.TrackingSteps = append(td.TrackingSteps, &commondata.TrackingStep{
			Datetime: t,

			Location: strings.TrimSpace(location),
			Message:  strings.TrimSpace(description),
		})
	})

	td.ApplyCommonTypeMappings(descriptionMappings)

	return td, nil
}
