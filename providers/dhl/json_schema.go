package dhl

type DHLResponse struct {
	Shipments []*DHLShipment `json:"shipments"`
}

type DHLShipment struct {
	Destination *DHLLocation `json:"destination"`
	Origin      *DHLLocation `json:"origin"`
	Events      []*DHLEvent  `json:"events"`
}

type DHLEvent struct {
	Description string       `json:"description"`
	Location    *DHLLocation `json:"location"`
	StatusCode  string       `json:"statusCode"`
	Timestamp   string       `json:"timestamp"`
}

type DHLLocation struct {
	Address *DHLAddress `json:"address"`
}

func (loc *DHLLocation) String() string {
	if loc.Address == nil {
		return ""
	}
	return loc.Address.String()
}

type DHLAddress struct {
	AddressLocality string `json:"addressLocality"`
}

func (addr *DHLAddress) String() string {
	if addr == nil {
		return ""
	}
	return addr.AddressLocality
}
