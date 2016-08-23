package vpp

import "fmt"

type VPPError struct {
	ErrorMessage string `json:"errorMessage"`
	ErrorNumber int `json:"errorNumber"`
}

func (e VPPError) Error() string {
	return fmt.Sprintf("(%s) %s", e.ErrorNumber, e.ErrorMessage)
}