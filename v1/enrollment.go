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
	Purchased   EnrollmentStatus = "PURCHASED"
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
	ID       string           `json:"id"`
	Status   EnrollmentStatus `json:"status"`
	DeviceID string           `json:"deviceID,omitempty"`
}

type enrollmentsWrap struct {
	Enrollments []*Enrollment `json:"enrollments"`
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
	value := new(Enrollment)
	if err := json.Unmarshal(slurp, value); err != nil {
		return nil, err
	}
	return value, nil
}

type EnrollmentUpdate struct {
	Status   string `json:"status" binding:"required"`
	DeviceID string `json:"device_id,omitempty"`
	ClientID string `json:"client_id,omitempty"`
}

func (c *Client) UpdateEnrollmentByID(id string, update *EnrollmentUpdate) (*Enrollment, error) {
	if id == "" {
		return nil, errNilEnrollmentID
	}
	if update == nil {
		return nil, errNilEnrollmentUpdate
	}
	path := fmt.Sprintf("/safety/media/enrollments/%s", id)
	return c.updateEnrollmentByID(path, update, enrollmentV1API)
}

func (c *Client) updateEnrollmentByID(path string, update *EnrollmentUpdate, versions ...string) (*Enrollment, error) {
	blob, err := json.Marshal(update)
	if err != nil {
		return nil, err
	}
	fullURL := fmt.Sprintf("%s%s", c.baseURL(versions...), path)
	req, err := http.NewRequest("PATCH", fullURL, bytes.NewReader(blob))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	slurp, _, err := c.doReq(req)
	if err != nil {
		return nil, err
	}

	value := new(Enrollment)
	if err := json.Unmarshal(slurp, value); err != nil {
		return nil, err
	}

	return value, nil
}
