package wdc

import (
	"context"
	"fmt"
	"net/http"
)

//
// METHODS
//

// PageRefresh command causes the browser to reload the page in the current top-level browsing context.
//
// https://www.w3.org/TR/webdriver/#refresh
func (c *Client) PageRefresh(ctx context.Context) error {
	route := fmt.Sprintf("session/%s/refresh", c.session.ID)

	req, err := c.prepare(http.MethodPost, route, nil)
	if err != nil {
		return err
	}

	err = c.do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}

// PageURL command returns the current URL.
//
// https://www.w3.org/TR/webdriver/#get-current-url
func (c *Client) PageURL(ctx context.Context) (string, error) {
	route := fmt.Sprintf("session/%s/url", c.session.ID)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return "", err
	}

	res := new(value)

	err = c.do(ctx, req, res)
	if err != nil {
		return "", err
	}

	return res.Value, nil
}

// PageTitle command returns the document title of the current top-level browsing context.
//
// This is equivalent to calling document.title.
// https://www.w3.org/TR/webdriver/#get-title
func (c *Client) PageTitle(ctx context.Context) (string, error) {
	route := fmt.Sprintf("session/%s/title", c.session.ID)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return "", err
	}

	res := new(value)

	err = c.do(ctx, req, res)
	if err != nil {
		return "", err
	}

	return res.Value, nil
}