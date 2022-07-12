package yuntrack

type YunTrackRequest struct {
	NumberList          []string `json:"NumberList"`
	CaptchaVerification string   `json:"CaptchaVerification"`
	Year                float64  `json:"Year"`
}

type YunTrackResponse struct {
	ResultList []struct {
		ID        string  `json:"Id"`
		Status    float64 `json:"Status"`
		TrackInfo struct {
			WaybillNumber          string  `json:"WaybillNumber"`
			TrackingNumber         string  `json:"TrackingNumber"`
			RelatedNumber          string  `json:"RelatedNumber"`
			CustomerOrderNumber    string  `json:"CustomerOrderNumber"`
			ChannelCodeIn          string  `json:"ChannelCodeIn"`
			ChannelNameIn          string  `json:"ChannelNameIn"`
			ChannelEnNameIn        string  `json:"ChannelEnNameIn"`
			ChannelCodeOut         string  `json:"ChannelCodeOut"`
			ProviderName           string  `json:"ProviderName"`
			ProviderSite           string  `json:"ProviderSite"`
			Telephone              string  `json:"Telephone"`
			AdditionalNotes        string  `json:"AdditionalNotes"`
			CheckInDate            string  `json:"CheckInDate"`
			CheckOutDate           string  `json:"CheckOutDate"`
			PickupDate             string  `json:"PickupDate"`
			ChannelNameOut         string  `json:"ChannelNameOut"`
			ChannelEnNameOut       string  `json:"ChannelEnNameOut"`
			TrackingFlag           float64 `json:"TrackingFlag"`
			CustomerCode           string  `json:"CustomerCode"`
			DestinationCountryCode string  `json:"DestinationCountryCode"`
			OriginCountryCode      string  `json:"OriginCountryCode"`
			PostalCode             string  `json:"PostalCode"`
			Weight                 float64 `json:"Weight"`
			TrackingStatus         float64 `json:"TrackingStatus"`
			IntervalDays           float64 `json:"IntervalDays"`
			IntervalWorkdays       float64 `json:"IntervalWorkdays"`
			TrackEventCount        float64 `json:"TrackEventCount"`
			LastTrackEvent         struct {
				SortCode             float64     `json:"SortCode"`
				IsDeleted            bool        `json:"IsDeleted"`
				ProcessDate          string      `json:"ProcessDate"`
				ProcessDateTimestamp int64       `json:"ProcessDateTimestamp"`
				ProcessContent       string      `json:"ProcessContent"`
				ProcessContentCode   string      `json:"ProcessContentCode"`
				ProcessLocation      string      `json:"ProcessLocation"`
				ProcessCountry       string      `json:"ProcessCountry"`
				ProcessProvince      string      `json:"ProcessProvince"`
				ProcessCity          string      `json:"ProcessCity"`
				TrackingStatus       float64     `json:"TrackingStatus"`
				CreatedOn            string      `json:"CreatedOn"`
				FlowType             float64     `json:"FlowType"`
				AirportCode          interface{} `json:"AirportCode"`
				ReturnReason         interface{} `json:"ReturnReason"`
				TrackCode            string      `json:"TrackCode"`
				TrackCodeName        string      `json:"TrackCodeName"`
			} `json:"LastTrackEvent"`
			CreatedOn         string      `json:"CreatedOn"`
			UpdatedOn         string      `json:"UpdatedOn"`
			Remarks           interface{} `json:"Remarks"`
			TrackEventDetails []struct {
				ProcessLocation interface{} `json:"ProcessLocation"`
				CreatedOn       string      `json:"CreatedOn"`
				ProcessContent  string      `json:"ProcessContent"`
			} `json:"TrackEventDetails"`
			SystemCode         float64     `json:"SystemCode"`
			ArrivedOn          string      `json:"ArrivedOn"`
			TotalPackageNumber string      `json:"TotalPackageNumber"`
			IsVirtual          bool        `json:"IsVirtual"`
			QueryMode          float64     `json:"QueryMode"`
			DisplayMode        float64     `json:"DisplayMode"`
			BillingWeight      float64     `json:"BillingWeight"`
			Amount             float64     `json:"Amount"`
			SellerCode         string      `json:"SellerCode"`
			AirWaybillNumber   interface{} `json:"AirWaybillNumber"`
			OrderSource        string      `json:"OrderSource"`
			PackageNumber      string      `json:"PackageNumber"`
			IsServerPickup     bool        `json:"IsServerPickup"`
		} `json:"TrackInfo"`
		TrackData struct {
			ChildCount       float64 `json:"ChildCount"`
			TrackStatus      string  `json:"TrackStatus"`
			IsAlert          bool    `json:"IsAlert"`
			DetailingID      string  `json:"DetailingId"`
			ProcessGroupList []struct {
				ProcessGroupDate  string `json:"ProcessGroupDate"`
				ProcessDetailList []struct {
					ProcessDate    string `json:"ProcessDate"`
					ProcessContent string `json:"ProcessContent"`
				} `json:"ProcessDetailList"`
			} `json:"ProcessGroupList"`
		} `json:"TrackData"`
	} `json:"ResultList"`
}
