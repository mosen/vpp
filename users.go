package vpp

import (
	"github.com/satori/go.uuid"
)

type VPPUser struct {
	UserID int `json:"userId,omitempty"`
	Email string `json:"email,omitempty"`
	Status string `json:"status,omitempty"`
	InviteURL string `json:"inviteUrl,omitempty"`
	InviteCode string `json:"inviteCode,omitempty"`
	ClientUserIdStr string `json:"clientUserIdStr,omitempty"`
}

func NewUser(email string) *VPPUser {
	return &VPPUser{
		ClientUserIdStr: uuid.NewV4().String(),
		Email: email,
	}
}

type UsersService interface {
	RegisterUser(user *VPPUser) error
}

type registerVPPUserSrvResponse struct {
	Status int `json:"status,omitempty"`
	User *VPPUser `json:"user,omitempty"`
}

type usersService struct {
	client *vppClient
}

func (svc *usersService) RegisterUser(user *VPPUser) error {

}

func (svc *usersService) GetUser() (*VPPUser, error) {

}

func (svc *usersService) GetUsers() ([]VPPUser, error) {

}

func (svc *usersService) RetireUser(user *VPPUser) error {

}

func (svc *usersService) EditUser(user *VPPUser) error {

}