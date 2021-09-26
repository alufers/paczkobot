package mock

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/alufers/paczkobot/commondata"
)

type MockProvider struct {
}

var MockNumberRegex = regexp.MustCompile(`^mock-`)

func (ip *MockProvider) GetName() string {
	return "mock"
}

func (ip *MockProvider) MatchesNumber(trackingNumber string) bool {
	return MockNumberRegex.MatchString(trackingNumber)
}

func (ip *MockProvider) Track(ctx context.Context, trackingNumber string) (*commondata.TrackingData, error) {
	f, err := os.Open("mock-package.json")
	if err != nil {
		return nil, fmt.Errorf("error opening mock-package.json: %v", err)
	}
	defer f.Close()
	out := &commondata.TrackingData{}
	if err = json.NewDecoder(f).Decode(out); err != nil {
		return nil, fmt.Errorf("error decoding mock-package.json: %v", err)
	}
	out.ProviderName = ip.GetName()
	out.ShipmentNumber = trackingNumber
	return out, nil
}
