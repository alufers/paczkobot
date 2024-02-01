package usps

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/commonerrors"
	providerutil "github.com/alufers/paczkobot/provider_util"

	"github.com/PuerkitoBio/goquery"
)

var descriptionMappings = map[string]commondata.CommonTrackingStepType{
	"Shipping Label Created, USPS Awaiting Item": commondata.CommonTrackingStepType_INFORMATION_PREPARED,
	"Accepted at USPS Origin Facility":           commondata.CommonTrackingStepType_SENT,
	"In Transit to Next Facility":                commondata.CommonTrackingStepType_IN_TRANSIT,
	"Departed USPS Regional Facility":            commondata.CommonTrackingStepType_IN_TRANSIT,
	"Arrived at USPS Regional Facility":          commondata.CommonTrackingStepType_IN_TRANSIT,
	"Processing at USPS Facility":                commondata.CommonTrackingStepType_IN_TRANSIT,
	"Arrived at Post Office":                     commondata.CommonTrackingStepType_OUT_FOR_DELIVERY,
}

type USPSProvider struct{}

func (gp *USPSProvider) GetName() string {
	return "usps"
}

var uspsRegexes = []*regexp.Regexp{
	regexp.MustCompile(`^(\d{22})$`),
	regexp.MustCompile(`^(\d{10})$`),
	regexp.MustCompile(`^(EC\d{9}US)$`),
	regexp.MustCompile(`^(CP\d{9}US)$`),
	regexp.MustCompile(`^(\d{22}EA\d{9}US)$`),
}

func (gp *USPSProvider) MatchesNumber(trackingNumber string) bool {
	// remove all spaces
	trackingNumber = strings.ReplaceAll(trackingNumber, " ", "")
	for _, regex := range uspsRegexes {
		if regex.MatchString(trackingNumber) {
			return true
		}
	}
	return false
}

func (gp *USPSProvider) Track(ctx context.Context, trackingNumber string) (*commondata.TrackingData, error) {
	doc, err := providerutil.FetchGoqueryDocument(
		ctx,
		gp.GetName(),
		"https://tools.usps.com/go/TrackConfirmAction.action?tLabels="+url.QueryEscape(trackingNumber),
		true,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTML response from Deutsche Post: %w", err)
	}
	datatables := doc.Find(".tracking-progress-bar-status-container").First()
	if datatables.Length() <= 0 {
		return nil, commonerrors.NotFoundError
	}
	td := &commondata.TrackingData{
		ShipmentNumber: trackingNumber,
		ProviderName:   gp.GetName(),
		TrackingSteps:  []*commondata.TrackingStep{},
	}
	datatables.Find(".tb-step").Each(func(i int, row *goquery.Selection) {
		date := strings.ReplaceAll(row.Find(".tb.date").Text(), " ", "")

		description := row.Find(".tb-status-detail").Text()
		location := row.Find(".tb-location").Text()

		// February 1, 2024, 12:26 am
		t, err := time.Parse("January 2, 2006, 3:04 pm", date)
		if err != nil {
			log.Printf("error while parsing date from USPS Post: %v", err)
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
