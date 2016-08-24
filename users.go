package vpp

import (
	"github.com/satori/go.uuid"
)

const (
	RegStatusRegistered string = "Registered"
	RegStatusAssociated        = "Associated"
	RegStatusRetired           = "Retired"
)

type VPPUser struct {
	UserID          int    `json:"userId,omitempty"`
	Email           string `json:"email,omitempty"`
	Status          string `json:"status,omitempty"`
	InviteURL       string `json:"inviteUrl,omitempty"`
	InviteCode      string `json:"inviteCode,omitempty"`
	ClientUserIdStr string `json:"clientUserIdStr,omitempty"`
}

func NewUser(email string, id string) *VPPUser {
	if id == "" {
		id = uuid.NewV4().String()
	}
	return &VPPUser{
		ClientUserIdStr: id,
		Email:           email,
	}
}

type UsersService interface {
	RegisterUser(user *VPPUser) (*VPPUser, error)
	GetUser(*VPPUser) error
}

type usersService struct {
	client *vppClient
}

type registerVPPUserSrvRequest struct {
	*VPPUser
	SToken string `json:"sToken"`
}

type registerVPPUserSrvResponse struct {
	Status int      `json:"status,omitempty"`
	User   *VPPUser `json:"user,omitempty"`
	*VPPError
}

func (s *usersService) RegisterUser(user *VPPUser) (*VPPUser, error) {
	sToken, err := s.client.Config.SToken.Base64String()
	if err != nil {
		return nil, err
	}
	var response *registerVPPUserSrvResponse
	var request *registerVPPUserSrvRequest = &registerVPPUserSrvRequest{
		VPPUser: user,
		SToken:  sToken,
	}

	req, err := s.client.NewRequest("POST", s.client.Config.serviceConfig.RegisterUserSrvURL, request)
	if err != nil {
		return nil, err
	}

	err = s.client.Do(req, &response)
	if err != nil {
		return nil, err
	}

	if response.Status == -1 {
		return nil, response.VPPError
	}

	return response.User, nil
}

type getVPPUserSrvRequest struct {
	UserID          int    `json:"userId,omitempty"`
	ClientUserIdStr string `json:"clientUserIdStr,omitempty"`
	ITSIdHash       string `json:"itsIdHash,omitempty"`
	SToken          string `json:"sToken"`
}

type getVPPUserSrvResponse struct {
	Status int      `json:"status,omitempty"`
	User   *VPPUser `json:"user,omitempty"`
	*VPPError
}

func (s *usersService) GetUser(user *VPPUser) error {
	sToken, err := s.client.Config.SToken.Base64String()
	if err != nil {
		return err
	}
	var response *getVPPUserSrvResponse
	var request *getVPPUserSrvRequest = &getVPPUserSrvRequest{
		ClientUserIdStr: user.ClientUserIdStr,
		SToken:          sToken,
	}

	req, err := s.client.NewRequest("POST", s.client.Config.serviceConfig.GetUserSrvURL, request)
	if err != nil {
		return err
	}

	err = s.client.Do(req, &response)
	if err != nil {
		return err
	}

	if response.Status == -1 {
		return response.VPPError
	}

	user.UserID = response.User.UserID
	user.Email = response.User.Email
	user.Status = response.User.Status
	user.InviteURL = response.User.InviteURL
	user.InviteCode = response.User.InviteCode

	return nil
}

//
//func (svc *usersService) GetUsers() ([]VPPUser, error) {
//
//}
//
//func (svc *usersService) RetireUser(user *VPPUser) error {
//
//}
//
//func (svc *usersService) EditUser(user *VPPUser) error {
//
//}
