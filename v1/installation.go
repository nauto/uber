package uber

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	path = "/v1/safety/media/installations"
)

type CreateInstallation struct {
	DeviceID     string `json:"device_id,binding:required"`
	LicensePlate string `json:"license_plate,binding:required"`
	Vin          string `json:"vin,binding:required"`
}

type Installation struct {
	ID       string `json:"id,binding:required"`
	DeviceID string `json:"device_id,binding:required"`
}

type NewInstallation struct {
	Installation *Installation `json:"installation,binding:required"`
}

type Installations struct {
	Installations []*Installation `json:"installations,binding:required"`
	Limit         int             `json:"limit,omitempty"`
	Offset        int             `json:"offset,omitempty"`
}

func (c *Client) CreateInstallation(installation CreateInstallation) (int, *NewInstallation, error) {
	blob, err := json.Marshal(installation)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	req, err := http.NewRequest("POST", path, bytes.NewReader(blob))
	if err != nil {
		return req.Response.StatusCode, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, _, err := c.doReq(req)
	if err != nil {
		return req.Response.StatusCode, nil, err
	}

	value := new(NewInstallation)
	if err := json.Unmarshal(resp, value); err != nil {
		return req.Response.StatusCode, nil, err
	}

	return req.Response.StatusCode, value, nil
}

func (c *Client) DeleteInstallationByID(installationID string) (int, error) {
	if installationID == "" {
		return http.StatusBadRequest, errors.New("expecting non empty installation id")
	}
	fullURL := fmt.Sprintf("%s/%s", path, installationID)
	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		return req.Response.StatusCode, err
	}
	req.Header.Set("Content-Type", "application/json")
	_, _, err = c.doReq(req)
	if err != nil {
		return req.Response.StatusCode, err
	}

	return req.Response.StatusCode, nil
}

func (c *Client) GetInstallationsByDeviceID(deviceID string) (int, *Installations, error) {
	if deviceID == "" {
		return http.StatusBadRequest, nil, errors.New("expecting non empty device id")
	}
	fullURL := fmt.Sprintf("%s?device_id=%s", path, deviceID)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return req.Response.StatusCode, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, _, err := c.doReq(req)
	if err != nil {
		return req.Response.StatusCode, nil, err
	}

	installations := new(Installations)
	if err := json.Unmarshal(resp, installations); err != nil {
		return req.Response.StatusCode, nil, err
	}

	return req.Response.StatusCode, installations, nil
}
