package vpp

import "net/url"

type ServiceConfig struct {
	InvitationEmailURL url.URL `json:"invitationEmailUrl"`
	RegisterUserSrvURL url.URL `json:"registerUserSrvUrl"`
	EditUserSrvURL url.URL `json:"editUserSrvUrl"`
	GetUserSrvURL url.URL `json:"getUserSrvUrl"`
	RetireUserSrvURL url.URL `json:"retireUserSrvUrl"`
	GetUsersSrvURL url.URL `json:"getUsersSrvUrl"`
	GetLicensesSrvURL url.URL `json:"getLicensesSrvUrl"`
	AssociateLicenseSrvURL url.URL `json:"associateLicenseSrvUrl"`
	DisassociateLicenseSrvURL url.URL `json:"disassociateLicenseSrvUrl"`
}


type ConfigService interface {
	ServiceConfig() (*ServiceConfig, error)
}

type configService struct {
	client *vppClient
}

func (cs *configService) ServiceConfig() (*ServiceConfig, error) {

}