package paczkobot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/commonerrors"
	"github.com/alufers/paczkobot/inpostextra"
	"github.com/alufers/paczkobot/providers"
	"github.com/alufers/paczkobot/providers/inpost"
)

type InpostScannerService struct {
	app *BotApp
}

func NewInpostScannerService(app *BotApp) *InpostScannerService {
	return &InpostScannerService{app: app}
}

func (s *InpostScannerService) ScanUserPackages(ctx context.Context, creds *inpostextra.InpostCredentials) error {
	resp, err := s.app.InpostService.GetUserParcels(ctx, s.app.DB, creds)
	if err != nil {
		return err
	}

	go func() {
		for _, parcel := range resp.Parcels {

			inpostProv := &inpost.InpostProvider{}
			d, err := providers.InvokeProvider(ctx, inpostProv, parcel.ShipmentNumber)
			if errors.Is(err, commonerrors.NotFoundError) {
				continue
			}
			if err != nil {
				log.Printf("failed to invoke inpost provider for shipment number %v: %v", parcel.ShipmentNumber, err)
				continue
			}
			lastStep := &commondata.TrackingStep{}
			if len(d.TrackingSteps) > 0 {
				lastStep = d.TrackingSteps[len(d.TrackingSteps)-1]
			}
			provider := &FollowedPackageProvider{
				ProviderName:       inpostProv.GetName(),
				LastStatusValue:    lastStep.Message,
				LastStatusDate:     lastStep.Datetime,
				LastStatusLocation: lastStep.Location,
			}

			followedPackage := &FollowedPackage{
				InpostCredentials: creds,
				FromName:          parcel.Sender.Name,
			}

			followErr := s.app.FollowService.FollowPackage(
				context.Background(),
				parcel.ShipmentNumber,
				creds.TelegramUserID,
				creds.TelegramChatID,
				[]*FollowedPackageProvider{provider},
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
