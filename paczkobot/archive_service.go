package paczkobot

import (
	"fmt"

	"github.com/alufers/paczkobot/commondata"
)

type ArchiveService struct {
	App *BotApp
}

func NewArchiveService(a *BotApp) *ArchiveService {
	return &ArchiveService{
		App: a,
	}
}

func (s *ArchiveService) FetchAndArchivePackagesForUser(telegramUserID int64) error {
	followedPackages := []FollowedPackageTelegramUser{}

	if err := s.App.DB.Where("telegram_user_id = ? AND archived = ?", telegramUserID, false).
		Preload("FollowedPackage").
		Preload("FollowedPackage.FollowedPackageProviders").
		Find(&followedPackages).Error; err != nil {
		return fmt.Errorf("failed to query DB: %w", err)
	}
	return s.ArchivePackagesIfNeeded(followedPackages)
}

// ArchivePackagesIfNeeded checks whether the passed packages meet the criteria to be archived
// FollowedPackage and FollowedPackage.FollowedPackageProviders must be preloaded
func (s *ArchiveService) ArchivePackagesIfNeeded(packages []FollowedPackageTelegramUser) error {
	for _, p := range packages {
		if s.shouldArchivePackage(p) {
			if err := s.App.DB.Model(&p).Update("archived", true).Error; err != nil {
				return fmt.Errorf("failed to update package: %w", err)
			}
		}
	}
	return nil
}

func (s *ArchiveService) shouldArchivePackage(p FollowedPackageTelegramUser) bool {
	if p.Archived {
		return false
	}
	if p.FollowedPackage.Inactive {
		return true
	}
	for _, provider := range p.FollowedPackage.FollowedPackageProviders {
		for _, s := range commondata.CommonTrackingStepsToArchive {
			if s == provider.LastStatusCommonType {
				return true
			}
		}
	}
	return false
}
