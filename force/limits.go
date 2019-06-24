package force

// Limits is a slice of Limit
type Limits map[string]Limit

// Limit defines the limits
type Limit struct {
	Remaining float64
	Max       float64
}

// GetLimits is a func that returns the limits
func (forceAPI *API) GetLimits() (limits *Limits, err error) {
	uri := forceAPI.apiResources[limitsKey]

	limits = &Limits{}
	err = forceAPI.Get(uri, nil, limits)

	return
}
