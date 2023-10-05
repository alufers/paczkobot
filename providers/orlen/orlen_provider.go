package orlen

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/commonerrors"
	"github.com/alufers/paczkobot/httphelpers"
)

// https://www.fedex.com/pl-pl/online/domestic-tracking.html

var commonMappings = map[string]commondata.CommonTrackingStepType{
	"Twoja paczka jest w trakcie przygotowania przez nadawcę ":      commondata.CommonTrackingStepType_INFORMATION_PREPARED,
	"Paczka w trakcie przygotowania":                                commondata.CommonTrackingStepType_INFORMATION_PREPARED,
	"Z przyjemnością informujemy, że Twoja paczka została nadana. ": commondata.CommonTrackingStepType_SENT,
	"Paczka nadana":                 commondata.CommonTrackingStepType_SENT,
	"Paczka w sortowni regionalnej": commondata.CommonTrackingStepType_IN_TRANSIT,
}

type OrlenProvider struct{}

func (op *OrlenProvider) GetName() string {
	return "orlen"
}

func (op *OrlenProvider) MatchesNumber(trackingNumber string) bool {
	return len(trackingNumber) > 2
}

func (op *OrlenProvider) Track(ctx context.Context, trackingNumber string) (*commondata.TrackingData, error) {
	client := httphelpers.NewClientWithLogger()
	trackingReq, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://nadaj.orlenpaczka.pl/parcel/api-status?id="+
			url.QueryEscape(trackingNumber)+
			"&jsonp=callback&_="+fmt.Sprintf("%d", time.Now().Unix()),
		nil,
	)
	if err != nil {
		return nil, err
	}
	commondata.SetCommonHTTPHeaders(&trackingReq.Header)

	resp, err := client.Do(trackingReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotFound {
		return nil, commonerrors.NotFoundError
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("orlen: unexpected status code: %d", resp.StatusCode)
	}

	allData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// strip the callback() from the response
	reg := regexp.MustCompile(`^callback[(]|[)];$`)
	strData := reg.ReplaceAllString(string(allData), "")

	respData := &OrlenResponse{}
	err = json.Unmarshal([]byte(strData), respData)
	if err != nil {
		return nil, err
	}

	if respData.Status != "OK" {
		return nil, commonerrors.NotFoundError
	}

	td := &commondata.TrackingData{
		ShipmentNumber: trackingNumber,
		ProviderName:   op.GetName(),
		TrackingSteps:  []*commondata.TrackingStep{},
	}

	for _, ev := range respData.History {
		// example date: 20-06-2023, 21:23

		datetime, _ := time.Parse("02-01-2006, 15:04", ev.Date)
		commonStepType, ok := commonMappings[ev.Label]
		if !ok {
			commonStepType, ok = commonMappings[ev.LabelShort]
			if !ok {
				commonStepType = commondata.CommonTrackingStepType_UNKNOWN
			}
		}
		td.TrackingSteps = append(td.TrackingSteps, &commondata.TrackingStep{
			Datetime:   datetime,
			Message:    ev.Label,
			CommonType: commonStepType,
		})
	}

	if len(td.TrackingSteps) == 0 {
		return nil, commonerrors.NotFoundError
	}

	return td, nil
}
