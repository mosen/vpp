package vpp


type VPPAsset struct {
	AdamIdStr string `json:"adamIdStr"`
	AssignedCount int `json:"assignedCount"`
	AvailableCount int `json:"availableCount"`
	DeviceAssignable bool `json:"deviceAssignable"`
	IsIrrevocable bool `json:"isIrrevocable"`
	PricingParam string `json:"pricingParam"`
	ProductTypeId int `json:"productTypeId"`
	ProductTypeName string `json:"productTypeName"`
	RetiredCount int `json:"retiredCount"`
	TotalCount int `json:"totalCount"`
}

type AssetsService interface {

}

type assetsService struct {
	client *vppClient
}


func (as *assetsService) GetAssets() {

}