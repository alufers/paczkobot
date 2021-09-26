package paczkobot

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/providers"
	"gorm.io/gorm"
)

type TrackingService struct {
	app *BotApp
}

func NewTrackingService(app *BotApp) *TrackingService {
	return &TrackingService{
		app: app,
	}
}

func (ts *TrackingService) InvokeProviderAndNotifyFollowers(ctx context.Context, provider providers.Provider, trackingNumber string) (result *commondata.TrackingData, err error) {
	result, err = providers.InvokeProvider(ctx, provider, trackingNumber)
	if err != nil {
		return
	}

	followedPackage := &FollowedPackage{}

	if err := ts.app.DB.Where("tracking_number = ?", trackingNumber).
		Preload("FollowedPackageProviders").
		Preload("FollowedPackageTelegramUsers").
		First(&followedPackage).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to find followers for package: %w", err)
		}
	} else {
		var providerToUpdate *FollowedPackageProvider
		for _, packageProvider := range followedPackage.FollowedPackageProviders {
			if packageProvider.ProviderName == provider.GetName() {
				providerToUpdate = packageProvider
			}
		}
		// create a new provider if it doesn't exist
		if providerToUpdate == nil {
			providerToUpdate = &FollowedPackageProvider{
				FollowedPackage: followedPackage,
				ProviderName:    provider.GetName(),
			}
			followedPackage.FollowedPackageProviders = append(followedPackage.FollowedPackageProviders, providerToUpdate)
		}
		lastTrackingStep := result.TrackingSteps[len(result.TrackingSteps)-1]
		if providerToUpdate.LastStatusValue != lastTrackingStep.Message || math.Abs(float64(providerToUpdate.LastStatusDate.Sub(lastTrackingStep.Datetime))) > float64(time.Minute) {
			providerToUpdate.LastStatusValue = lastTrackingStep.Message
			providerToUpdate.LastStatusDate = lastTrackingStep.Datetime
			providerToUpdate.LastStatusLocation = lastTrackingStep.Location

			if err := ts.app.DB.Save(&providerToUpdate).Error; err != nil {
				return nil, fmt.Errorf("failed to save followed package: %w", err)
			}
		}
		ts.app.NotificationsService.NotifyProviderStatusChanged(provider, followedPackage)
	}

	return result, nil
}
