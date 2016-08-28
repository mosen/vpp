package vpp

import "errors"

// VPPLicense describes a licensed product (VPPAsset) and its (optional) association with a VPP User.
type VPPLicense struct {
	LicenseID     string `json:"licenseIdStr,omitempty"`
	IsIrrevocable bool   `json:"isIrrevocable"`
	*VPPAsset
	*VPPUser
}

// LicenseAssociation describes an association between a (VPP user OR device serial) and a license
type LicenseAssociation struct {
	ClientUserIDStr string `json:"clientUserIdStr,omitempty"`
	LicenseIDStr    string `json:"licenseIdStr,omitempty"`
	SerialNumber    string `json:"serialNumber,omitempty"`
	*VPPError
}

type getVPPLicensesSrvRequest struct {
	*getLicensesRequestOpts
	*BatchRequest
	SToken string `json:"sToken"`
}

type getVPPLicensesSrvResponse struct {
	Status          Status       `json:"status,omitempty"`
	TotalCount      int          `json:"totalCount,omitempty"`
	TotalBatchCount int          `json:"totalBatchCount,omitempty"`
	Licenses        []VPPLicense `json:"licenses,omitempty"`
	*VPPError
}

// LicensesService describes an interface that can manage VPP licenses
type LicensesService interface {
	GetLicenses(opts ...GetLicensesOption) ([]VPPLicense, error)
	AssociateLicense(user *VPPUser, license *VPPLicense) error
	DisassociateLicense(user *VPPUser, license *VPPLicense) error
}

type licensesService struct {
	client *vppClient
	sToken string
}

type getLicensesRequestOpts struct {
	AssignedOnly bool         `json:"assignedOnly,omitempty"`
	AdamID       int          `json:"adamId,omitempty"`
	PricingParam PricingParam `json:"pricingParam,omitempty"`
}

// GetLicensesOption describes the signature of the closure returned by a function adding an argument to GetLicenses
type GetLicensesOption func(*getLicensesRequestOpts) error

// GetLicenses retrieves a list of available VPP licenses. The result can optionally be filtered by the application id
// and/or its assigned status.
func (s *licensesService) GetLicenses(batch *BatchRequest, opts ...GetLicensesOption) ([]VPPLicense, error) {
	requestOpts := &getLicensesRequestOpts{}
	for _, option := range opts {
		if err := option(requestOpts); err != nil {
			return nil, err
		}
	}
	var request *getVPPLicensesSrvRequest = &getVPPLicensesSrvRequest{
		getLicensesRequestOpts: requestOpts,
		SToken:                 s.sToken,
		BatchRequest:           batch,
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
	AdamIDStr    string       `json:"adamIdStr"`
	PricingParam PricingParam `json:"pricingParam"`

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
	PricingParam    PricingParam         `json:"pricingParam"`
	ProductTypeName string               `json:"productTypeName"`
	IsIrrevocable   bool                 `json:"isIrrevocable"`
	Associations    []LicenseAssociation `json:"associations,omitempty"`
	Disassociations []LicenseAssociation `json:"disassociations,omitempty"`
}

// licenseOperations describes a list of operations on a single asset (app or book).
type licenseOperations struct {
	Asset *VPPAsset

	AssociateUsers         []VPPUser
	AssociateSerialNumbers []string

	DisassociateUsers         []VPPUser
	DisassociateSerialNumbers []string
	DisassociateLicenseIDs    []string
}

// NewLicenseOperations creates a new batch of license operations (assigning/freeing)
func (op *licenseOperations) NewLicenseOperations(asset *VPPAsset) licenseOperations {
	return &licenseOperations{
		Asset: asset,
	}
}

// AssignUser adds a VPP user to the list of license holders for the current asset.
func (op *licenseOperations) AssignUser(user *VPPUser) error {
	if len(op.AssociateSerialNumbers) > 0 {
		return errors.New("you cannot assign licenses to both users and devices in the same license operation")
	}

	op.AssociateUsers = append(op.AssociateUsers, user)
	return nil
}

// AssignSerialNumber adds a device to the list of license holders for the current asset.
func (op *licenseOperations) AssignSerialNumber(serialNumber string) error {
	if len(op.AssociateUsers) > 0 {
		return errors.New("you cannot assign licenses to both users and devices in the same license operation")
	}

	op.AssociateSerialNumbers = append(op.AssociateSerialNumbers, serialNumber)
	return nil
}

// UnassignUser removes a VPP user from the list of license holders for the current asset.
func (op *licenseOperations) UnassignUser(user *VPPUser) error {
	if len(op.DisassociateSerialNumbers) > 0|len(op.DisassociateLicenseIDs) > 0 {
		return errors.New("you can only unassign licenses from either users, serial numbers, or license ids in a single operation")
	}

	op.DisassociateUsers = append(op.DisassociateUsers, user)
	return nil
}

// UnassignSerialNumber removes a device from the list of license holders for the current asset.
func (op *licenseOperations) UnassignSerialNumber(serialNumber string) error {
	if len(op.DisassociateUsers) > 0|len(op.DisassociateLicenseIDs) > 0 {
		return errors.New("you can only unassign licenses from either users, serial numbers, or license ids in a single operation")
	}

	op.DisassociateSerialNumbers = append(op.DisassociateSerialNumbers, serialNumber)
	return nil
}

// UnassignLicenseID removes a VPP license by its identifier.
func (op *licenseOperations) UnassignLicenseID(licenseID string) error {
	if len(op.DisassociateUsers) > 0|len(op.DisassociateSerialNumbers) > 0 {
		return errors.New("you can only unassign licenses from either users, serial numbers, or license ids in a single operation")
	}

	op.DisassociateLicenseIDs = append(op.DisassociateLicenseIDs, licenseID)
	return nil
}

// ManageLicenses is used for the bulk addition/removal of VPP licenses.
// The notify argument controls whether the user will be notified of the removal.
//func (s *licensesService) ManageLicenses(operations *licenseOperations, notify bool) ([]LicenseAssociation, error) {
//
//}

type associateVPPLicenseWithVPPUserSrvRequest struct {
	UserID       int          `json:"userId,omitempty"`
	ClientUserID string       `json:"clientUserIdStr,omitempty"`
	AdamID       string       `json:"adamId,omitempty"`
	LicenseID    string       `json:"licenseId,omitempty"`
	PricingParam PricingParam `json:"pricingParam,omitempty"`
	SToken       string       `json:"sToken"`
}

type associateVPPLicenseWithVPPUserSrvResponse struct {
	Status  Status      `json:"status"`
	License *VPPLicense `json:"license,omitempty"` // TODO: user also included in this object, should we parse it?
	User    *VPPUser    `json:"user,omitempty"`
	*VPPError
}

// DEPRECATED: Associate an (available) VPP license with a MDM system user.
func (s *licensesService) AssociateLicense(user *VPPUser, license *VPPLicense) error {
	var response *associateVPPLicenseWithVPPUserSrvResponse
	var request *associateVPPLicenseWithVPPUserSrvRequest = &associateVPPLicenseWithVPPUserSrvRequest{
		SToken: s.sToken,
	}

	// UserID takes precedence over ClientUserID
	if user.UserID != 0 {
		request.UserID = user.UserID
	} else {
		request.ClientUserID = user.ClientUserIdStr
	}

	// LicenseID takes precedence over AdamID
	if license.LicenseID != "" {
		request.LicenseID = license.LicenseID
	} else {
		request.AdamID = license.AdamID
	}

	req, err := s.client.NewRequest("POST", s.client.Config.serviceConfig.AssociateLicenseSrvURL, request)
	if err != nil {
		return err
	}

	err = s.client.Do(req, &response)
	if err != nil {
		return err
	}

	if response.Status == StatusErr {
		return response.VPPError
	}

	user.Status = response.User.Status
	return nil
}

type disassociateVPPLicenseFromVPPUserSrvRequest struct {
	UserID    int    `json:"userId"`
	LicenseID string `json:"licenseId"`
	SToken    string `json:"sToken"`
}

type disassociateVPPLicenseFromVPPUserSrvResponse struct {
	Status  Status      `json:"status"`
	License *VPPLicense `json:"license,omitempty"`
	User    *VPPUser    `json:"user,omitempty"`
	*VPPError
}

// DEPRECATED: Disassociate a license from a VPP user.
func (s *licensesService) DisassociateLicense(user *VPPUser, license *VPPLicense) error {
	var response *disassociateVPPLicenseFromVPPUserSrvResponse
	var request *disassociateVPPLicenseFromVPPUserSrvRequest = &disassociateVPPLicenseFromVPPUserSrvRequest{
		UserID:    user.UserID,
		LicenseID: license.LicenseID,
		SToken:    s.sToken,
	}

	req, err := s.client.NewRequest("POST", s.client.Config.serviceConfig.DisassociateLicenseSrvURL, request)
	if err != nil {
		return err
	}

	err = s.client.Do(req, &response)
	if err != nil {
		return err
	}

	if response.Status == StatusErr {
		return response.VPPError
	}

	user.Status = response.User.Status
	return nil
}
