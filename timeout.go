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

type Timeout struct {
	Implicit time.Duration
	PageLoad time.Duration
	Script   time.Duration
}

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

type timeoutResponse struct {
	Value struct {
		Implicit uint `json:"implicit"`
		PageLoad uint `json:"pageLoad"`
		Script   uint `json:"script"`
	} `json:"value"`
}

//
// METHODS
//

// Timeouts command is used to return the timeouts associated with the current session.
//
// https://www.w3.org/TR/webdriver/#get-timeouts
func (c *Client) Timeouts(ctx context.Context) (Timeout, error) {
	route := fmt.Sprintf("session/%s/timeouts", c.session.ID)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return Timeout{}, err
	}

	res := new(timeoutResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return Timeout{}, err
	}

	return Timeout{
		Implicit: time.Duration(res.Value.Implicit) * time.Millisecond,
		PageLoad: time.Duration(res.Value.PageLoad) * time.Millisecond,
		Script:   time.Duration(res.Value.Script) * time.Millisecond,
	}, nil
}

// TimeoutElementFind command is used to set the amount of time d the driver should wait when searching for elements.
//
// https://www.w3.org/TR/webdriver/#set-timeouts
func (c *Client) TimeoutElementFind(ctx context.Context, d time.Duration) error {
	r := &timeoutRequest{
		Implicit: uint(d / time.Millisecond),
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

// TimeoutPageLoad command is used to set the amount of time d to interrupt a navigation attempt.
//
// https://www.w3.org/TR/webdriver/#set-timeouts
func (c *Client) TimeoutPageLoad(ctx context.Context, d time.Duration) error {
	r := &timeoutPageLoadRequest{
		PageLoad: uint(d / time.Millisecond),
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

// TimeoutScript command is used to set the amount of time d to interrupt a script that is being evaluated.
//
// https://www.w3.org/TR/webdriver/#set-timeouts
func (c *Client) TimeoutScript(ctx context.Context, d time.Duration) error {
	r := &timeoutScriptRequest{
		Script: uint(d / time.Millisecond),
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
