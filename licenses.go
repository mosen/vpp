package vpp

type VPPLicense struct {
	LicenseID       string `json:"licenseIdStr,omitempty"`
	AdamID          string `json:"adamId,omitempty"`
	ProductTypeID   int    `json:"productTypeId,omitempty"`
	PricingParam    string `json:"pricingParam,omitempty"`
	ProductTypeName string `json:"productTypeName"`
	IsIrrevocable   bool   `json:"isIrrevocable"`
}

type getVPPLicensesSrvRequest struct {
}

type getVPPLicensesSrvResponse struct {
	TotalCount      int          `json:"totalCount,omitempty"`
	TotalBatchCount int          `json:"totalBatchCount,omitempty"`
	Licenses        []VPPLicense `json:"licenses,omitempty"`
	*VPPError
}

type LicensesService interface {
}

type licensesService struct {
	client *vppClient
}

func (svc *licensesService) GetLicenses() {

}
