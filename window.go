package wdc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

//
// TYPES
//

// WindowSize represents size and position of a window.
type WindowSize struct {
	Height int `json:"height"`
	Width  int `json:"width"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

// WindowID represents an id of a window.
type WindowID string

//
// REQUESTS
//

type windowRequest struct {
	Name string `json:"name"`
}

//
// RESPONSES
//

type windowsResponse struct {
	Value []WindowID `json:"value"`
}

type windowSizeResponse struct {
	Value *WindowSize `json:"value"`
}

//
// METHODS
//

// WindowID command is used to retrieve the current window handle.
//
// https://www.w3.org/TR/webdriver/#get-window-handle
func (c *Client) WindowID(ctx context.Context) (WindowID, error) {
	route := fmt.Sprintf("session/%s/window", c.session.ID)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return "", err
	}

	res := new(stringValue)

	err = c.do(ctx, req, res)
	if err != nil {
		return "", err
	}

	return WindowID(res.Value), nil
}

// WindowIDs command is used to retrieve the list of all window handles available to the session.
//
// https://www.w3.org/TR/webdriver/#get-window-handles
func (c *Client) WindowIDs(ctx context.Context) ([]WindowID, error) {
	route := fmt.Sprintf("session/%s/window/handles", c.session.ID)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return nil, err
	}

	res := new(windowsResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	return res.Value, nil
}

// WindowNew command is used to create new window.
//
// https://www.w3.org/TR/webdriver/#new-window
func (c *Client) WindowNew(ctx context.Context) (WindowID, error) {
	route := fmt.Sprintf("session/%s/window/new", c.session.ID)

	req, err := c.prepare(http.MethodPost, route, nil)
	if err != nil {
		return "", err
	}

	res := new(stringValue)

	err = c.do(ctx, req, res)
	if err != nil {
		return "", err
	}

	return WindowID(res.Value), nil
}

// WindowClose command is used to close a window.
//
// https://www.w3.org/TR/webdriver/#close-window
func (c *Client) WindowClose(ctx context.Context) error {
	route := fmt.Sprintf("session/%s/window", c.session.ID)

	req, err := c.prepare(http.MethodDelete, route, nil)
	if err != nil {
		return err
	}

	return c.do(ctx, req, nil)
}

// WindowClose command is used to close a window.
//
// https://www.w3.org/TR/webdriver/#switch-to-window
func (c *Client) WindowSwitch(ctx context.Context, wid string) error {
	if wid == "" {
		return ErrorWindowIDIsEmpty
	}

	r := &windowRequest{Name: wid}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return err
	}

	route := fmt.Sprintf("session/%s/window", c.session.ID)

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return err
	}

	return c.do(ctx, req, nil)
}

// WindowSize command is used to get the size of the specified window.
//
// https://www.w3.org/TR/webdriver/#get-window-rect
func (c *Client) WindowSize(ctx context.Context) (*WindowSize, error) {
	route := fmt.Sprintf("session/%s/window/rect", c.session.ID)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return nil, err
	}

	res := new(windowSizeResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	return res.Value, nil
}

// WindowResize command is used to change the size of the specified window.
//
// https://www.w3.org/TR/webdriver/#set-window-rect
func (c *Client) WindowResize(ctx context.Context, w, h, x, y int) error {
	if w == 0 {
		return ErrorWindowWidthIsEmpty
	}
	if h == 0 {
		return ErrorWindowHeightIsEmpty
	}

	r := &WindowSize{Width: w, Height: h, X: x, Y: y}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return err
	}

	route := fmt.Sprintf("session/%s/window/rect", c.session.ID)

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return err
	}

	return c.do(ctx, req, nil)
}

// WindowMaximize command is used to maximize the specified window.
//
// https://www.w3.org/TR/webdriver/#maximize-window
func (c *Client) WindowMaximize(ctx context.Context, wid string) (*WindowSize, error) {
	if wid == "" {
		return nil, ErrorWindowIDIsEmpty
	}

	route := fmt.Sprintf("session/%s/window/maximize", c.session.ID)

	r := &windowRequest{Name: wid}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return nil, err
	}

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return nil, err
	}

	res := new(windowSizeResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	return res.Value, nil
}

// WindowMinimize command is used to minimize the specified window.
//
// https://www.w3.org/TR/webdriver/#minimize-window
func (c *Client) WindowMinimize(ctx context.Context, wid string) (*WindowSize, error) {
	if wid == "" {
		return nil, ErrorWindowIDIsEmpty
	}

	route := fmt.Sprintf("session/%s/window/minimize", c.session.ID)

	r := &windowRequest{Name: wid}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return nil, err
	}

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return nil, err
	}

	res := new(windowSizeResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	return res.Value, nil
}

// WindowFullscreen command is used to fullscreen the specified window.
//
// https://www.w3.org/TR/webdriver/#fullscreen-window
func (c *Client) WindowFullscreen(ctx context.Context, wid string) (*WindowSize, error) {
	if wid == "" {
		return nil, ErrorWindowIDIsEmpty
	}

	route := fmt.Sprintf("session/%s/window/fullscreen", c.session.ID)

	r := &windowRequest{Name: wid}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return nil, err
	}

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return nil, err
	}

	res := new(windowSizeResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	return res.Value, nil
}
