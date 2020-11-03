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

type navigateRequest struct {
	URL string `json:"url"`
}

//
// METHODS
//

// NavigateTo command causes the user agent to navigate the current top-level browsing context to a new location.
//
// https://www.w3.org/TR/webdriver/#navigate-to
func (c *Client) NavigateTo(ctx context.Context, url string) error {
	if url == "" {
		return ErrorURLIsRequired
	}

	r := &navigateRequest{URL: url}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return err
	}

	route := fmt.Sprintf("session/%s/url", c.session.ID)

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

// NavigateBack command causes the browser to traverse one step backward in the joint session history of the current top-level browsing context.
//
// This is equivalent to pressing the back button in the browser chrome or invoking window.history.back.
// https://www.w3.org/TR/webdriver/#back
func (c *Client) NavigateBack(ctx context.Context) error {
	route := fmt.Sprintf("session/%s/back", c.session.ID)

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

// NavigateForward command causes the browser to traverse one step forwards in the joint session history of the current top-level browsing context.
//
// This is equivalent to pressing the forward button in the browser chrome or invoking window.history.forward.
// https://www.w3.org/TR/webdriver/#forward
func (c *Client) NavigateForward(ctx context.Context) error {
	route := fmt.Sprintf("session/%s/forward", c.session.ID)

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
