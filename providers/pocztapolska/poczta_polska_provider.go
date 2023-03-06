package pocztapolska

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/alufers/paczkobot/commondata"

	"github.com/alufers/paczkobot/commonerrors"
	"github.com/alufers/paczkobot/providers/pocztapolska/sledzeniehttpbinding"
)

var codeMappings = map[string]commondata.CommonTrackingStepType{
	"P_NAD": commondata.CommonTrackingStepType_SENT,
	"P_WD":  commondata.CommonTrackingStepType_OUT_FOR_DELIVERY,
	"P_D":   commondata.CommonTrackingStepType_DELIVERED,
}

type PocztaPolskaProvider struct{}

func (ip *PocztaPolskaProvider) GetName() string {
	return "poczta-polska"
}

func (ip *PocztaPolskaProvider) MatchesNumber(trackingNumber string) bool {
	return true
}

type SledzEnvelope struct {
	XMLName      xml.Name `xml:"Envelope"`
	EnvelopeAttr string   `xml:"xmlns:SOAP-ENV,attr"`
	NSAttr       string   `xml:"xmlns:ns,attr"`
	TNSAttr      string   `xml:"xmlns:tns,attr,omitempty"`
	URNAttr      string   `xml:"xmlns:urn,attr,omitempty"`
	XSIAttr      string   `xml:"xmlns:xsi,attr,omitempty"`

	Body *SledzEnvelopeBody `xml:"Body"`
}

type SledzEnvelopeBody struct {
	XMLName                  xml.Name                                       `xml:"Body"`
	SprawdzPrzesylkeResponse *sledzeniehttpbinding.SprawdzPrzesylkeResponse `xml:"sprawdzPrzesylkeResponse"`
}

func EscapeXML(d string) string {
	buf := &bytes.Buffer{}
	xml.Escape(buf, []byte(d))
	return buf.String()
}

func (ip *PocztaPolskaProvider) Track(ctx context.Context, trackingNumber string) (*commondata.TrackingData, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", "https://tt.poczta-polska.pl/Sledzenie/services/Sledzenie?wsdl", strings.NewReader(fmt.Sprintf(`
	<soapenv:Envelope   xmlns:sled="http://sledzenie.pocztapolska.pl" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
    <soapenv:Header>
        <wsse:Security           soapenv:mustUnderstand="1"          xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd">
            <wsse:UsernameToken wsu:Id="UsernameToken-2"             xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
                <wsse:Username>sledzeniepp</wsse:Username>
                <wsse:Password                  Type="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText">PPSA</wsse:Password>
                <wsse:Nonce                  EncodingType="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-soap-message-security-1.0#Base64Binary">X41PkdzntfgpowZsKegMFg==</wsse:Nonce>
                <wsu:Created>2011-12-08T07:59:28.656Z</wsu:Created>
            </wsse:UsernameToken>
        </wsse:Security>
    </soapenv:Header>
    <soapenv:Body>
        <sled:sprawdzPrzesylke>
            <sled:numer>%v</sled:numer>
        </sled:sprawdzPrzesylke>
    </soapenv:Body>
</soapenv:Envelope>
	`, EscapeXML(trackingNumber))))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}
	req.Header.Add("Content-type", "text/xml")
	httpResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, commonerrors.NewNetworkError(ip.GetName(), req)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP status code: %v", httpResponse.StatusCode)
	}
	data, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read XMl response from PP: %w", err)
	}
	// log.Printf("xmlData = %v", string(data))

	resp := &SledzEnvelope{}
	if err := xml.Unmarshal(data, resp); err != nil {
		return nil, fmt.Errorf("failed to decode XMl response from PP: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to request poczta polska: %v", err)
	}
	td := &commondata.TrackingData{
		ShipmentNumber: trackingNumber,
		ProviderName:   ip.GetName(),
		TrackingSteps:  []*commondata.TrackingStep{},
	}
	if resp.Body.SprawdzPrzesylkeResponse.Return == nil ||
		resp.Body.SprawdzPrzesylkeResponse.Return.DanePrzesylki == nil ||
		resp.Body.SprawdzPrzesylkeResponse.Return.DanePrzesylki.Zdarzenia == nil {
		return nil, commonerrors.NotFoundError
	}
	for _, z := range resp.Body.SprawdzPrzesylkeResponse.Return.DanePrzesylki.Zdarzenia.Zdarzenie {
		t, _ := time.Parse("2006-01-02 15:04", *z.Czas)
		var commonType commondata.CommonTrackingStepType
		if v, ok := codeMappings[*z.Kod]; ok {
			commonType = v
		} else {
			commonType = commondata.CommonTrackingStepType_UNKNOWN
		}
		td.TrackingSteps = append(td.TrackingSteps, &commondata.TrackingStep{
			Datetime:   t,
			CommonType: commonType,
			Message:    *z.Nazwa,
			Location:   *z.Jednostka.Nazwa,
		})
	}
	up := resp.Body.SprawdzPrzesylkeResponse.Return.DanePrzesylki.UrzadPrzezn
	if up != nil && up.Nazwa != nil {
		td.Destination = *up.Nazwa
	}
	return td, nil
}
