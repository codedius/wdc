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

type frameChangeRequest struct {
	ID int `json:"id"`
}

// SwitchToFrame command switches the active frame to a nested frame by index fi. The active frame receives commands.
//
// https://www.w3.org/TR/webdriver/#switch-to-frame
func (c *Client) SwitchToFrame(ctx context.Context, fi int) error {
	r := &frameChangeRequest{ID: fi}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return err
	}

	route := fmt.Sprintf("session/%s/frame", c.session.ID)

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return err
	}

	return c.do(ctx, req, nil)
}

// SwitchToParentFrame command switches the active frame to the parent frame. The active frame receives commands.
//
// https://www.w3.org/TR/webdriver/#switch-to-parent-frame
func (c *Client) SwitchToParentFrame(ctx context.Context) error {
	route := fmt.Sprintf("/session/%s/frame/parent", c.session.ID)

	req, err := c.prepare(http.MethodPost, route, nil)
	if err != nil {
		return err
	}

	return c.do(ctx, req, nil)
}
