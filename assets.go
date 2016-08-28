package vpp

type PricingParam string

const (
	PricingParamStd  PricingParam = "STDQ" // Standard Quality
	PricingParamPlus              = "PLUS" // High Quality (Books only)
)

// VPPAssetAssignment represents a single asset/product and its currently available/assigned licenses totals.
type VPPAssetAssignment struct {
	AdamIdStr        string       `json:"adamIdStr"`
	AssignedCount    int          `json:"assignedCount"`
	AvailableCount   int          `json:"availableCount"`
	DeviceAssignable bool         `json:"deviceAssignable"`
	IsIrrevocable    bool         `json:"isIrrevocable"`
	PricingParam     PricingParam `json:"pricingParam"`
	ProductTypeId    int          `json:"productTypeId"`
	ProductTypeName  string       `json:"productTypeName"`
	RetiredCount     int          `json:"retiredCount"`
	TotalCount       int          `json:"totalCount"`
}

// VPPAsset represents a single licenseable item
type VPPAsset struct {
	AdamID          string       `json:"adamIdStr,omitempty"`
	ProductTypeID   int          `json:"productTypeId,omitempty"`
	PricingParam    PricingParam `json:"pricingParam,omitempty"`
	ProductTypeName string       `json:"productTypeName"`
}

// AssetsService describes an interface that is capable of reporting on VPP assets and their license counts.
type AssetsService interface {
	GetAssets(includeLicenseCounts bool) ([]VPPAssetAssignment, error)
}

type assetsService struct {
	client *vppClient
	sToken string
}

type getVPPAssetsSrvRequest struct {
	IncludeLicenseCounts bool   `json:"includeLicenseCounts,omitempty"`
	SToken               string `json:"sToken"`
}

type getVPPAssetsSrvResponse struct {
	Status Status               `json:"status,omitempty"`
	Assets []VPPAssetAssignment `json:"assets,omitempty"`
	*VPPError
}

// Get a list of VPP assets and their license counts.
// Note that if you use includeLicenseCounts the client may be deemed as overloading the VPP service and the next
// request may be delayed.
func (s *assetsService) GetAssets(includeLicenseCounts bool) ([]VPPAssetAssignment, error) {
	var response *getVPPAssetsSrvResponse
	var request *getVPPAssetsSrvRequest = &getVPPAssetsSrvRequest{
		IncludeLicenseCounts: includeLicenseCounts,
		SToken:               s.sToken,
	}
	req, err := s.client.NewRequest("POST", s.client.Config.serviceConfig.GetVPPAssetsSrvURL, request)
	if err != nil {
		return nil, err
	}

	err = s.client.Do(req, &response)
	if err != nil {
		return nil, err
	}

	if response.Status == StatusErr {
		return nil, response.VPPError
	}

	return response.Assets, nil
}
