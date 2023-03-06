package paczkobot

import (
	"fmt"

	"github.com/alufers/paczkobot/commondata"
	"github.com/xeonx/timeago"
)

type PackagePrinterService struct{}

func NewPackagePrinterService() *PackagePrinterService {
	return &PackagePrinterService{}
}

// PrintPackages prints packages for being listed in commands such as /packages
func (s *PackagePrinterService) PrintPackages(followedPackages []FollowedPackageTelegramUser) string {
	packagesText := ""
	for _, p := range followedPackages {
		customName := p.CustomName

		if p.FollowedPackage.FromName != "" {
			if customName != "" {
				customName += " "
			}
			customName += fmt.Sprintf("from %s", p.FollowedPackage.FromName)
		}
		if customName != "" {
			customName = fmt.Sprintf(" <i>(%s)</i>", customName)
		}
		packagesText += fmt.Sprintf("<b>%v</b>%v", p.FollowedPackage.TrackingNumber, customName)
		for i, prov := range p.FollowedPackage.FollowedPackageProviders {
			emojiText := ""
			if prov.LastStatusCommonType != commondata.CommonTrackingStepType_UNKNOWN {
				emojiText = commondata.CommonTrackingStepTypeEmoji[prov.LastStatusCommonType]
				if emojiText != "" {
					emojiText = " " + emojiText
				}
			}
			packagesText += fmt.Sprintf(" %v (<i>%v%v %v</i>)",
				prov.ProviderName,
				emojiText,
				prov.LastStatusValue,
				timeago.English.Format(prov.LastStatusDate))
			if i != len(p.FollowedPackage.FollowedPackageProviders)-1 {
				packagesText += ", "
			}
		}

		packagesText += "\n"
	}

	return packagesText
}
