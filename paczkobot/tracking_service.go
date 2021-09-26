package paczkobot

import (
	"context"

	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/providers"
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

	// followedPackage := &FollowedPackage{}

	// if ts.app.DB.Where("tracking_number")

	return nil, nil
}
