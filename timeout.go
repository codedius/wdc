package wdc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//
// REQUESTS
//

type timeoutRequest struct {
	Implicit uint `json:"implicit"`
}

type timeoutPageLoadRequest struct {
	PageLoad uint `json:"pageLoad"`
}

type timeoutScriptRequest struct {
	Script uint `json:"script"`
}

//
// RESPONSES
//

type TimeoutsResponse struct {
	Implicit uint `json:"implicit"`
	PageLoad uint `json:"pageLoad"`
	Script   uint `json:"script"`
}

//
// METHODS
//

// Timeouts command is used to return the timeouts associated with the current session.
//
// https://www.w3.org/TR/webdriver/#get-timeouts
func (c *Client) Timeouts(ctx context.Context) (*TimeoutsResponse, error) {
	route := fmt.Sprintf("session/%s/timeouts", c.session.ID)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return nil, err
	}

	res := new(TimeoutsResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// TimeoutElementFind command is used to set the amount of time the driver should wait when searching for elements.
//
// https://www.w3.org/TR/webdriver/#set-timeouts
func (c *Client) TimeoutElementFind(ctx context.Context, timeout time.Duration) error {
	r := &timeoutRequest{
		Implicit: uint(timeout / time.Millisecond),
	}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return err
	}

	route := fmt.Sprintf("session/%s/timeouts", c.session.ID)

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return err
	}

	return c.do(ctx, req, nil)
}

// TimeoutPageLoad command is used to set the amount of time to interrupt a navigation attempt.
//
// https://www.w3.org/TR/webdriver/#set-timeouts
func (c *Client) TimeoutPageLoad(ctx context.Context, timeout time.Duration) error {
	r := &timeoutPageLoadRequest{
		PageLoad: uint(timeout / time.Millisecond),
	}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return err
	}

	route := fmt.Sprintf("session/%s/timeouts", c.session.ID)

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return err
	}

	return c.do(ctx, req, nil)
}

// TimeoutScript command is used to set the amount of time to interrupt a script that is being evaluated.
//
// https://www.w3.org/TR/webdriver/#set-timeouts
func (c *Client) TimeoutScript(ctx context.Context, timeout time.Duration) error {
	r := &timeoutScriptRequest{
		Script: uint(timeout / time.Millisecond),
	}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return err
	}

	route := fmt.Sprintf("session/%s/timeouts", c.session.ID)

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return err
	}

	return c.do(ctx, req, nil)
}
