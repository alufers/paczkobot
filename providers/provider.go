package providers

import (
	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/providers/caniao"
	"github.com/alufers/paczkobot/providers/inpost"
	"github.com/alufers/paczkobot/providers/pocztapolska"
	"github.com/alufers/paczkobot/providers/postnl"
)

var AllProviders = []Provider{
	&inpost.InpostProvider{},
	&pocztapolska.PocztaPolskaProvider{},
	&postnl.PostnlProvider{},
	&caniao.CaniaoProvider{},
}

type Provider interface {
	GetName() string
	MatchesNumber(trackingNumber string) bool
	Track(trackingNumber string) (*commondata.TrackingData, error)
}
