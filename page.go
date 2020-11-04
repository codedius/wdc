package wdc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

//
// REQUESTS
//

type scriptRequest struct {
	Script string        `json:"script"`
	Args   []interface{} `json:"args"`
}

//
// METHODS
//

// PageRefresh command is used to refresh the current page.
//
// https://www.w3.org/TR/webdriver/#refresh
func (c *Client) PageRefresh(ctx context.Context) error {
	route := fmt.Sprintf("session/%s/refresh", c.session.ID)

	req, err := c.prepare(http.MethodPost, route, nil)
	if err != nil {
		return err
	}

	return c.do(ctx, req, nil)
}

// PageURL command is used to retrieve the URL of the current page.
//
// https://www.w3.org/TR/webdriver/#get-current-url
func (c *Client) PageURL(ctx context.Context) (string, error) {
	route := fmt.Sprintf("session/%s/url", c.session.ID)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return "", err
	}

	res := new(stringValue)

	err = c.do(ctx, req, res)
	if err != nil {
		return "", err
	}

	return res.Value, nil
}

// PageTitle command is used to get the current page title.
//
// This is equivalent to calling document.title.
// https://www.w3.org/TR/webdriver/#get-title
func (c *Client) PageTitle(ctx context.Context) (string, error) {
	route := fmt.Sprintf("session/%s/title", c.session.ID)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return "", err
	}

	res := new(stringValue)

	err = c.do(ctx, req, res)
	if err != nil {
		return "", err
	}

	return res.Value, nil
}

// PageSource command is used to get the current page source.
//
// https://www.w3.org/TR/webdriver/#get-page-source
func (c *Client) PageSource(ctx context.Context) (string, error) {
	route := fmt.Sprintf("session/%s/title", c.session.ID)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return "", err
	}

	res := new(stringValue)

	err = c.do(ctx, req, res)
	if err != nil {
		return "", err
	}

	return res.Value, nil
}

// PageScreenshot command is used to take a screenshot of the current page and return base64 string.
//
// https://www.w3.org/TR/webdriver/#take-screenshot
func (c *Client) PageScreenshot(ctx context.Context) (string, error) {
	route := fmt.Sprintf("session/%s/screenshot", c.session.ID)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return "", err
	}

	res := new(stringValue)

	err = c.do(ctx, req, res)
	if err != nil {
		return "", err
	}

	return res.Value, nil
}

// PageScript command is used to inject a snippet of JavaScript into the page for execution in the context of the currently selected frame.
//
// The executed script is assumed to be synchronous and the result of evaluating the script is returned to the client.
// https://www.w3.org/TR/webdriver/#execute-script
func (c *Client) PageScript(ctx context.Context, s string, args []interface{}) (string, error) {
	if s == "" {
		return "", ErrorScriptIsEmpty
	}

	r := &scriptRequest{Script: s, Args: args}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return "", err
	}

	route := fmt.Sprintf("session/%s/execute/sync", c.session.ID)

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return "", err
	}

	res := new(stringValue)

	err = c.do(ctx, req, res)
	if err != nil {
		return "", err
	}

	return res.Value, nil
}

// PageScriptAsync command is used to inject a snippet of JavaScript into the page for execution in the context of the currently selected frame.
//
// The executed script is assumed to be asynchronous and must signal that is done by invoking the provided callback, which is always provided as the final argument to the function. The value to this callback will be returned to the client.
// https://www.w3.org/TR/webdriver/#execute-script
func (c *Client) PageScriptAsync(ctx context.Context, s string, args []interface{}) (string, error) {
	if s == "" {
		return "", ErrorScriptIsEmpty
	}

	r := &scriptRequest{Script: s, Args: args}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return "", err
	}

	route := fmt.Sprintf("session/%s/execute/async", c.session.ID)

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return "", err
	}

	res := new(stringValue)

	err = c.do(ctx, req, res)
	if err != nil {
		return "", err
	}

	return res.Value, nil
}
