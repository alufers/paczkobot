package providers

import (
	"context"

	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/providers/caniao"
	"github.com/alufers/paczkobot/providers/dhl"
	"github.com/alufers/paczkobot/providers/dpdcompl"
	"github.com/alufers/paczkobot/providers/gls"
	"github.com/alufers/paczkobot/providers/inpost"
	"github.com/alufers/paczkobot/providers/pocztapolska"
	"github.com/alufers/paczkobot/providers/postnl"
	"github.com/alufers/paczkobot/providers/ups"
)

var AllProviders = []Provider{
	&inpost.InpostProvider{},
	&pocztapolska.PocztaPolskaProvider{},
	&postnl.PostnlProvider{},
	&caniao.CaniaoProvider{},
	&dpdcompl.DpdComPlProvider{},
	&ups.UPSProvider{},
	&dhl.DHLProvider{},
	&gls.GLSProvider{},
}

type Provider interface {
	GetName() string
	MatchesNumber(trackingNumber string) bool
	Track(ctx context.Context, trackingNumber string) (*commondata.TrackingData, error)
}

func GetProviderByName(name string) Provider {
	for _, provider := range AllProviders {
		if provider.GetName() == name {
			return provider
		}
	}
	return nil
}
