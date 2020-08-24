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
	"strings"
)

type EnrollmentStatus string

const (
	Offered     EnrollmentStatus = "OFFERED"
	Consented   EnrollmentStatus = "CONSENTED"
	Purchased   EnrollmentStatus = "PURCHASED"
	Scheduled   EnrollmentStatus = "SCHEDULED"
	Enabled     EnrollmentStatus = "ENABLED"
	Deactivated EnrollmentStatus = "DEACTIVATED"
	Blocked     EnrollmentStatus = "BLOCKED"
)

const enrollmentV1API = "v1"

var (
	errNilEnrollmentID     = errors.New("expecting a non-empty enrollment id")
	errNilEnrollmentUpdate = errors.New("expecting a non-empty update map")
)

type Enrollment struct {
	ID           string           `json:"id"`
	Status       EnrollmentStatus `json:"status"`
	DeviceID     string           `json:"device_id,omitempty"`
	LicensePlate string           `json:"license_plate,omitempty"`
	Vin          string           `json:"vin,omitempty"`
}

type enrollmentsWrap struct {
	Enrollments []*Enrollment `json:"enrollments"`
}

type enrollmentWrap struct {
	Enrollment *Enrollment `json:"enrollment"`
}

func (c *Client) Enrollments(query url.Values) ([]*Enrollment, error) {
	// support email for now
	path := "/safety/media/enrollments"
	if email := query.Get("email"); email != "" {
		path = strings.Join([]string{path, "?email=", strings.ToLower(email)}, "")
	}
	return c.enrollments(path, enrollmentV1API)
}

func (c *Client) enrollments(path string, versions ...string) ([]*Enrollment, error) {
	fullURL := fmt.Sprintf("%s%s", c.baseURL(versions...), path)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	slurp, _, err := c.doReq(req)
	if err != nil {
		return nil, err
	}
	values := new(enrollmentsWrap)
	if err := json.Unmarshal(slurp, values); err != nil {
		return nil, err
	}
	return values.Enrollments, nil
}

func (c *Client) EnrollmentByID(id string) (*Enrollment, error) {
	if id == "" {
		return nil, errNilEnrollmentID
	}
	path := fmt.Sprintf("/safety/media/enrollments/%s", id)
	return c.enrollmentByID(path, enrollmentV1API)
}

func (c *Client) enrollmentByID(path string, versions ...string) (*Enrollment, error) {
	fullURL := fmt.Sprintf("%s%s", c.baseURL(versions...), path)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	slurp, _, err := c.doReq(req)
	if err != nil {
		return nil, err
	}
	value := new(enrollmentWrap)
	if err := json.Unmarshal(slurp, value); err != nil {
		return nil, err
	}
	return value.Enrollment, nil
}

type EnrollmentUpdate struct {
	Status       EnrollmentStatus `json:"status" binding:"required"`
	DeviceID     string           `json:"device_id,omitempty"`
	LicensePlate string           `json:"license_plate,omitempty"`
	Vin          string           `json:"vin,omitempty"`
	ClientID     string           `json:"client_id" binding:"required"`
}

func (c *Client) UpdateEnrollmentByID(enrollmentID string, update *EnrollmentUpdate) (int, *Enrollment, error) {
	if enrollmentID == "" {
		return http.StatusBadRequest, nil, errNilEnrollmentID
	}
	if update == nil {
		return http.StatusBadRequest, nil, errNilEnrollmentUpdate
	}
	path := fmt.Sprintf("/safety/media/enrollments/%s", enrollmentID)
	return c.updateEnrollmentByID(path, update, enrollmentV1API)
}

func (c *Client) updateEnrollmentByID(path string, update *EnrollmentUpdate, versions ...string) (int, *Enrollment, error) {
	blob, err := json.Marshal(update)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	fullURL := fmt.Sprintf("%s%s", c.baseURL(versions...), path)
	req, err := http.NewRequest("PATCH", fullURL, bytes.NewReader(blob))
	if err != nil {
		return req.Response.StatusCode, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	slurp, _, err := c.doReq(req)
	if err != nil {
		return req.Response.StatusCode, nil, err
	}

	value := new(Enrollment)
	if err := json.Unmarshal(slurp, value); err != nil {
		return req.Response.StatusCode, nil, err
	}

	return req.Response.StatusCode, value, nil
}
