package providers

import (
	"context"

	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/providers/deutsche_post"
	"github.com/alufers/paczkobot/providers/dhl"
	"github.com/alufers/paczkobot/providers/dpdcompl"
	"github.com/alufers/paczkobot/providers/fedex_pl"
	"github.com/alufers/paczkobot/providers/geis_pl"
	"github.com/alufers/paczkobot/providers/gls"
	"github.com/alufers/paczkobot/providers/inpost"
	"github.com/alufers/paczkobot/providers/orlen"
	"github.com/alufers/paczkobot/providers/packeta"
	"github.com/alufers/paczkobot/providers/pocztapolska"
	"github.com/alufers/paczkobot/providers/postnl"
	"github.com/alufers/paczkobot/providers/ups"
	"github.com/alufers/paczkobot/providers/yuntrack"
)

var AllProviders = []Provider{
	&inpost.InpostProvider{},
	&pocztapolska.PocztaPolskaProvider{},
	&postnl.PostnlProvider{},
	// &caniao.CaniaoProvider{},
	&dpdcompl.DpdComPlProvider{},
	&ups.UPSProvider{},
	&dhl.DHLProvider{},
	&gls.GLSProvider{},
	&yuntrack.YunTrack{},
	&packeta.PacketaProvider{},
	&fedex_pl.FedexPlProvider{},
	&geis_pl.GeisPlProvider{},
	&orlen.OrlenProvider{},
	&deutsche_post.DeutschePostProvider{},
}

type Provider interface {
	GetName() string
	MatchesNumber(trackingNumber string) bool
	Track(ctx context.Context, trackingNumber string) (*commondata.TrackingData, error)
}

type ProviderWithAutoArchive interface {
	Provider
	ShouldAutoArchive(*commondata.TrackingData) bool
}

func GetProviderByName(name string) Provider {
	for _, provider := range AllProviders {
		if provider.GetName() == name {
			return provider
		}
	}
	return nil
}
