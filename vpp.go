package vpp

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const (
	libraryVersion = "0.0.1"
	defaultBaseURL = "https://vpp.itunes.apple.com/WebObjects/MZFinance.woa/wa"
	userAgent      = "micromdm/" + libraryVersion
	mediaType      = "application/json;charset=UTF8"
)

type Status int

const (
	StatusErr Status = iota - 1
	StatusOk
)

type SToken struct {
	Token      string `json:"token"`
	ExpDateStr string `json:"expDate"`
	OrgName    string `json:"orgName"`
}

func (t *SToken) Base64String() (string, error) {

	jsonValue, err := json.Marshal(t)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, buf)
	encoder.Write(jsonValue)
	encoder.Close()

	return string(buf.Bytes()), nil
}

type Config struct {
	URL           *url.URL
	SToken        *SToken
	debug         bool
	serviceConfig *ServiceConfig
}

type VPPError struct {
	ErrorMessage string `json:"errorMessage"`
	ErrorNumber  int    `json:"errorNumber"`
}

func (e VPPError) Error() string {
	return fmt.Sprintf("(%d) %s", e.ErrorNumber, e.ErrorMessage)
}

type VPPClient interface {
	AssetsService
	ConfigService
	LicensesService
	MetadataService
	UsersService
}

type vppClient struct {
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	UserAgent string

	Config *Config

	assetsService
	configService
	licensesService
	metadataService
	usersService
}

func (c *vppClient) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)

	return req, nil
}

// Do sends an API request and returns the API response.
func (c *vppClient) Do(req *http.Request, into interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	// The VPP service may issue either a 3xx (redirect) or 503 (unavailable) with a Retry-After header
	// if the service is overloaded or this client is causing too much load.
	// The Retry-After header may be in seconds or as a date.

	//if resp.StatusCode == http.StatusPermanentRedirect || resp.StatusCode == http.StatusTemporaryRedirect {
	//	retryAfter := resp.Header.Get("Retry-After")
	//}
	//
	//if resp.StatusCode == http.StatusServiceUnavailable {
	//	retryAfter := resp.Header.Get("Retry-After")
	//}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("VPP API Error: %v", string(body))
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	return decodeJSON(c.Config.debug, resp.Body, into)
}

func decodeJSON(debug bool, body io.Reader, into interface{}) error {
	var dec *json.Decoder
	if debug {
		dec = json.NewDecoder(io.TeeReader(body, os.Stdout))
	} else {
		dec = json.NewDecoder(body)
	}

	return dec.Decode(into)
}

func NewVPPClient(config *Config) (VPPClient, error) {
	if config.URL == nil {
		config.URL, _ = url.Parse(defaultBaseURL)
	}

	c := &vppClient{client: http.DefaultClient, BaseURL: config.URL, Config: config}
	sToken, err := config.SToken.Base64String()
	if err != nil {
		return nil, err
	}

	c.configService = configService{client: c, sToken: sToken}

	serviceConfig, err := c.ServiceConfig()
	if err != nil {
		return nil, err
	}
	c.Config.serviceConfig = serviceConfig

	c.assetsService = assetsService{client: c, sToken: sToken}
	c.licensesService = licensesService{client: c, sToken: sToken}
	c.metadataService = metadataService{client: c, sToken: sToken}
	c.usersService = usersService{client: c, sToken: sToken}

	return c, nil
}
