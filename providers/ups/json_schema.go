package ups

type UPSJsonSchema struct {
	StatusCode      string          `json:"statusCode"`
	StatusText      string          `json:"statusText"`
	IsLoggedInUser  bool            `json:"isLoggedInUser"`
	TrackedDateTime string          `json:"trackedDateTime"`
	IsBcdnMultiView bool            `json:"isBcdnMultiView"`
	ReturnToDetails ReturnToDetails `json:"returnToDetails"`
	TrackDetails    []TrackDetails  `json:"trackDetails"`
}
type ReturnToDetails struct {
	ReturnToURL interface{} `json:"returnToURL"`
	ReturnToApp interface{} `json:"returnToApp"`
}
type ScheduledDeliverDateDetail struct {
	MonthCMSKey string `json:"monthCMSKey"`
	DayNum      string `json:"dayNum"`
}
type ShipToAddress struct {
	StreetAddress1     string `json:"streetAddress1"`
	StreetAddress2     string `json:"streetAddress2"`
	StreetAddress3     string `json:"streetAddress3"`
	City               string `json:"city"`
	State              string `json:"state"`
	Province           string `json:"province"`
	Country            string `json:"country"`
	ZipCode            string `json:"zipCode"`
	CompanyName        string `json:"companyName"`
	AttentionName      string `json:"attentionName"`
	IsAddressCorrected bool   `json:"isAddressCorrected"`
	IsReturnAddress    bool   `json:"isReturnAddress"`
	IsHoldAddress      bool   `json:"isHoldAddress"`
}
type ServiceInformation struct {
	ServiceName      string      `json:"serviceName"`
	ServiceLink      string      `json:"serviceLink"`
	ServiceAttribute interface{} `json:"serviceAttribute"`
}
type AdditionalInformation struct {
	ServiceInformation        ServiceInformation `json:"serviceInformation"`
	Weight                    string             `json:"weight"`
	WeightUnit                string             `json:"weightUnit"`
	CodInformation            interface{}        `json:"codInformation"`
	ShippedOrBilledDate       string             `json:"shippedOrBilledDate"`
	ReferenceNumbers          interface{}        `json:"referenceNumbers"`
	PostalServiceTrackingID   string             `json:"postalServiceTrackingID"`
	AlternateTrackingNumbers  interface{}        `json:"alternateTrackingNumbers"`
	OtherRequestedServices    []string           `json:"otherRequestedServices"`
	DescriptionOfGood         string             `json:"descriptionOfGood"`
	CargoReady                string             `json:"cargoReady"`
	FileNumber                string             `json:"fileNumber"`
	OriginPort                string             `json:"originPort"`
	DestinationPort           string             `json:"destinationPort"`
	EstimatedArrival          string             `json:"estimatedArrival"`
	EstimatedDeparture        string             `json:"estimatedDeparture"`
	PoNumber                  string             `json:"poNumber"`
	BlNumber                  string             `json:"blNumber"`
	AppointmentMade           interface{}        `json:"appointmentMade"`
	AppointmentRequested      interface{}        `json:"appointmentRequested"`
	AppointmentRequestedRange interface{}        `json:"appointmentRequestedRange"`
	Manifest                  string             `json:"manifest"`
	IsSmallPackage            bool               `json:"isSmallPackage"`
	ShipmentVolume            interface{}        `json:"shipmentVolume"`
	NumberOfPallets           interface{}        `json:"numberOfPallets"`
	ShipmentCategory          string             `json:"shipmentCategory"`
	PkgSequenceNum            interface{}        `json:"pkgSequenceNum"`
	PkgIdentificationCode     interface{}        `json:"pkgIdentificationCode"`
	PkgID                     interface{}        `json:"pkgID"`
	Product                   interface{}        `json:"product"`
	NumberOfPieces            interface{}        `json:"numberOfPieces"`
	Wwef                      bool               `json:"wwef"`
	WwePostal                 bool               `json:"wwePostal"`
	ShowAltTrkLink            bool               `json:"showAltTrkLink"`
	WweParcel                 bool               `json:"wweParcel"`
	DeliveryPreference        interface{}        `json:"deliveryPreference"`
	LiftGateOnDelivery        interface{}        `json:"liftGateOnDelivery"`
	LiftGateOnPickup          interface{}        `json:"liftGateOnPickup"`
	PickupDropOffDate         interface{}        `json:"pickupDropOffDate"`
	PickupPreference          interface{}        `json:"pickupPreference"`
}
type AttentionNeeded struct {
	Actions            []interface{} `json:"actions"`
	IsCorrectMyAddress bool          `json:"isCorrectMyAddress"`
}
type Milestone struct {
	Name              string      `json:"name"`
	IsCurrent         bool        `json:"isCurrent"`
	IsCompleted       bool        `json:"isCompleted"`
	SupplementaryText interface{} `json:"supplementaryText"`
	IsRFIDIcon        bool        `json:"isRFIDIcon"`
}
type ShipmentProgressActivities struct {
	Date                          string      `json:"date"`
	Time                          string      `json:"time"`
	Location                      string      `json:"location"`
	ActivityScan                  interface{} `json:"activityScan"`
	Milestone                     Milestone   `json:"milestone"`
	IsInOverViewTable             bool        `json:"isInOverViewTable"`
	ActivityAdditionalDescription interface{} `json:"activityAdditionalDescription"`
	Trailer                       interface{} `json:"trailer"`
	IsDisplayPodLink              bool        `json:"isDisplayPodLink"`
	IsRFIDIconEvent               bool        `json:"isRFIDIconEvent"`
	ActCode                       interface{} `json:"actCode"`
}
type NotLoggedInMyChoicePage struct {
	DeliveryChangesOptions []string `json:"deliveryChangesOptions"`
	IsDriverInstructions   bool     `json:"isDriverInstructions"`
	IsDeliveryChanges      bool     `json:"isDeliveryChanges"`
	IsSignUpLogin          bool     `json:"isSignUpLogin"`
	IsInfonoticeNote       bool     `json:"isInfonoticeNote"`
}
type DeliveryOptions struct {
	IsNotLoggedInMyChoicePage  bool                    `json:"isNotLoggedInMyChoicePage"`
	IsInfoNoticePage           bool                    `json:"isInfoNoticePage"`
	IsLoggedInNoneMyChoicePage bool                    `json:"isLoggedInNoneMyChoicePage"`
	IsContactOnlyPage          bool                    `json:"isContactOnlyPage"`
	IsRedirect                 bool                    `json:"isRedirect"`
	Login                      string                  `json:"login"`
	SignUp                     string                  `json:"signUp"`
	DeliveryChanges            interface{}             `json:"deliveryChanges"`
	SiiEligible                interface{}             `json:"siiEligible"`
	UpgradeToUpsGround         interface{}             `json:"upgradeToUpsGround"`
	ContactUps                 interface{}             `json:"contactUps"`
	Redirect                   interface{}             `json:"redirect"`
	MyChoiceTandCURL           interface{}             `json:"myChoiceTandCUrl"`
	DoappURL                   string                  `json:"doappUrl"`
	DcrEligible                bool                    `json:"dcrEligible"`
	NotLoggedInMyChoicePage    NotLoggedInMyChoicePage `json:"notLoggedInMyChoicePage"`
	IsIdentityVerification     bool                    `json:"isIdentityVerification"`
	Text                       interface{}             `json:"text"`
	Name                       string                  `json:"name"`
	URL                        interface{}             `json:"url"`
}
type UserOptions struct {
	DeliveryOptions DeliveryOptions `json:"deliveryOptions"`
}
type LanguageOptions struct {
	Locale   string `json:"locale"`
	Language string `json:"language"`
}
type SendUpdatesOptions struct {
	BridgePageType                             string            `json:"bridgePageType"`
	ModalType                                  string            `json:"modalType"`
	MyChoicePreferencesLink                    string            `json:"myChoicePreferencesLink"`
	IsDisplayCurrentStatus                     bool              `json:"isDisplayCurrentStatus"`
	IsDisplayUnforeseenEventsOrDelays          bool              `json:"isDisplayUnforeseenEventsOrDelays"`
	IsDisplayShipmentDelivered                 bool              `json:"isDisplayShipmentDelivered"`
	IsDisplayPackageReadyForPickup             bool              `json:"isDisplayPackageReadyForPickup"`
	IsPreCheckedCurrentStatus                  bool              `json:"isPreCheckedCurrentStatus"`
	IsPreCheckedUnforeseenEventsOrDelays       bool              `json:"isPreCheckedUnforeseenEventsOrDelays"`
	IsPreCheckedShipmentDelivered              bool              `json:"isPreCheckedShipmentDelivered"`
	IsPreCheckedPackageReadyForPickup          bool              `json:"isPreCheckedPackageReadyForPickup"`
	IsDayBeforeDeliveryAlert                   bool              `json:"isDayBeforeDeliveryAlert"`
	IsDayOfDeliveryAlert                       bool              `json:"isDayOfDeliveryAlert"`
	IsDeliveryConfirmationAlert                bool              `json:"isDeliveryConfirmationAlert"`
	IsPackageAvailableForPickupAlert           bool              `json:"isPackageAvailableForPickupAlert"`
	IsDeliveryScheduleUpdateAlert              bool              `json:"isDeliveryScheduleUpdateAlert"`
	IsPreCheckedDayBeforeDeliveryAlert         bool              `json:"isPreCheckedDayBeforeDeliveryAlert"`
	IsPreCheckedDayOfDeliveryAlert             bool              `json:"isPreCheckedDayOfDeliveryAlert"`
	IsPreCheckedDeliveryConfirmationAlert      bool              `json:"isPreCheckedDeliveryConfirmationAlert"`
	IsPreCheckedPackageAvailableForPickupAlert bool              `json:"isPreCheckedPackageAvailableForPickupAlert"`
	IsPreCheckedDeliveryScheduleUpdateAlert    bool              `json:"isPreCheckedDeliveryScheduleUpdateAlert"`
	LanguageOptions                            []LanguageOptions `json:"languageOptions"`
	DefaultSelectedLanguage                    string            `json:"defaultSelectedLanguage"`
	Email                                      interface{}       `json:"email"`
	PhoneNumber                                interface{}       `json:"phoneNumber"`
	IsSMSSupported                             bool              `json:"isSMSSupported"`
	Recipients                                 interface{}       `json:"recipients"`
	Text                                       interface{}       `json:"text"`
	Name                                       string            `json:"name"`
	URL                                        interface{}       `json:"url"`
}
type Promo struct {
	IsPackagePromotion bool        `json:"isPackagePromotion"`
	IsShipperPromotion bool        `json:"isShipperPromotion"`
	ProductImage       interface{} `json:"productImage"`
	Title              interface{} `json:"title"`
	Description        interface{} `json:"description"`
	ShipperURL         interface{} `json:"shipperURL"`
	ShipperLogoURL     interface{} `json:"shipperLogoURL"`
}
type AsrInformation struct {
	AllowDriverRelease interface{} `json:"allowDriverRelease"`
	ProcessEDN         interface{} `json:"processEDN"`
}
type TrackDetails struct {
	ErrorCode                     interface{}                  `json:"errorCode"`
	ErrorText                     interface{}                  `json:"errorText"`
	RequestedTrackingNumber       string                       `json:"requestedTrackingNumber"`
	TrackingNumber                string                       `json:"trackingNumber"`
	IsMobileDevice                bool                         `json:"isMobileDevice"`
	PackageStatus                 string                       `json:"packageStatus"`
	PackageStatusType             string                       `json:"packageStatusType"`
	PackageStatusCode             string                       `json:"packageStatusCode"`
	ProgressBarType               string                       `json:"progressBarType"`
	ProgressBarPercentage         string                       `json:"progressBarPercentage"`
	SimplifiedText                string                       `json:"simplifiedText"`
	ScheduledDeliveryDayCMSKey    string                       `json:"scheduledDeliveryDayCMSKey"`
	ScheduledDeliveryDate         string                       `json:"scheduledDeliveryDate"`
	ScheduledDeliverDateDetail    ScheduledDeliverDateDetail   `json:"scheduledDeliverDateDetail"`
	NoEstimatedDeliveryDateLabel  interface{}                  `json:"noEstimatedDeliveryDateLabel"`
	ScheduledDeliveryTime         string                       `json:"scheduledDeliveryTime"`
	ScheduledDeliveryTimeEODLabel string                       `json:"scheduledDeliveryTimeEODLabel"`
	PackageCommitedTime           string                       `json:"packageCommitedTime"`
	EndOfDayResCMSKey             interface{}                  `json:"endOfDayResCMSKey"`
	DeliveredDayCMSKey            string                       `json:"deliveredDayCMSKey"`
	DeliveredDate                 string                       `json:"deliveredDate"`
	DeliveredDateDetail           interface{}                  `json:"deliveredDateDetail"`
	DeliveredTime                 string                       `json:"deliveredTime"`
	ReceivedBy                    string                       `json:"receivedBy"`
	LeaveAt                       interface{}                  `json:"leaveAt"`
	AlertCount                    int                          `json:"alertCount"`
	IsEligibleViewMoreAlerts      bool                         `json:"isEligibleViewMoreAlerts"`
	CdiLeaveAt                    interface{}                  `json:"cdiLeaveAt"`
	LeftAt                        string                       `json:"leftAt"`
	ShipToAddress                 ShipToAddress                `json:"shipToAddress"`
	ShipFromAddress               interface{}                  `json:"shipFromAddress"`
	ConsigneeAddress              interface{}                  `json:"consigneeAddress"`
	SignatureTrackingURL          interface{}                  `json:"signatureTrackingUrl"`
	TrackHistoryDescription       interface{}                  `json:"trackHistoryDescription"`
	AdditionalInformation         AdditionalInformation        `json:"additionalInformation"`
	SpecialInstructions           interface{}                  `json:"specialInstructions"`
	ProofOfDeliveryURL            interface{}                  `json:"proofOfDeliveryUrl"`
	UpsAccessPoint                interface{}                  `json:"upsAccessPoint"`
	AdditionalPackagesCount       interface{}                  `json:"additionalPackagesCount"`
	AttentionNeeded               AttentionNeeded              `json:"attentionNeeded"`
	ShipmentProgressActivities    []ShipmentProgressActivities `json:"shipmentProgressActivities"`
	TrackingNumberType            string                       `json:"trackingNumberType"`
	PreAuthorizedForReturnData    interface{}                  `json:"preAuthorizedForReturnData"`
	ShipToAddressLblKey           string                       `json:"shipToAddressLblKey"`
	TrackSummaryView              interface{}                  `json:"trackSummaryView"`
	SenderShipperNumber           string                       `json:"senderShipperNumber"`
	InternalKey                   string                       `json:"internalKey"`
	UserOptions                   UserOptions                  `json:"userOptions"`
	SendUpdatesOptions            SendUpdatesOptions           `json:"sendUpdatesOptions"`
	MyChoiceUpSellLink            string                       `json:"myChoiceUpSellLink"`
	BcdnNumber                    interface{}                  `json:"bcdnNumber"`
	Promo                         Promo                        `json:"promo"`
	WhatsNextText                 interface{}                  `json:"whatsNextText"`
	PackageStatusTimeLbl          string                       `json:"packageStatusTimeLbl"`
	DeSepcialTranslation          bool                         `json:"deSepcialTranslation"`
	PackageStatusTime             string                       `json:"packageStatusTime"`
	MyChoiceToken                 interface{}                  `json:"myChoiceToken"`
	ShowMycTerms                  bool                         `json:"showMycTerms"`
	EnrollNum                     string                       `json:"enrollNum"`
	ShowConfirmWindow             bool                         `json:"showConfirmWindow"`
	ConfirmWindowLbl              interface{}                  `json:"confirmWindowLbl"`
	ConfirmWindowLink             interface{}                  `json:"confirmWindowLink"`
	FollowMyDelivery              interface{}                  `json:"followMyDelivery"`
	FileClaim                     interface{}                  `json:"fileClaim"`
	ViewClaim                     interface{}                  `json:"viewClaim"`
	FlightInformation             interface{}                  `json:"flightInformation"`
	VoyageInformation             interface{}                  `json:"voyageInformation"`
	ViewDeliveryReceipt           interface{}                  `json:"viewDeliveryReceipt"`
	IsInWatchList                 bool                         `json:"isInWatchList"`
	IsHistoryUpdateRequire        bool                         `json:"isHistoryUpdateRequire"`
	ConsumerHub                   string                       `json:"consumerHub"`
	CampusShip                    interface{}                  `json:"campusShip"`
	AsrInformation                AsrInformation               `json:"asrInformation"`
	IsSuppressDetailTab           bool                         `json:"isSuppressDetailTab"`
	IsUpsPremierPackage           bool                         `json:"isUpsPremierPackage"`
	LastSensorLocation            interface{}                  `json:"lastSensorLocation"`
	IsPremierStyleEligible        bool                         `json:"isPremierStyleEligible"`
	IsEDW                         bool                         `json:"isEDW"`
}
