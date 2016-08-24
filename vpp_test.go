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
	vppSimURL *url.URL
)

func setup() {
	expDateStr = "2017-08-24T22:11:47+10:00"
	//expDate = time.Parse("", expDateStr)

	//expDate = time.Now()
	//expDate = expDate.AddDate(1, 0, 0)
	vppSimURL, _ = url.Parse("http://localhost:9001")

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

//func TestNewVPPClient(t *testing.T) {
//	setup()
//	defer teardown()
//
//	_, err := NewVPPClient(config)
//	if err != nil {
//		t.Fatal(err)
//	}
//}

//func TestConfigService_ServiceConfig(t *testing.T) {
//	setup()
//	defer teardown()
//
//	vppClient, err := NewVPPClient(config)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	_, err = vppClient.ServiceConfig()
//	//fmt.Printf("%#v\n", sconfig)
//}

//func TestUsersService_RegisterUser(t *testing.T) {
//	setup()
//	defer teardown()
//
//	vppClient, err := NewVPPClient(config)
//	if err != nil {
//		t.Error(err)
//	}
//
//	user := NewUser("test@localhost", "")
//	vppUser, err := vppClient.RegisterUser(user)
//	if err != nil {
//		t.Error(err)
//	}
//	fmt.Printf("%#v\n", vppUser)
//}

func TestUsersService_GetUser(t *testing.T) {
	setup()
	defer teardown()

	vppClient, err := NewVPPClient(config)
	if err != nil {
		t.Error(err)
	}

	user := &VPPUser{ClientUserIdStr: "b5497f2e-6d96-4ac6-8578-c96414433145"}
	err = vppClient.GetUser(user)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%#v\n", user)
}

//func TestSToken_Base64String(t *testing.T) {
//	setup()
//	defer teardown()
//
//	encoded, err := config.SToken.Base64String()
//	if err != nil {
//		t.Error(err)
//	}
//	fmt.Printf("sToken encoded: %s\n", encoded)
//}
