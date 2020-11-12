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

type Status struct {
	Status int         `json:"status"`
	Value  StatusValue `json:"value"`
}

//
// METHODS
//

// Status command is used to return the server's current status.
//
// https://www.w3.org/TR/webdriver/#status
func (c *Client) Status(ctx context.Context) (Status, error) {
	req, err := c.prepare(http.MethodGet, "status", nil)
	if err != nil {
		return Status{}, err
	}

	res := new(Status)

	err = c.do(ctx, req, res)
	if err != nil {
		return Status{}, err
	}

	return *res, nil
}
