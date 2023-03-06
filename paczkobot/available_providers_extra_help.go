package paczkobot

import (
	"fmt"
	"strings"

	"github.com/alufers/paczkobot/providers"
)

type AvailableProvidersExtraHelp struct{}

func (a *AvailableProvidersExtraHelp) Help() string {
	providerNames := []string{}
	for _, p := range providers.AllProviders {
		providerNames = append(providerNames, p.GetName())
	}
	return fmt.Sprintf("Available providers:\n%v", strings.Join(providerNames, ", "))
}
