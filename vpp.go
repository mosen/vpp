package vpp

import (
	"net/url"
	"net/http"
	"fmt"
)

const (
	defaultBaseURL = "https://vpp.itunes.apple.com/WebObjects/MZFinance.woa/wa"
)

type Config struct {
	url *url.URL
	ServerToken string
}

type VPPError struct {
	Number int `json:"errorNumber,omitempty"`
	Message string `json:"errorMessage,omitempty"`
}

func (v *VPPError) Err() string {
	return fmt.Sprintf("%v: %v", v.Number, v.Message)
}

type VPPClient interface {
	UsersService
}

type vppClient struct {
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	UserAgent string

	Config *Config

	usersService
}

func NewVPPClient(config *Config) (VPPClient, error) {
	if config.url == nil {
		config.url, _ = url.Parse(defaultBaseURL)
	}

	c := &vppClient{client: http.DefaultClient, BaseURL: config.url, Config: config}
	c.usersService = usersService{client: c}

	return c, nil
}

