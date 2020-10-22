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

// Timeouts returns the timeouts associated with the current session.
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

// TimeoutElementFind specifies a time to wait for the element location strategy to complete when locating an element.
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

	err = c.do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}

// TimeoutPageLoad provides the timeout limit used to interrupt an explicit navigation attempt.
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

	err = c.do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}

// TimeoutScript specifies when to interrupt a script that is being evaluated.
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

	err = c.do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}
