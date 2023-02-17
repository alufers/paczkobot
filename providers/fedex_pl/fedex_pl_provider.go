package fedex_pl

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/alufers/paczkobot/commondata"
)

// https://www.fedex.com/pl-pl/online/domestic-tracking.html

var commonMappings = map[string]commondata.CommonTrackingStepType{
	"Przesyłka wydana do doręczenia kurierowi FedEx.": commondata.CommonTrackingStepType_OUT_FOR_DELIVERY,
	"Przesyłka w oddziale FedEx.":                     commondata.CommonTrackingStepType_IN_TRANSIT,
	"Przesyłka doręczona do odbiorcy.":                commondata.CommonTrackingStepType_DELIVERED,
}

var (
	// not used
	// FedexClientID     = "l7xx474b79016a4d4ec5a60bf7a7e5e7e6fe"
	// FedexClientSecret = "448399ccafaa4f62a4ed202fc5ef3a01"

	FedexStaticAPIKey = "l7xx492b4e2b8682483c979222bdd33216cf"
)

type FedexPlProvider struct {
}

func (fp *FedexPlProvider) GetName() string {
	return "fedex-pl"
}

func (fp *FedexPlProvider) MatchesNumber(trackingNumber string) bool {
	return true
}

func (fp *FedexPlProvider) Track(ctx context.Context, trackingNumber string) (*commondata.TrackingData, error) {

	trackingReq, err := http.NewRequestWithContext(ctx, "GET", "https://api1.emea.fedex.com/fds2-tracking/trck-v1/info?trackingKey="+url.QueryEscape(trackingNumber), nil)
	if err != nil {
		return nil, err
	}
	commondata.SetCommonHTTPHeaders(&trackingReq.Header)
	trackingReq.Header.Add("apikey", FedexStaticAPIKey)
	resp, err := http.DefaultClient.Do(trackingReq)
	if err != nil {
		return nil, err
	}

	respData := &FedexPlTrackingResponse{}
	err = json.NewDecoder(resp.Body).Decode(respData)
	if err != nil {
		return nil, err
	}

	td := &commondata.TrackingData{
		ShipmentNumber: trackingNumber,
		ProviderName:   fp.GetName(),
		TrackingSteps:  []*commondata.TrackingStep{},
		Destination:    respData.DeliveryDepot,
		Weight:         respData.Weight,
		SentFrom:       respData.ShipmentDepot,
	}

	for _, ev := range respData.Events {
		// example date: 2023-02-16T21:40:14+0100

		datetime, _ := time.Parse("2006-01-02T15:04:05-0700", ev.EventDate)
		td.TrackingSteps = append(td.TrackingSteps, &commondata.TrackingStep{
			Location: ev.Depot,
			Datetime: datetime,
			Message:  ev.EventName,
		})
	}

	td.ApplyCommonTypeMappings(commonMappings)

	return td, nil
}

// not used again xD
// func (fp *FedexPlProvider) getOauthToken(ctx context.Context) (string, error) {
// 	// make a GET reqeuest to https://www.fedex.com/etc/services/getapigconfigs.jsonp

// 	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.fedex.com/etc/services/getapigconfigs.jsonp", nil)
// 	if err != nil {
// 		return "", err
// 	}
// 	commondata.SetCommonHTTPHeaders(&req.Header)
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return "", err
// 	}
// 	respData := &FedexPlApiConfigs{}
// 	err = json.NewDecoder(resp.Body).Decode(respData)
// 	if err != nil {
// 		return "", err
// 	}
// 	return respData.ClientID, nil
// }

// Lol, used the wrong endpoint. Might be useful later.
// func (fp *FedexPlProvider) getOauthToken(ctx context.Context) (string, error) {
// 	var tokenReqParams = url.Values{
// 		"client_id":     []string{FedexClientID},
// 		"client_secret": []string{FedexClientSecret},
// 		"grant_type":    []string{"client_credentials"},
// 	}
// 	tokenReqUrl, err := url.Parse("https://api.fedex.com/auth/oauth/v2/token")
// 	if err != nil {
// 		return "", err
// 	}
// 	tokenReqUrl.RawQuery = tokenReqParams.Encode()
// 	oauthTokenRequest, err := http.NewRequestWithContext(ctx, "POST", tokenReqUrl.String(), nil)
// 	if err != nil {
// 		return "", err
// 	}
// 	commondata.SetCommonHTTPHeaders(&oauthTokenRequest.Header)
// 	oauthTokenRequest.Header.Add("Content-type", "application/x-www-form-urlencoded")
// 	openIdTokenResponse, err := http.DefaultClient.Do(oauthTokenRequest)
// 	if err != nil {
// 		return "", err
// 	}
// 	if openIdTokenResponse.StatusCode != http.StatusOK {
// 		errorContent, err := ioutil.ReadAll(openIdTokenResponse.Body)
// 		if err != nil {
// 			return "", fmt.Errorf("failed to get OAuth token: %s", openIdTokenResponse.Status)
// 		}
// 		return "", fmt.Errorf("failed to get OAuth token (status: %s), error: %v", openIdTokenResponse.Status, string(errorContent))
// 	}

// 	resp := &FedexPlTokenResponse{}
// 	err = json.NewDecoder(openIdTokenResponse.Body).Decode(resp)
// 	if err != nil {
// 		return "", err
// 	}
// 	return resp.AccessToken, nil
// }
