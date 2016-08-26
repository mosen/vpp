package vpp

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"net/url"
)

const (
	serviceConfigPath = "VPPServiceConfigSrv"
)

type ServiceConfig struct {
	InvitationEmailURL               string     `json:"invitationEmailUrl"`
	RegisterUserSrvURL               string     `json:"registerUserSrvUrl"`
	EditUserSrvURL                   string     `json:"editUserSrvUrl"`
	GetUserSrvURL                    string     `json:"getUserSrvUrl"`
	RetireUserSrvURL                 string     `json:"retireUserSrvUrl"`
	GetUsersSrvURL                   string     `json:"getUsersSrvUrl"`
	GetLicensesSrvURL                string     `json:"getLicensesSrvUrl"`
	AssociateLicenseSrvURL           string     `json:"associateLicenseSrvUrl"`
	DisassociateLicenseSrvURL        string     `json:"disassociateLicenseSrvUrl"`
	ClientConfigSrvURL               string     `json:"clientConfigSrvUrl"`
	ErrorCodes                       []VPPError `json:"errorCodes"`
	GetVPPAssetsSrvURL               string     `json:"getVPPAssetsSrvUrl"`
	InvitationEmailUrl               string     `json:"invitationEmailUrl"`
	ManageVPPLicensesByAdamIdSrvURL  string     `json:"manageVPPLicensesByAdamIdSrvUrl"`
	MaxBatchAssociateLicenseCount    int        `json:"maxBatchAssociateLicenseCount"`
	MaxBatchDisassociateLicenseCount int        `json:"maxBatchDisassociateLicenseCount"`
	Status                           int        `json:"status"`
	VPPWebsiteUrl                    string     `json:"vppWebsiteUrl"`
}

// ClientContext represents the information about the current MDM which is stored with the VPP service to ensure
// that two MDM services are not managing the same VPP account.
type ClientContext struct {
	Hostname string `json:"hostname,omitempty"`
	GUID     string `json:"guid,omitempty"`
}

func NewClientContext(hostname string) *ClientContext {
	guid := uuid.NewV4()
	return &ClientContext{
		Hostname: hostname,
		GUID:     guid.String(),
	}
}

type ConfigService interface {
	ServiceConfig() (*ServiceConfig, error)
	ClientContext() (string, error)
	UpdateClientContext(clientContext *ClientContext) (string, error)
}

type configService struct {
	client *vppClient
	sToken string
}

func (s *configService) ServiceConfig() (*ServiceConfig, error) {
	var response ServiceConfig
	path, _ := url.Parse(serviceConfigPath)
	u := s.client.BaseURL.ResolveReference(path)
	req, err := s.client.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	err = s.client.Do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type vppClientConfigSrvRequest struct {
	ClientContext string `json:"clientContext,omitempty"`
	SToken        string `json:"sToken"`
}

type vppClientConfigSrvResponse struct {
	Status        Status `json:"status"`
	ClientContext string `json:"clientContext,omitempty"`
	CountryCode   string `json:"countryCode,omitempty"`
	*VPPError
}

// ClientContext retrieves the current clientContext value by posting an empty request to clientConfigSrvURL.
func (s *configService) ClientContext() (string, error) {
	var response *vppClientConfigSrvResponse
	var request *vppClientConfigSrvRequest = &vppClientConfigSrvRequest{
		SToken: s.sToken,
	}

	req, err := s.client.NewRequest("POST", s.client.Config.serviceConfig.ClientConfigSrvURL, request)
	if err != nil {
		return "", err
	}

	err = s.client.Do(req, &response)
	if err != nil {
		return "", err
	}

	if response.Status == StatusErr {
		return "", response.VPPError
	}

	return response.ClientContext, nil
}

// UpdateClientContext updates the clientContext with the VPP Service.
// The clientContext is used to indicate which product is managing this VPP account so that two products are not
// simultaneously attempting to associate/disassociate licenses.
func (s *configService) UpdateClientContext(clientContext *ClientContext) (string, error) {
	var clientContextBytes []byte
	clientContextBytes, err := json.Marshal(&clientContext)
	if err != nil {
		return "", err
	}

	var response *vppClientConfigSrvResponse
	var request *vppClientConfigSrvRequest = &vppClientConfigSrvRequest{
		ClientContext: string(clientContextBytes),
		SToken:        s.sToken,
	}

	req, err := s.client.NewRequest("POST", s.client.Config.serviceConfig.ClientConfigSrvURL, request)
	if err != nil {
		return "", err
	}

	err = s.client.Do(req, &response)
	if err != nil {
		return "", err
	}

	if response.Status == StatusErr {
		return "", response.VPPError
	}

	return response.CountryCode, nil
}
