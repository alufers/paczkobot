package inpostextra

import "gorm.io/gorm"

type InpostService interface {
	SendSMSCode(phoneNumber string) error
	ConfirmSMSCode(phoneNumber string, code string) (*InpostCredentials, error)
	Authenticate(creds *InpostCredentials) error
	ReauthenticateIfNeeded(db *gorm.DB, creds *InpostCredentials) error
	GetParcel(db *gorm.DB, creds *InpostCredentials, shipmentNumber string) (*InpostParcel, error)
	GetUserParcels(db *gorm.DB, creds *InpostCredentials) (*GetTrackedParcelsResponse, error)
	OpenParcelLocker(db *gorm.DB, creds *InpostCredentials, shipmentNumber string) error
}
