package vpp

import (
	"fmt"
	"testing"
)

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
