package vpp

import (
	"testing"
	"net/url"
	"time"
)

var (
	config *Config
	expDate *time.Time
)

func setup() {
	expDate = time.Now()
	expDate = expDate.AddDate(1, 0, 0)

	config = &Config{
		URL: url.Parse("http://localhost:9001"),
		SToken: &SToken{
			Token: "VGhpcyBpcyBhIHNhbXBsZSB0ZXh0IHdoaWNoIHdhcyB1c2VkIHRvIGNyZWF0ZSB0aGUgc2ltdWxhdG9yIHRva2VuCg==",
			ExpDate: expDate.String(),
			OrgName: "Example Inc.",
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

	sconfig, err := vppClient.ServiceConfig()

}

