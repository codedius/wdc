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
// TYPES
//

type Cookie struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Path     string `json:"path"`
	Domain   string `json:"domain"`
	Secure   bool   `json:"secure"`
	HTTPOnly bool   `json:"httpOnly"`
	Expiry   int    `json:"expiry"`
}

//
// REQUESTS
//

type cookieRequest struct {
	Cookie Cookie `json:"cookie"`
}

//
// RESPONSES
//

type cookieResponse struct {
	Value Cookie `json:"value"`
}

type cookiesResponse struct {
	Value []Cookie `json:"value"`
}

//
// METHODS
//

// Cookie command is used to retrieve the cookie with the given name n.
//
// https://www.w3.org/TR/webdriver/#get-named-cookie
func (c *Client) Cookie(ctx context.Context, n string) (Cookie, error) {
	if n == "" {
		return Cookie{}, errors.New("cookie name is empty")
	}

	route := fmt.Sprintf("session/%s/cookie/%s", c.session.ID, n)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return Cookie{}, err
	}

	res := new(cookieResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return Cookie{}, err
	}

	return res.Value, nil
}

// CookieSet command is used to set a cookie v.
//
// https://www.w3.org/TR/webdriver/#add-cookie
func (c *Client) CookieSet(ctx context.Context, v Cookie) error {
	if v.Name == "" {
		return errors.New("cookie name field is empty")
	}
	if v.Value == "" {
		return errors.New("cookie value field is empty")
	}
	if v.Path == "" {
		return errors.New("cookie path field is empty")
	}
	if v.Domain == "" {
		return errors.New("cookie domain field is empty")
	}
	if v.Expiry == 0 {
		return errors.New("cookie expiry field is empty")
	}

	r := &cookieRequest{Cookie: v}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return err
	}

	route := fmt.Sprintf("session/%s/cookie", c.session.ID)

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return err
	}

	res := new(cookieResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return err
	}

	return nil
}

// CookieDelete command is used to delete the cookie with the given name n.
//
// https://www.w3.org/TR/webdriver/#delete-cookie
func (c *Client) CookieDelete(ctx context.Context, n string) error {
	if n == "" {
		return errors.New("cookie name is empty")
	}

	route := fmt.Sprintf("session/%s/cookie/%s", c.session.ID, n)

	req, err := c.prepare(http.MethodDelete, route, nil)
	if err != nil {
		return err
	}

	return c.do(ctx, req, nil)
}

// Cookies command is used to retrieve all cookies visible to the current page.
//
// https://www.w3.org/TR/webdriver/#get-all-cookies
func (c *Client) Cookies(ctx context.Context) ([]Cookie, error) {
	route := fmt.Sprintf("session/%s/cookie", c.session.ID)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return nil, err
	}

	res := new(cookiesResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	return res.Value, nil
}

// CookiesDelete command is used to delete all cookies visible to the current page.
//
// https://www.w3.org/TR/webdriver/#delete-all-cookies
func (c *Client) CookiesDelete(ctx context.Context) error {
	route := fmt.Sprintf("session/%s/cookie", c.session.ID)

	req, err := c.prepare(http.MethodDelete, route, nil)
	if err != nil {
		return err
	}

	return c.do(ctx, req, nil)
}
