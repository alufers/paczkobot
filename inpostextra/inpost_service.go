package inpostextra

import (
	"context"

	"gorm.io/gorm"
)

type InpostService interface {
	SendSMSCode(ctx context.Context, phoneNumber string) error
	ConfirmSMSCode(ctx context.Context, phoneNumber string, code string) (*InpostCredentials, error)
	Authenticate(ctx context.Context, creds *InpostCredentials) error
	ReauthenticateIfNeeded(ctx context.Context, db *gorm.DB, creds *InpostCredentials) error
	GetParcel(ctx context.Context, db *gorm.DB, creds *InpostCredentials, shipmentNumber string) (*InpostParcel, error)
	GetUserParcels(ctx context.Context, db *gorm.DB, creds *InpostCredentials) (*GetTrackedParcelsResponse, error)
	OpenParcelLocker(ctx context.Context, db *gorm.DB, creds *InpostCredentials, shipmentNumber string) error
}
