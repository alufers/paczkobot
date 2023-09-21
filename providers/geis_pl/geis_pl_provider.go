package geis_pl

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/commonerrors"
	providerutil "github.com/alufers/paczkobot/provider_util"

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
	doc, err := providerutil.FetchGoqueryDocument(
		ctx,
		gp.GetName(),
		"https://www.geis.pl/en/detail-of-cargo?packNumber="+url.QueryEscape(trackingNumber),
		true,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTML response from DPD: %w", err)
	}
	datatables := doc.Find(".trace211 table").First()
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
