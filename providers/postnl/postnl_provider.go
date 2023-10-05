package postnl

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
	"github.com/alufers/paczkobot/httphelpers"

	"github.com/PuerkitoBio/goquery"
)

type PostnlProvider struct{}

func (pp *PostnlProvider) GetName() string {
	return "postnl"
}

func (pp *PostnlProvider) MatchesNumber(trackingNumber string) bool {
	return true
}

func (pp *PostnlProvider) Track(ctx context.Context, trackingNumber string) (*commondata.TrackingData, error) {
	//await fetch("https://postnl.post/details/", {
	//"credentials": "include",
	//"headers": {
	//"User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:86.0) Gecko/20100101 Firefox/86.0",
	//"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
	//"Accept-Language": "en-US,en;q=0.5",
	//"Content-Type": "application/x-www-form-urlencoded",
	//"Upgrade-Insecure-Requests": "1",
	//"Pragma": "no-cache",
	//"Cache-Control": "no-cache"
	//},
	//"referrer": "https://postnl.post/",
	//"body": "barcodes=LT504554145NL",
	//"method": "POST",
	//"mode": "cors"
	//});

	client := httphelpers.NewClientWithLogger()
	requestData := url.Values{}
	requestData.Set("barcodes", trackingNumber)

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://postnl.post/details/",
		strings.NewReader(requestData.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}
	commondata.SetCommonHTTPHeaders(&req.Header)
	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	httpResponse, err := client.Do(req)
	if err != nil {
		return nil, commonerrors.NewNetworkError(pp.GetName(), req)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP status code %v", httpResponse.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTML response from PostNL: %w", err)
	}
	datatables := doc.Find("table#datatables")
	if datatables.Length() <= 0 {
		return nil, commonerrors.NotFoundError
	}
	td := &commondata.TrackingData{
		ShipmentNumber: trackingNumber,
		ProviderName:   pp.GetName(),
		Destination:    strings.TrimSpace(datatables.Find("td.country:nth-child(3)").Text()),
		TrackingSteps:  []*commondata.TrackingStep{},
	}
	datatables.Find("tr.detail").Each(func(i int, row *goquery.Selection) {
		description := row.Find("td:nth-child(2)").Text()
		t, err := time.Parse("02-01-2006 15:04", strings.TrimSpace(row.Find("td.date").Text()))
		if err != nil {
			log.Printf("error while parsing date from postNL: %v", err)
		}
		td.TrackingSteps = append(td.TrackingSteps, &commondata.TrackingStep{
			Datetime:   t,
			CommonType: commondata.CommonTrackingStepType_UNKNOWN,
			Message:    description,
		})
	})

	return td, nil
}
