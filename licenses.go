package vpp

type VPPLicense struct {
	LicenseID       string `json:"licenseIdStr,omitempty"`
	AdamID          string `json:"adamId,omitempty"`
	ProductTypeID   int    `json:"productTypeId,omitempty"`
	PricingParam    string `json:"pricingParam,omitempty"`
	ProductTypeName string `json:"productTypeName"`
	IsIrrevocable   bool   `json:"isIrrevocable"`
}

type LicenseAssociation struct {
	ClientUserIDStr string `json:"clientUserIdStr,omitempty"`
	LicenseIDStr    string `json:"licenseIdStr,omitempty"`
	SerialNumber    string `json:"serialNumber,omitempty"`
	*VPPError
}

type getVPPLicensesSrvRequest struct {
	*getLicensesRequestOpts
	SToken string `json:"sToken"`
}

type getVPPLicensesSrvResponse struct {
	Status          Status       `json:"status,omitempty"`
	TotalCount      int          `json:"totalCount,omitempty"`
	TotalBatchCount int          `json:"totalBatchCount,omitempty"`
	Licenses        []VPPLicense `json:"licenses,omitempty"`
	*VPPError
}

type LicensesService interface {
	GetLicenses(opts ...GetLicensesOption) ([]VPPLicense, error)
}

type licensesService struct {
	client *vppClient
}

type getLicensesRequestOpts struct {
	*BatchRequestOpts
	AssignedOnly bool   `json:"assignedOnly,omitempty"`
	AdamID       int    `json:"adamId,omitempty"`
	PricingParam string `json:"pricingParam,omitempty"`
}

// GetLicensesOption describes the signature of the closure returned by a function adding an argument to GetLicenses
type GetLicensesOption func(*getLicensesRequestOpts) error

// GetLicenses retrieves a list of available VPP licenses. The result can optionally be filtered by the application id
// and/or its assigned status.
func (s *licensesService) GetLicenses(opts ...GetLicensesOption) ([]VPPLicense, error) {
	sToken, err := s.client.Config.SToken.Base64String()
	if err != nil {
		return nil, err
	}
	requestOpts := &getLicensesRequestOpts{}
	for _, option := range opts {
		if err := option(requestOpts); err != nil {
			return nil, err
		}
	}
	var request *getVPPLicensesSrvRequest = &getVPPLicensesSrvRequest{
		getLicensesRequestOpts: requestOpts,
		SToken:                 sToken,
	}
	var response getVPPLicensesSrvResponse
	req, err := s.client.NewRequest("POST", s.client.Config.serviceConfig.GetLicensesSrvURL, request)
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
	return response.Licenses, nil
}

type manageVPPLicensesByAdamIdSrvRequest struct {
	AdamIDStr    string `json:"adamIdStr"`
	PricingParam string `json:"pricingParam"`

	// Only one of the below is required
	AssociateClientIDStrs  []string `json:"associateClientIdStrs,omitempty"`
	AssociateSerialNumbers []string `json:"associateSerialNumbers,omitempty"`

	// Only one of the below is required
	DisassociateClientUserIDStrs []string `json:"disassociateClientIdStrs,omitempty"`
	DisassociateLicenseIDStrs    []string `json:"disassociateLicenseIdStrs,omitempty"`
	DisassociateSerialNumbers    []string `json:"disassociateSerialNumbers,omitempty"`

	NotifyDisassociation bool   `json:"notifyDisassociation,omitempty"`
	SToken               string `json:"sToken"`
}

type manageVPPLicensesByAdamIdSrvResponse struct {
	Status          Status               `json:"status"`
	AdamIDStr       string               `json:"adamIdStr"`
	ProductTypeID   int                  `json:"productTypeId"`
	PricingParam    string               `json:"pricingParam"`
	ProductTypeName string               `json:"productTypeName"`
	IsIrrevocable   bool                 `json:"isIrrevocable"`
	Associations    []LicenseAssociation `json:"associations,omitempty"`
	Disassociations []LicenseAssociation `json:"disassociations,omitempty"`
}

// ManageLicenses is used for the bulk addition/removal of VPP licenses.
//func (s *licensesService) ManageLicenses() error {
//
//}
