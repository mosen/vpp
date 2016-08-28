package vpp

import (
	"bytes"
	"fmt"
	"net/url"
	"os/exec"
	"testing"
)

var (
	config             *Config
	expDateStr         string
	vppSimURL          *url.URL
	clientIdStrFixture string
)

// Get the current VPPsim token (assuming vppsim is running from this directory)
func VPPSimToken() string {
	tokenCmd := exec.Command("./vppsim", "stoken", "-port", "9001")
	var tokenOut bytes.Buffer
	tokenCmd.Stdout = &tokenOut
	err := tokenCmd.Run()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Got token from VPPsim: %s", tokenOut.String())
	return tokenOut.String()
}

func setup() {
	// get SToken from running VPPsim if possible
	tokenStr := VPPSimToken()

	vppSimURL, _ = url.Parse("http://localhost:9001")
	clientIdStrFixture = "eefe2e28-2f64-46ce-9bf5-726f5eda50a0"

	config = &Config{
		URL:    vppSimURL,
		SToken: tokenStr,
		debug:  true,
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

// Make sure that the client adheres to Retry-After header in seconds
func TestVppClient_Do_RetryAfterSecs(t *testing.T) {
	setup()
	defer teardown()

	vppClient, err := NewVPPClient(config)
	if err != nil {
		t.Fatal(err)
	}

	requestBody := getVPPUserSrvRequest{SToken: config.SToken}
	req, err := vppClient.NewRequest("POST", fmt.Sprintf("%s/getVPPUsersSrv", vppSimURL), &requestBody)
	if err != nil {
		t.Error(err)
	}

	var response getVPPUsersSrvResponse
	err = vppClient.Do(req, &response)
	if err != nil {
		t.Error(err)
	}

	if response.VPPError != nil {
		t.Errorf("VPP error: %s", response.VPPError.Error())
	}
}

// Make sure that the client adheres to Retry-After header in HTTP Date format
func TestVppClient_Do_RetryAfterHTTPDate(t *testing.T) {
	setup()
	defer teardown()

	vppClient, err := NewVPPClient(config)
	if err != nil {
		t.Fatal(err)
	}

	requestBody := getVPPUserSrvRequest{SToken: config.SToken}
	req, err := vppClient.NewRequest("POST", fmt.Sprintf("%s/getVPPUsersSrv", vppSimURL), &requestBody)
	if err != nil {
		t.Error(err)
	}

	var response getVPPUsersSrvResponse
	err = vppClient.Do(req, &response)
	if err != nil {
		t.Error(err)
	}

	if response.VPPError.ErrorMessage != "" {
		t.Errorf("VPP error: %s", response.VPPError.Error())
	}
}
