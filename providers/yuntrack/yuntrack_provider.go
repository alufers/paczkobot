package yuntrack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/commonerrors"
)

type YunTrack struct {
}

func (p *YunTrack) GetName() string {
	return "yuntrack"
}

func (p *YunTrack) MatchesNumber(trackingNumber string) bool {
	// ymmv
	return true
}

func (p *YunTrack) Track(ctx context.Context, trackingNumber string) (*commondata.TrackingData, error) {

	bodyData := &YunTrackRequest{
		CaptchaVerification: "",
		Year:                0,
		NumberList:          []string{trackingNumber},
	}
	bodyBytes, err := json.Marshal(bodyData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body data: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("https://services.yuntrack.com/Track/Query"),
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	/*
		Accept-Encoding: gzip, deflate, br
		Accept-Language: en-US,en;q=0.9,de-DE;q=0.8,de;q=0.7,pl-PL;q=0.6,pl;q=0.5
		Authorization: Nebula token:undefined
		Cache-Control: no-cache
		Connection: keep-alive
		Content-Length: 71
		Content-Type: application/json
		Host: services.yuntrack.com
		Origin: https://www.yuntrack.com
		Pragma: no-cache
		Referer: https://www.yuntrack.com/
		Sec-Fetch-Dest: empty
		Sec-Fetch-Mode: cors
		Sec-Fetch-Site: same-site
		Sec-GPC: 1
		User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36
	*/
	req.Header.Set("Accept", "application/json, text/plain, */*")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,de-DE;q=0.8,de;q=0.7,pl-PL;q=0.6,pl;q=0.5")
	req.Header.Set("Authorization", "Nebula token:undefined")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", "services.yuntrack.com")
	req.Header.Set("Origin", "https://www.yuntrack.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "https://www.yuntrack.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-GPC", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, commonerrors.NewNetworkError(p.GetName(), req)
	}
	if res.StatusCode == http.StatusNotFound {
		return nil, commonerrors.NotFoundError
	}
	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		log.Printf("YUNTRACK : %v", string(body))
		return nil, fmt.Errorf("HTTP status code %v", res.StatusCode)
	}

	decoder := json.NewDecoder(res.Body)
	decodedBody := &YunTrackResponse{}
	if err := decoder.Decode(decodedBody); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}
	if len(decodedBody.ResultList) == 0 || decodedBody.ResultList[0].Status == 0 {
		return nil, commonerrors.NotFoundError
	}
	firstResult := decodedBody.ResultList[0]

	trackingData := &commondata.TrackingData{
		ProviderName:   p.GetName(),
		TrackingSteps:  []*commondata.TrackingStep{},
		Destination:    firstResult.TrackInfo.DestinationCountryCode,
		ShipmentNumber: firstResult.TrackData.DetailingID,
	}

	for _, event := range firstResult.TrackInfo.TrackEventDetails {
		datetime, _ := time.Parse("2006-01-02T15:04:05", event.CreatedOn)
		location, _ := event.ProcessLocation.(string)

		trackingData.TrackingSteps = append(trackingData.TrackingSteps, &commondata.TrackingStep{
			Datetime:   datetime,
			CommonType: commondata.CommonTrackingStepType_UNKNOWN,
			Message:    event.ProcessContent,
			Location:   location,
		})
	}

	return trackingData, nil
}
