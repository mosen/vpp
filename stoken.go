package vpp

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
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

func DecodeSToken(base64str string) (*SToken, error) {
	var jsonValue []byte
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(jsonValue))
	decoder.Read([]byte(base64str))

	token := &SToken{}
	json.Unmarshal(jsonValue, token)

	return token, nil
}
