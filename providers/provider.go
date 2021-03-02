package providers

import (
	"context"
	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/providers/caniao"
	"github.com/alufers/paczkobot/providers/dpdcompl"
	"github.com/alufers/paczkobot/providers/inpost"
	"github.com/alufers/paczkobot/providers/pocztapolska"
	"github.com/alufers/paczkobot/providers/postnl"
)

var AllProviders = []Provider{
	&inpost.InpostProvider{},
	&pocztapolska.PocztaPolskaProvider{},
	&postnl.PostnlProvider{},
	&caniao.CaniaoProvider{},
	&dpdcompl.DpdComPlProvider{},
}

type Provider interface {
	GetName() string
	MatchesNumber(trackingNumber string) bool
	Track(ctx context.Context, trackingNumber string) (*commondata.TrackingData, error)
}
