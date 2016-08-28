package vpp

import (
	"github.com/satori/go.uuid"
)

const (
	RegStatusRegistered string = "Registered"
	RegStatusAssociated        = "Associated"
	RegStatusRetired           = "Retired"
)

// VPPUser describes the attributes of a VPP user.
// In most cases, ClientUserIdStr should be used over UserID
type VPPUser struct {
	UserID          int    `json:"userId,omitempty"`
	Email           string `json:"email,omitempty"`
	Status          string `json:"status,omitempty"`
	InviteURL       string `json:"inviteUrl,omitempty"`
	InviteCode      string `json:"inviteCode,omitempty"`
	ClientUserIdStr string `json:"clientUserIdStr,omitempty"`
	ITSIdHash       string `json:"itsIdHash,omitempty"` // empty if no iTunes account has been associated
}

// NewUser creates a new user by generating a UUID.
// Under normal circumstances you should use a UUID that is authoritative for your identity platform
func NewUser(email string, id string) *VPPUser {
	if id == "" {
		id = uuid.NewV4().String()
	}
	return &VPPUser{
		ClientUserIdStr: id,
		Email:           email,
	}
}

// UsersService interface describes the methods available as part of the VPP user management API.
type UsersService interface {
	RegisterUser(user *VPPUser) (*VPPUser, error)
	GetUser(*VPPUser) error
	GetUsers(batch *BatchRequest, opts ...GetUsersOption) ([]VPPUser, error)
	RetireUser(user *VPPUser) error
	EditUser(user *VPPUser) error
}

type usersService struct {
	client *vppClient
	sToken string
}

type registerVPPUserSrvRequest struct {
	*VPPUser
	SToken string `json:"sToken"`
}

type registerVPPUserSrvResponse struct {
	Status Status   `json:"status,omitempty"`
	User   *VPPUser `json:"user,omitempty"`
	*VPPError
}

// RegisterUser registers a new VPP user
func (s *usersService) RegisterUser(user *VPPUser) (*VPPUser, error) {
	var response *registerVPPUserSrvResponse
	var request *registerVPPUserSrvRequest = &registerVPPUserSrvRequest{
		VPPUser: user,
		SToken:  s.sToken,
	}

	req, err := s.client.NewRequest("POST", s.client.Config.serviceConfig.RegisterUserSrvURL, request)
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

	return response.User, nil
}

type getVPPUserSrvRequest struct {
	UserID          int    `json:"userId,omitempty"`
	ClientUserIdStr string `json:"clientUserIdStr,omitempty"`
	ITSIdHash       string `json:"itsIdHash,omitempty"`
	SToken          string `json:"sToken"`
}

type getVPPUserSrvResponse struct {
	Status Status   `json:"status,omitempty"`
	User   *VPPUser `json:"user,omitempty"`
	*VPPError
}

// GetUser gets any number of users associated with the given ClientIdStr
func (s *usersService) GetUser(user *VPPUser) error {
	var response *getVPPUserSrvResponse
	var request *getVPPUserSrvRequest = &getVPPUserSrvRequest{
		ClientUserIdStr: user.ClientUserIdStr,
		SToken:          s.sToken,
	}

	req, err := s.client.NewRequest("POST", s.client.Config.serviceConfig.GetUserSrvURL, request)
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

	user.UserID = response.User.UserID
	user.Email = response.User.Email
	user.Status = response.User.Status
	user.InviteURL = response.User.InviteURL
	user.InviteCode = response.User.InviteCode

	return nil
}

type getUsersRequestOpts struct {
	IncludeRetired int `json:"includeRetired"`
}

// GetUsersOption describes the signature of the closure returned by a function adding an argument to GetUsers
type GetUsersOption func(*getUsersRequestOpts) error

// IncludeRetired is an argument given to GetUsers to include users that have been retired.
// Retiring a user disassociates a VPP user ID from its iTunes account and releases all revocable licenses.
func IncludeRetired(include bool) GetUsersOption {
	return func(opts *getUsersRequestOpts) error {
		if include == true {
			opts.IncludeRetired = 1
		} else {
			opts.IncludeRetired = 0
		}
		return nil
	}
}

type getVPPUsersSrvRequest struct {
	*getUsersRequestOpts
	*BatchRequest
	SToken string `json:"sToken"`
}

type getVPPUsersSrvResponse struct {
	Status     Status    `json:"status,omitempty"`
	Users      []VPPUser `json:"users,omitempty"`
	TotalCount int       `json:"totalCount,omitempty"` // An estimate of the records returned. Will not appear if batchToken exists
	*VPPError
	BatchToken         string `json:"batchToken,omitempty"`
	SinceModifiedToken string `json:"sinceModifiedToken,omitempty"`
}

// GetUsers obtains a list of all known users from the VPP server
func (s *usersService) GetUsers(batch *BatchRequest, opts ...GetUsersOption) ([]VPPUser, error) {
	requestOpts := &getUsersRequestOpts{}
	for _, option := range opts {
		if err := option(requestOpts); err != nil {
			return nil, err
		}
	}
	var request *getVPPUsersSrvRequest = &getVPPUsersSrvRequest{
		getUsersRequestOpts: requestOpts,
		BatchRequest:        batch,
		SToken:              s.sToken,
	}
	var response getVPPUsersSrvResponse
	req, err := s.client.NewRequest("POST", s.client.Config.serviceConfig.GetUsersSrvURL, request)
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

	batch.BatchToken = response.BatchToken
	batch.SinceModifiedToken = response.SinceModifiedToken

	return response.Users, nil
}

type retireVPPUserSrvRequest struct {
	UserId          int    `json:"userId,omitempty"`
	ClientUserIDStr string `json:"clientUserIdStr,omitempty"`
	SToken          string `json:"sToken"`
}

type retireVPPUserSrvResponse struct {
	Status Status   `json:"status,omitempty"`
	User   *VPPUser `json:"user,omitempty"`
	*VPPError
}

// RetireUser disassociates our user id with an itunes user id. All revocable licenses are then freed.
func (s *usersService) RetireUser(user *VPPUser) error {
	var response *retireVPPUserSrvResponse
	var request *retireVPPUserSrvRequest = &retireVPPUserSrvRequest{
		UserId:          user.UserID,
		ClientUserIDStr: user.ClientUserIdStr,
		SToken:          s.sToken,
	}
	req, err := s.client.NewRequest("POST", s.client.Config.serviceConfig.RetireUserSrvURL, request)
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

type editVPPUserSrvRequest struct {
	UserId          int    `json:"userId,omitempty"`
	ClientUserIDStr string `json:"clientUserIdStr,omitempty"`
	Email           string `json:"email"`
	SToken          string `json:"sToken"`
}

type editVPPUserSrvResponse struct {
	Status Status   `json:"status,omitempty"`
	User   *VPPUser `json:"user,omitempty"`
	*VPPError
}

// EditUser edits the e-mail address associated with a user
func (s *usersService) EditUser(user *VPPUser) error {
	var response *editVPPUserSrvResponse
	var request *editVPPUserSrvRequest = &editVPPUserSrvRequest{
		UserId:          user.UserID,
		ClientUserIDStr: user.ClientUserIdStr,
		Email:           user.Email,
		SToken:          s.sToken,
	}
	req, err := s.client.NewRequest("POST", s.client.Config.serviceConfig.EditUserSrvURL, request)
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

	user = response.User
	return nil
}
