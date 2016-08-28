package vpp

import (
	"fmt"
	"testing"
)

func TestUsersService_RegisterUser(t *testing.T) {
	setup()
	defer teardown()

	vppClient, err := NewVPPClient(config)
	if err != nil {
		t.Error(err)
	}

	user := NewUser("test@localhost", clientIdStrFixture)
	vppUser, err := vppClient.RegisterUser(user)
	if err != nil {
		t.Error(err)
	}

	if vppUser.UserID == 0 {
		t.Error("Expected UserId but got empty value")
	}

	if vppUser.Status != RegStatusRegistered {
		t.Errorf("Expected status registered but got %s", vppUser.Status)
	}

	t.Logf("Registered user successfully with ClientIDStr: %s", vppUser.ClientUserIdStr)
}

func TestUsersService_GetUser(t *testing.T) {
	setup()
	defer teardown()

	vppClient, err := NewVPPClient(config)
	if err != nil {
		t.Error(err)
	}

	user := &VPPUser{ClientUserIdStr: clientIdStrFixture}
	err = vppClient.GetUser(user)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%#v\n", user)
}

func TestUsersService_GetUsers(t *testing.T) {
	setup()
	defer teardown()

	vppClient, err := NewVPPClient(config)
	if err != nil {
		t.Error(err)
	}

	batch := &BatchRequest{}
	_, err = vppClient.GetUsers(batch)
	if err != nil {
		t.Error(err)
	}

	if batch.BatchToken == "" {
		t.Error("response did not return a batch token")
	}
}

func TestUsersService_RetireUser(t *testing.T) {
	setup()
	defer teardown()

	vppClient, err := NewVPPClient(config)
	if err != nil {
		t.Error(err)
	}

	user := &VPPUser{ClientUserIdStr: clientIdStrFixture}
	err = vppClient.RetireUser(user)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%#v\n", user)
}

func TestUsersService_EditUser(t *testing.T) {
	setup()
	defer teardown()

	vppClient, err := NewVPPClient(config)
	if err != nil {
		t.Error(err)
	}

	user := &VPPUser{
		ClientUserIdStr: clientIdStrFixture,
		Email:           "updated@email",
	}
	err = vppClient.EditUser(user)
	if err != nil {
		t.Error(err)
	}
}
