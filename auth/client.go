package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"os"
)

var errLoginGagal = errors.New("login gagal")

type Payload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Mac      string `json:"mac"`
	BotID    int    `json:"bot_id"`
	Version  string `json:"version"`
}

type AuthClient struct {
	Endpoint string
	client   *http.Client
}

func getMacAddress() string {
	ifas, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}

	return as[0]
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return hostname
}

func (c *AuthClient) Login(email string, password string, botID int, version string) error {
	payload := Payload{
		Name:     getHostname(),
		Mac:      getMacAddress(),
		Email:    email,
		Password: password,
		BotID:    botID,
		Version:  version,
	}

	raw, err := json.Marshal(&payload)
	// log.Println(string(raw))
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.Endpoint, bytes.NewBuffer(raw))
	req.Header.Set("Content-Type", "aplication/json")
	req.Header.Set("Accept", "aplication/json")

	if err != nil {
		return err
	}

	res, _ := c.client.Do(req)
	if res.StatusCode != 200 {
		return errLoginGagal
	}

	return nil
}

func NewAuthClient(endpoint string) *AuthClient {
	return &AuthClient{
		Endpoint: endpoint,
		client:   &http.Client{},
	}
}
