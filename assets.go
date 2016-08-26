package vpp

type PricingParam string

const (
	PricingParamStd  PricingParam = "STDQ" // Standard Quality
	PricingParamPlus              = "PLUS" // High Quality (Books only)
)

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

type VPPAsset struct {
	AdamID          string       `json:"adamIdStr,omitempty"`
	ProductTypeID   int          `json:"productTypeId,omitempty"`
	PricingParam    PricingParam `json:"pricingParam,omitempty"`
	ProductTypeName string       `json:"productTypeName"`
}

type AssetsService interface {
}

type assetsService struct {
	client *vppClient
	sToken string
}

func (as *assetsService) GetAssets() {

}
