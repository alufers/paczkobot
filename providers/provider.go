package providers

var AllProviders = []Provider{
	&InpostProvider{},
	&PocztaPolskaProvider{},
}

type Provider interface {
	GetName() string
	MatchesNumber(trackingNumber string) bool
	Track(trackingNumber string) (*TrackingData, error)
}
