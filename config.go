package vpp

import "net/url"

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

type ConfigService interface {
	ServiceConfig() (*ServiceConfig, error)
}

type configService struct {
	client *vppClient
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
