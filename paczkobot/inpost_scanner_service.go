package paczkobot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alufers/paczkobot/inpostextra"
)

type InpostScannerService struct {
	app *BotApp
}

func NewInpostScannerService(app *BotApp) *InpostScannerService {
	return &InpostScannerService{app: app}
}

func (s *InpostScannerService) ScanUserPackages(creds *inpostextra.InpostCredentials) error {
	parcels, err := s.app.InpostService.GetUserParcels(s.app.DB, creds)
	if err != nil {
		return err
	}

	go func() {
		for _, parcel := range parcels {
			followedPackage := &FollowedPackage{
				InpostCredentials: creds,
				FromName:          parcel.SenderName,
			}
			lastStep := parcel.StatusHistory[len(parcel.StatusHistory)-1]
			followErr := s.app.FollowService.FollowPackage(
				context.Background(),
				parcel.ShipmentNumber,
				creds.TelegramUserID,
				creds.TelegramChatID,
				[]*FollowedPackageProvider{
					{
						ProviderName:    "inpost",
						LastStatusDate:  lastStep.Date,
						LastStatusValue: lastStep.Status,
					},
				},
				followedPackage,
			)
			if followErr != nil {
				log.Printf("failed to follow package %v after inpost scan for user id %v: %v", parcel.ShipmentNumber, creds.TelegramUserID, followErr)
				continue
			}

		}
	}()

	creds.LastScan.Time = time.Now()
	creds.LastScan.Valid = true
	if err := s.app.DB.Save(creds).Error; err != nil {
		return fmt.Errorf("failed to save InpostCredentials: %v", err)
	}

	return nil
}
