package vpp

import (
	"fmt"
	"net/url"
	"testing"
)

var (
	config     *Config
	expDateStr string
	//expDate time.Time
	vppSimURL          *url.URL
	clientIdStrFixture string
)

func setup() {
	// If you want to use VPPsim, you have to paste the expDate time here
	expDateStr = "2017-08-26T10:24:11+10:00"
	vppSimURL, _ = url.Parse("http://localhost:9001")
	clientIdStrFixture = "eefe2e28-2f64-46ce-9bf5-726f5eda50a0"

	config = &Config{
		URL: vppSimURL,
		SToken: &SToken{
			Token:      "VGhpcyBpcyBhIHNhbXBsZSB0ZXh0IHdoaWNoIHdhcyB1c2VkIHRvIGNyZWF0ZSB0aGUgc2ltdWxhdG9yIHRva2VuCg==",
			ExpDateStr: expDateStr,
			OrgName:    "Example Inc.",
		},
		debug: true,
	}
}

func teardown() {

}

func TestNewVPPClient(t *testing.T) {
	setup()
	defer teardown()

	_, err := NewVPPClient(config)
	if err != nil {
		t.Fatal(err)
	}
}

func TestConfigService_ServiceConfig(t *testing.T) {
	setup()
	defer teardown()

	vppClient, err := NewVPPClient(config)
	if err != nil {
		t.Fatal(err)
	}

	_, err = vppClient.ServiceConfig()
	//fmt.Printf("%#v\n", sconfig)
}

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

	_, err = vppClient.GetUsers()
	if err != nil {
		t.Error(err)
	}

	//fmt.Printf("%#v\n", users)
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

func TestSToken_Base64String(t *testing.T) {
	setup()
	defer teardown()

	encoded, err := config.SToken.Base64String()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("sToken encoded: %s\n", encoded)
}

func TestLicensesService_GetLicenses(t *testing.T) {
	setup()
	defer teardown()

	vppClient, err := NewVPPClient(config)
	if err != nil {
		t.Error(err)
	}

	licenses, err := vppClient.GetLicenses()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%#v\n", licenses)
}

func TestConfigService_ClientContext(t *testing.T) {
	setup()
	defer teardown()

	vppClient, err := NewVPPClient(config)
	if err != nil {
		t.Error(err)
	}

	clientContext, err := vppClient.ClientContext()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%#v\n", clientContext)
}

func TestConfigService_UpdateClientContext(t *testing.T) {
	setup()
	defer teardown()

	vppClient, err := NewVPPClient(config)
	if err != nil {
		t.Error(err)
	}

	cc := NewClientContext("localhost")
	countryCode, err := vppClient.UpdateClientContext(cc)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%#v\n", countryCode)
}

func TestLicensesService_AssociateLicense(t *testing.T) {
	setup()
	defer teardown()

	vppClient, err := NewVPPClient(config)
	if err != nil {
		t.Error(err)
	}

	if err := vppClient.AssociateLicense(&VPPUser{ClientUserIdStr: clientIdStrFixture}, &VPPLicense{LicenseID: "1"}); err != nil {
		t.Error(err)
	}
}
