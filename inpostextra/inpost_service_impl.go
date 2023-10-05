package inpostextra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

var (
	ErrReauthenticationRequired = fmt.Errorf("reauthentication required")
	ErrRefreshTokenExpired      = fmt.Errorf("refresh token expired")
)

type InpostServiceImpl struct {
	BaseURL    string
	httpClient *http.Client
}

func NewInpostService(
	httpClient *http.Client,
) InpostService {
	return &InpostServiceImpl{
		BaseURL:    "https://api-inmobile-pl.easypack24.net",
		httpClient: httpClient,
	}
}

func (s *InpostServiceImpl) makeJSONRequest(
	ctx context.Context,
	creds *InpostCredentials, method string, path string, data interface{}, respRef interface{},
) error {
	var bodyBuf io.Reader
	if method != "GET" && method != "HEAD" {
		body, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("error marshalling JSON data: %s", err)
		}
		bodyBuf = bytes.NewBuffer(body)
	}
	// log.Printf("Sending %s request to %s", method, s.BaseURL+path)
	if bodyBuf != nil {
		log.Printf("Body: %s", bodyBuf)
	}
	httpReq, err := http.NewRequestWithContext(ctx, method, s.BaseURL+path, bodyBuf)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %s", err)
	}
	if bodyBuf != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	httpReq.Header.Set("Accept", "application/json")
	if creds != nil {
		httpReq.Header.Set("Authorization", creds.AuthToken)
	}
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		data, _ := ioutil.ReadAll(resp.Body)
		dataLen := 1024
		if len(data) < dataLen {
			dataLen = len(data)
		}
		return fmt.Errorf("HTTP request failed with status code %s: %v", resp.Status, string(data[0:dataLen]))
	}
	if respRef != nil {
		bodyBuf := bytes.NewBuffer(nil)
		_, err = io.Copy(bodyBuf, resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %s", err)
		}
		// log.Printf("Response: %s", bodyBuf.String())
		err = json.Unmarshal(bodyBuf.Bytes(), respRef)

		if err != nil {
			return fmt.Errorf("error decoding JSON response: %s", err)
		}
	}
	return nil
}

func (s *InpostServiceImpl) SendSMSCode(ctx context.Context, phoneNumber string) error {
	phoneNumber = NormalizePhoneNumber(phoneNumber)

	// make a post request to /v1/sendSMSCode with the phone number sent as JSON
	data := map[string]string{
		"phoneNumber": phoneNumber,
	}
	return s.makeJSONRequest(ctx, nil, "POST", "/v1/sendSMSCode", data, nil)
}

func (s *InpostServiceImpl) ConfirmSMSCode(ctx context.Context, phoneNumber string, code string) (*InpostCredentials, error) {
	phoneNumber = NormalizePhoneNumber(phoneNumber)
	// make a post request to /v1/confirmSMSCode with the phone number and code sent as JSON
	data := map[string]string{
		"phoneNumber": phoneNumber,
		"smsCode":     code,
		"phoneOS":     "Android",
	}

	out := &struct {
		AuthToken    string `json:"authToken"`
		RefreshToken string `json:"refreshToken"`
	}{}
	err := s.makeJSONRequest(ctx, nil, "POST", "/v1/confirmSMSCode", data, out)
	if err != nil {
		return nil, err
	}

	creds := &InpostCredentials{
		AuthToken:    out.AuthToken,
		RefreshToken: out.RefreshToken,
		PhoneNumber:  phoneNumber,
	}
	return creds, nil
}

// Authenticate uses the refresh token to get a new access token.
func (s *InpostServiceImpl) Authenticate(ctx context.Context, creds *InpostCredentials) error {
	// make a post request to /v1/authenticate with the access token and refresh token sent as JSON
	data := map[string]string{
		"phoneOS":      "Android",
		"refreshToken": creds.RefreshToken,
	}
	out := &struct {
		AuthToken                string     `json:"authToken"`
		ReauthenticationRequired bool       `json:"reauthenticationRequired"`
		RefreshTokenExpiryDate   *time.Time `json:"refreshTokenExpiryDate"`
	}{}
	err := s.makeJSONRequest(ctx, nil, "POST", "/v1/authenticate", data, out)
	if err != nil {
		return err
	}
	if out.ReauthenticationRequired {
		return ErrReauthenticationRequired
	}
	log.Printf("RefreshTokenExpiryDate: %v", out.RefreshTokenExpiryDate)
	if out.RefreshTokenExpiryDate != nil && out.RefreshTokenExpiryDate.Before(time.Now()) {
		return ErrRefreshTokenExpired
	}
	creds.AuthToken = out.AuthToken
	return nil
}

func (s *InpostServiceImpl) ReauthenticateIfNeeded(ctx context.Context, db *gorm.DB, creds *InpostCredentials) error {
	if creds.RefreshToken == "" {
		return fmt.Errorf("refresh token is empty")
	}

	log.Printf("[ReauthenticateIfNeeded] Creds: %+v", creds)

	tok, _ := jwt.ParseWithClaims(strings.TrimPrefix(creds.AuthToken, "Bearer "), &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("unused"), nil
	})
	claims := tok.Claims.(*jwt.StandardClaims)
	expirationDate := time.Unix(claims.ExpiresAt, 0)
	if expirationDate.Before(time.Now()) {
		log.Printf("Inpost token for phone number %v expired %v ago, refreshing...", creds.PhoneNumber, time.Since(expirationDate))
		err := s.Authenticate(ctx, creds)
		if err != nil {
			return err
		}
		if err := db.Save(creds).Error; err != nil {
			return fmt.Errorf("error saving credentials for phone number %v: %v", creds.PhoneNumber, err)
		}
		return nil
	}

	return nil
}

func (s *InpostServiceImpl) GetParcel(ctx context.Context, db *gorm.DB, creds *InpostCredentials, shipmentNumber string) (*InpostParcel, error) {
	err := s.ReauthenticateIfNeeded(ctx, db, creds)
	if err != nil {
		return nil, err
	}

	// make a get request to /v1/parcel/{shipmentNumber}
	out := &InpostParcel{}
	err = s.makeJSONRequest(ctx, creds, "GET", fmt.Sprintf("/v3/parcels/tracked/%s", shipmentNumber), nil, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *InpostServiceImpl) GetUserParcels(ctx context.Context, db *gorm.DB, creds *InpostCredentials) (*GetTrackedParcelsResponse, error) {
	err := s.ReauthenticateIfNeeded(ctx, db, creds)
	if err != nil {
		return nil, err
	}

	// make a get request to parcels
	out := &GetTrackedParcelsResponse{}
	err = s.makeJSONRequest(ctx, creds, "GET", "/v3/parcels/tracked", nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *InpostServiceImpl) OpenParcelLocker(ctx context.Context, db *gorm.DB, creds *InpostCredentials, shipmentNumber string) error {
	parcel, err := s.GetParcel(ctx, db, creds, shipmentNumber)
	if err != nil {
		return fmt.Errorf("error getting parcel details before opening: %v", err)
	}
	if parcel.PickUpPoint == nil || parcel.PickUpPoint.Location == nil {
		return fmt.Errorf("parcel %v has no pickup point or location", shipmentNumber)
	}
	if parcel.Receiver == nil || parcel.Receiver.PhoneNumber == "" {
		return fmt.Errorf("parcel %v has no receiver or receiver phone number", shipmentNumber)
	}
	startOpenSessionRequest := map[string]any{
		"geoPoint": map[string]any{
			"accuracy":  10.0,
			"latitude":  parcel.PickUpPoint.Location.Latitude,
			"longitude": parcel.PickUpPoint.Location.Longitude,
		},
		"parcel": map[string]any{
			"shipmentNumber":      shipmentNumber,
			"openCode":            parcel.OpenCode,
			"receiverPhoneNumber": parcel.Receiver.PhoneNumber,
		},
	}
	validateResp := &ValidateCompartmentResponse{}
	err = s.makeJSONRequest(ctx, creds, "POST", "/v1/collect/validate", startOpenSessionRequest, validateResp)
	if err != nil {
		return fmt.Errorf("error validating compartment: %v", err)
	}

	openResp := make(map[string]any)
	err = s.makeJSONRequest(ctx, creds, "POST", "/v1/collect/compartment/open/"+validateResp.SessionUUID, map[string]any{
		"sessionUuid": validateResp.SessionUUID,
	}, &openResp)
	if err != nil {
		return fmt.Errorf("error opening compartment: %v", err)
	}
	return nil
}
