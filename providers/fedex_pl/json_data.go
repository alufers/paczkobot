package fedex_pl

// not used in the Polish domestic tracking
type FedexPlTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type FedexPlApiConfigs struct {
	ClientID            string `json:"clientID"`
	ClientSecret        string `json:"clientSecret"`
	GrantType           string `json:"grantType"`
	ApigURL             string `json:"apigURL"`
	RewardsAPIGEnabled  string `json:"rewardsAPIGEnabled"`
	RewardsClientID     string `json:"rewardsClientID"`
	RewardsClientSecret string `json:"rewardsClientSecret"`
}

type FedexPlTrackingResponse struct {
	Level            int         `json:"level"`
	Status           string      `json:"status"`
	Tuid             interface{} `json:"tuid"`
	Type             interface{} `json:"type"`
	TrackingKey      string      `json:"trackingKey"`
	DeliveryDate     string      `json:"deliveryDate"`
	DeliveryDepot    string      `json:"deliveryDepot"`
	ShipmentDate     string      `json:"shipmentDate"`
	ShipmentDepot    string      `json:"shipmentDepot"`
	Weight           float64     `json:"weight"`
	ParcelQty        int         `json:"parcelQty"`
	ShipmentNo       string      `json:"shipmentNo"`
	ProcessStep      float64     `json:"processStep"`
	ProcessType      string      `json:"processType"`
	ProcessStepStyle string      `json:"processStepStyle"`
	Events           []Events    `json:"events"`
	ShipmentExtID    interface{} `json:"shipmentExtId"`
}

type Events struct {
	Depot     string      `json:"depot"`
	Euid      interface{} `json:"euid"`
	EventDate string      `json:"eventDate"`
	EventName string      `json:"eventName"`
	EventExID int         `json:"eventExId"`
	BitMark   interface{} `json:"bitMark"`
	EventCode string      `json:"eventCode"`
	Act       int         `json:"act"`
	ArgsExtra interface{} `json:"argsExtra"`
	Signature interface{} `json:"signature"`
}
