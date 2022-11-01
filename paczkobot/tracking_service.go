package paczkobot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/inpostextra"
	"github.com/alufers/paczkobot/providers"
	"github.com/spf13/viper"
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
		if _, err := ts.notifyFollowersOfPackageIfNeeded(ctx, followedPackage, provider, result); err != nil {
			return nil, fmt.Errorf("failed to notify followers of package: %w", err)
		}
	}

	return result, nil
}

func (ts *TrackingService) notifyFollowersOfPackageIfNeeded(ctx context.Context, followedPackage *FollowedPackage, provider providers.Provider, result *commondata.TrackingData) (bool, error) {
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
	if len(result.TrackingSteps) == 0 {
		log.Printf("[WARN]: tracking result for %v is empty", followedPackage.TrackingNumber)
		return false, nil
	}
	lastTrackingStep := result.TrackingSteps[len(result.TrackingSteps)-1]
	if providerToUpdate.LastStatusValue != lastTrackingStep.Message || // message changed
		math.Abs(float64(providerToUpdate.LastStatusDate.Sub(lastTrackingStep.Datetime))) > float64(time.Minute) || // time changed more than 1 minute
		providerToUpdate.LastStatusCommonType != lastTrackingStep.CommonType { // type changed (usually due to a schema update in paczkobot)
		providerToUpdate.LastStatusValue = lastTrackingStep.Message
		providerToUpdate.LastStatusDate = lastTrackingStep.Datetime
		providerToUpdate.LastStatusLocation = lastTrackingStep.Location
		providerToUpdate.LastStatusCommonType = lastTrackingStep.CommonType

		if err := ts.app.DB.Save(&providerToUpdate).Error; err != nil {
			return false, fmt.Errorf("failed to save followed package provider: %w", err)
		}
		followedPackage.LastChange = time.Now()
		followedPackage.Inactive = false
		if err := ts.app.DB.Save(followedPackage).Error; err != nil {
			return false, fmt.Errorf("failed to save followed package: %w", err)
		}
		return true, ts.app.NotificationsService.NotifyProviderStatusChanged(provider, followedPackage)
	}

	return false, nil

}

func (ts *TrackingService) MarkPackagesWithoutChangesAsInactive() error {
	var followedPackages []FollowedPackage
	if err := ts.app.DB.Where("last_change < ? AND inactive = false", time.Now().Add(-viper.GetDuration("tracking.max_time_without_change"))).
		Find(&followedPackages).Error; err != nil {
		return fmt.Errorf("failed to find followed expired packages: %w", err)
	}

	for _, followedPackage := range followedPackages {
		followedPackage.Inactive = true
		if err := ts.app.DB.Save(&followedPackage).Error; err != nil {
			return fmt.Errorf("failed to mark package as inactive: %w", err)
		}
	}

	log.Printf("Marked %v packages as inactive", len(followedPackages))

	return nil
}

func (ts *TrackingService) RunAutomaticTrackingLoop() {
	log.Printf("Starting automatic tracking loop...")
	for {
		lastCheckStarted := time.Now()
		if err := ts.MarkPackagesWithoutChangesAsInactive(); err != nil {
			log.Printf("Failed to mark packages without changes as inactive: %v", err)
		}
		var followedPackages []*FollowedPackage
		if err := ts.app.DB.
			Where("inactive = false AND last_automatic_check < ?", time.Now().Add(-viper.GetDuration("tracking.automatic_tracking_check_interval"))).
			Order("last_automatic_check ASC").
			Limit(viper.GetInt("tracking.max_packages_per_automatic_tracking_check")).
			Preload("FollowedPackageProviders").
			Preload("FollowedPackageTelegramUsers").
			Find(&followedPackages).Error; err != nil {
			log.Printf("failed to find packages to track automatically: %v", err)
			return
		}

		totalWaitingPackages := int64(0)
		// count total number of packages to track
		if err := ts.app.DB.
			Model(&FollowedPackage{}).
			Where(
				"inactive = false AND last_automatic_check < ?",
				time.Now().Add(-viper.GetDuration("tracking.automatic_tracking_check_interval")),
			).
			Count(&totalWaitingPackages).Error; err != nil {
			log.Printf("failed to count packages waiting to be tracked automatically: %v", err)
			return
		}

		log.Printf("Checking %v packages for updates in this interval out of %v all packages elegible for tracking...", len(followedPackages), totalWaitingPackages)

		for _, followedPackage := range followedPackages {
			if err := ts.runAutomaticTrackingForPackage(followedPackage); err != nil {
				log.Printf("failed to track package %v automatically: %v", followedPackage.TrackingNumber, err)
			}
			time.Sleep(viper.GetDuration("tracking.delay_between_packages_in_automatic_tracking") - time.Duration(rand.Intn(10000)))
		}

		jitterModulo := int(viper.GetDuration("tracking.automatic_tracking_check_jitter")) - viper.GetInt("tracking.automatic_tracking_check_jitter")/2
		jitterValue := time.Second * 0
		if jitterModulo > 0 {
			jitterValue = time.Second * time.Duration(rand.Intn(jitterModulo))
		}

		timeToWait := viper.GetDuration("tracking.automatic_tracking_check_interval") - time.Since(lastCheckStarted) + jitterValue
		log.Printf("Automatic tracking finished, now scanning inpost accounts...")
		err := ts.ScanInpostAccounts()
		if err != nil {
			log.Printf("Failed to scan inpost accounts: %v", err)
		}
		log.Printf("Automatic tracking finished, next check scheduled in %v", timeToWait)
		if timeToWait > 0 {
			time.Sleep(timeToWait)
		}
	}
}

func (ts *TrackingService) runAutomaticTrackingForPackage(pkg *FollowedPackage) error {
	for _, prov := range pkg.FollowedPackageProviders {
		provider := providers.GetProviderByName(prov.ProviderName)
		if provider == nil {
			return fmt.Errorf("failed to find provider %v", prov.ProviderName)
		}
		ctx := context.Background()
		result, err := providers.InvokeProvider(ctx, provider, pkg.TrackingNumber)
		if err != nil {
			return fmt.Errorf("error from provider %v: %w", provider.GetName(), err)
		}
		didChange, err := ts.notifyFollowersOfPackageIfNeeded(ctx, pkg, provider, result)
		if err != nil {
			return fmt.Errorf("failed to notify followers of package: %w", err)
		}
		pkg.LastAutomaticCheck = time.Now()
		// save the package in the database
		if err := ts.app.DB.Save(pkg).Error; err != nil {
			return fmt.Errorf("failed to save package: %w", err)
		}
		if didChange {
			log.Printf("Package %v (%v) changed! -> %v", pkg.TrackingNumber, prov.ProviderName, result.TrackingSteps[len(result.TrackingSteps)-1].Message)
		}
	}
	return nil
}

func (ts *TrackingService) ScanInpostAccounts() error {
	log.Printf("Starting inpost account scan...")

	lastCheckStarted := time.Now()
	var inpostCreds = []*inpostextra.InpostCredentials{}
	if err := ts.app.DB.
		Where("last_scan < ? OR last_scan IS NULL", time.Now().Add(-viper.GetDuration("tracking.inpost_scan_interval"))).
		Find(&inpostCreds).Error; err != nil {
		return fmt.Errorf("failed to find inpost credentials to scan: %w", err)
	}
	for _, c := range inpostCreds {
		err := ts.app.InpostScannerService.ScanUserPackages(c)
		if err != nil {
			log.Printf("failed to scan user packages for phone number %v: %v", c.PhoneNumber, err)
		}
		time.Sleep(time.Second * 2)
	}
	log.Printf("Inpost account scan finished in %v", time.Since(lastCheckStarted))
	return nil
}
