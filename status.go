package wdc

import (
	"context"
	"net/http"
	"time"
)

//
// TYPES
//

type StatusValue struct {
	Ready   bool        `json:"ready"`
	Message string      `json:"message"`
	Build   StatusBuild `json:"build"`
	OS      StatusBuild `json:"os"`
	Java    StatusJava  `json:"java"`
}

type StatusBuild struct {
	Revision string    `json:"revision"`
	Time     time.Time `json:"time"`
	Version  string    `json:"version"`
}

type StatusOS struct {
	Arch    string `json:"arch"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type StatusJava struct {
	Java string `json:"java"`
}

//
// RESPONSES
//

type StatusResponse struct {
	Status int         `json:"status"`
	Value  StatusValue `json:"value"`
}

//
// METHODS
//

// Status returns information about whether a remote end is in a state in which it can create new sessions,
// but may additionally include arbitrary meta information that is specific to the implementation.
//
// https://www.w3.org/TR/webdriver/#status
func (c *client) Status(ctx context.Context) (*StatusResponse, error) {
	req, err := c.prepare(http.MethodGet, "status", nil)
	if err != nil {
		return nil, err
	}

	res := new(StatusResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
