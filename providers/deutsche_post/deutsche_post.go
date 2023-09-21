package deutsche_post

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
	"Shipment information uploaded to Deutsche Post": commondata.CommonTrackingStepType_INFORMATION_PREPARED,
	"Item received at Deutsche Post Mailterminal":    commondata.CommonTrackingStepType_SENT,
	"Departure to destination country":               commondata.CommonTrackingStepType_IN_TRANSIT,
	"Delivered":                                      commondata.CommonTrackingStepType_DELIVERED,
}

type DeutschePostProvider struct{}

func (gp *DeutschePostProvider) GetName() string {
	return "deutsche-post"
}

func (gp *DeutschePostProvider) MatchesNumber(trackingNumber string) bool {
	return true
}

func (gp *DeutschePostProvider) Track(ctx context.Context, trackingNumber string) (*commondata.TrackingData, error) {
	doc, err := providerutil.FetchGoqueryDocument(
		ctx,
		gp.GetName(),
		"https://www.packet.deutschepost.com/webapp/public/packet_traceit.xhtml?barcode="+url.QueryEscape(trackingNumber),
		true,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTML response from Deutsche Post: %w", err)
	}
	datatables := doc.Find(".gmpacketTraceItHistoryTable table").First()
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

		description := row.Find("td:nth-child(2)").Text()

		// 02.06.2022 10:39
		t, err := time.Parse("02.01.2006", date)
		if err != nil {
			log.Printf("error while parsing date from Detsche Post: %v", err)
		}

		td.TrackingSteps = append(td.TrackingSteps, &commondata.TrackingStep{
			Datetime: t,

			Message: strings.TrimSpace(description),
		})
	})

	td.Destination = doc.Find(".gmpacketTraceItDestinationTo .section>.section").First().Text()
	td.SentFrom = doc.Find(".gmpacketTraceItDestinationFrom .section>.section").First().Text()

	// replace newlines with spaces
	// and replace multiple spaces with single space
	re1 := regexp.MustCompile(`[\n\r]+`)
	re2 := regexp.MustCompile(`\s+`)
	td.Destination = re1.ReplaceAllString(td.Destination, " ")
	td.Destination = re2.ReplaceAllString(td.Destination, " ")
	td.SentFrom = re1.ReplaceAllString(td.SentFrom, " ")
	td.SentFrom = re2.ReplaceAllString(td.SentFrom, " ")

	td.ApplyCommonTypeMappings(descriptionMappings)

	return td, nil
}
