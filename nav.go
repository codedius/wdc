package wdc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

// NavigateTo command is used to navigate to a new URL url.
//
// https://www.w3.org/TR/webdriver/#navigate-to
func (c *Client) NavigateTo(ctx context.Context, url string) error {
	if url == "" {
		return errors.New("URL is empty")
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

	return c.do(ctx, req, nil)
}

// NavigateBack command is used to navigate backwards in the browser history, if possible.
//
// This is equivalent to pressing the back button in the browser or invoking window.history.back.
// https://www.w3.org/TR/webdriver/#back
func (c *Client) NavigateBack(ctx context.Context) error {
	route := fmt.Sprintf("session/%s/back", c.session.ID)

	req, err := c.prepare(http.MethodPost, route, nil)
	if err != nil {
		return err
	}

	return c.do(ctx, req, nil)
}

// NavigateForward command is used to navigate forwards in the browser history, if possible.
//
// This is equivalent to pressing the forward button in the browser or invoking window.history.forward.
// https://www.w3.org/TR/webdriver/#forward
func (c *Client) NavigateForward(ctx context.Context) error {
	route := fmt.Sprintf("session/%s/forward", c.session.ID)

	req, err := c.prepare(http.MethodPost, route, nil)
	if err != nil {
		return err
	}

	return c.do(ctx, req, nil)
}
