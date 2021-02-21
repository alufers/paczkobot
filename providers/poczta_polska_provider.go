package providers

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alufers/paczkobot/commonerrors"
	"github.com/alufers/paczkobot/providers/sledzeniehttpbinding"
	"github.com/davecgh/go-spew/spew"
)

type PocztaPolskaProvider struct {
}

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

func (ip *PocztaPolskaProvider) Track(trackingNumber string) (*TrackingData, error) {

	req, err := http.NewRequest("POST", "https://tt.poczta-polska.pl/Sledzenie/services/Sledzenie?wsdl", strings.NewReader(fmt.Sprintf(`
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
		return nil, fmt.Errorf("failed to make POST request to %v: %w", req.URL.String(), err)
	}

	if httpResponse.StatusCode != 200 {
		return nil, fmt.Errorf("POST request to PP failed with status code %v: %w", httpResponse.StatusCode, err)
	}
	data, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read XMl response from PP: %w", err)
	}
	log.Printf("xmlData = %v", string(data))

	var resp = &SledzEnvelope{}
	if err := xml.Unmarshal(data, resp); err != nil {
		return nil, fmt.Errorf("failed to decode XMl response from PP: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to request poczta polska: %v", err)
	}
	td := &TrackingData{
		ShipmentNumber: trackingNumber,
		ProviderName:   ip.GetName(),
		TrackingSteps:  []*TrackingStep{},
	}
	spew.Dump(resp)
	if resp.Body.SprawdzPrzesylkeResponse.Return.DanePrzesylki.Zdarzenia == nil {
		return nil, commonerrors.NotFoundError
	}
	for _, z := range resp.Body.SprawdzPrzesylkeResponse.Return.DanePrzesylki.Zdarzenia.Zdarzenie {
		log.Printf("%#v", *z)
		t, _ := time.Parse("2006-01-02 15:04", *z.Czas)
		td.TrackingSteps = append(td.TrackingSteps, &TrackingStep{
			Datetime:   t,
			CommonType: *z.Kod,
			Message:    *z.Nazwa,
			Location:   *z.Jednostka.Nazwa,
		})
	}
	return td, nil

}
