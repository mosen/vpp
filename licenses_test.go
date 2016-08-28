package vpp

import (
	"fmt"
	"testing"
)

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
