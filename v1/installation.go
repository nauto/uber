// Copyright 2017 orijtech. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package uber

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	path = "/v1/safety/media/installations"
)

var (
	ErrInvalidInstallationID = errors.New("expecting non empty installation id")
	ErrInvalidDeviceID       = errors.New("expecting non empty device id")
)

type CreateInstallation struct {
	DeviceID     string `json:"device_id,binding:required"`
	LicensePlate string `json:"license_plate,binding:required"`
	VIN          string `json:"vin,binding:required"`
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

func (c *Client) CreateInstallation(installation *CreateInstallation) (int, *NewInstallation, error) {
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
		return http.StatusBadRequest, ErrInvalidInstallationID
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

func (c *Client) GetInstallations(query url.Values) (int, *Installations, error) {
	fullURL := fmt.Sprintf("%s?%s", path, query.Encode())
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

func (c *Client) UpdateInstallationByID(installationID string, installation CreateInstallation) (int, error) {
	if installationID == "" {
		return http.StatusBadRequest, ErrInvalidInstallationID
	}

	blob, err := json.Marshal(installation)
	if err != nil {
		return http.StatusBadRequest, err
	}

	fullURL := fmt.Sprintf("%s/%s", path, installationID)
	req, err := http.NewRequest("PATCH", fullURL, bytes.NewReader(blob))
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
