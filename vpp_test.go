package vpp

import (
	"fmt"
	"net/url"
	"testing"
)

var (
	config             *Config
	expDateStr         string
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

func TestSToken_Base64String(t *testing.T) {
	setup()
	defer teardown()

	encoded, err := config.SToken.Base64String()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("sToken encoded: %s\n", encoded)
}
