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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const activityV1API = "v1"

var errInvalidMinMaxTS = errors.New("invalid min/max timestamp")

type Activity struct {
	StartTS time.Time `json:"start_time"`
	EndTS   time.Time `json:"end_time"`
}

type activitiesWrap struct {
	Activities []*Activity `json:"activity"`
}

func (c *Client) ActivitiesByID(id string, query url.Values) ([]*Activity, error) {
	if id == "" {
		return nil, errNilEnrollmentID
	}
	if query.Get("min") == "" || query.Get("max") == "" {
		return nil, errInvalidMinMaxTS
	}

	path := fmt.Sprintf("/safety/media/enrollments/%s/activity?", id) + query.Encode()
	return c.activitiesByID(path, activityV1API)
}

func (c *Client) activitiesByID(path string, versions ...string) ([]*Activity, error) {
	fullURL := fmt.Sprintf("%s%s", c.baseURL(versions...), path)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	slurp, _, err := c.doReq(req)
	if err != nil {
		return nil, err
	}
	values := new(activitiesWrap)
	if err := json.Unmarshal(slurp, values); err != nil {
		return nil, err
	}
	return values.Activities, nil
}
