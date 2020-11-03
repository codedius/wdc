// Package webdriver provides simple REST Client for a remote web driver server.
package wdc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var (
	ErrorSessionIsEmpty            = errors.New("session is empty")
	ErrorBaseURLIsEmpty            = errors.New("base URL is empty")
	ErrorLocatorStrategyIsRequired = errors.New("locator strategy is required")
	ErrorValueIsRequired           = errors.New("value is required")
	ErrorElementIsNotFound         = errors.New("element is not found")
	ErrorElementsAreNotFound       = errors.New("elements are not found")
	ErrorElementIDIsRequired       = errors.New("element ID is required")
	ErrorKeysAreRequired           = errors.New("keys are required")
	ErrorAttributeIsRequired       = errors.New("attribute is required")
	ErrorPropertyIsRequired        = errors.New("property is required")
	ErrorCSSPropertyIsRequired     = errors.New("CSS property is required")
	ErrorURLIsRequired             = errors.New("URL is required")
)

// Session data to make requests to the server.
type Session struct {
	// ID of a remote session
	ID string
	// URL of a web driver server
	URL string
}

// Client for a server API.
type Client struct {
	session *Session
	client  *http.Client
	url     *url.URL
}

// New returns a new web driver REST Client instance.
func New(s *Session) (*Client, error) {
	if s == nil {
		return nil, ErrorSessionIsEmpty
	}
	if s.URL == "" {
		return nil, ErrorBaseURLIsEmpty
	}

	s.URL = strings.TrimSuffix(s.URL, "/") + "/"

	httpc := http.DefaultClient

	u, err := url.Parse(s.URL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		session: s,
		client:  httpc,
		url:     u,
	}

	return c, nil
}

// prepare creates a server request.
func (c *Client) prepare(method string, path string, body io.Reader) (*http.Request, error) {
	p, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	u := c.url.ResolveReference(p)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// do sends a server request and returns server response.
//
// The provided ctx must be non-nil. If it is canceled or time out, ctx.Err() will be returned.
func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) error {
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		return err
	}
	defer safeclose(resp.Body)

	err = check(resp)
	if err != nil {
		return err
	}

	if v == nil {
		return nil
	}

	err = json.NewDecoder(resp.Body).Decode(v)
	if err == io.EOF {
		return nil // ignore EOF errors caused by empty response body
	}

	return err
}

// check checks the server response for errors.
func check(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	errResp := &ErrorResponse{}

	if data != nil {
		err = json.Unmarshal(data, errResp)
		if err != nil {
			return err
		}
	}

	return errResp
}

//
// ERRORS
//

type ErrorResponse struct {
	// SessionID is an ID of the WebDriver session.
	SessionID string `json:"sessionId"`
	// Value is an WebDriver error.
	Value ErrorValue `json:"value"`
}

// ErrorValue contains information about a failure of a command.
//
// https://www.w3.org/TR/webdriver/#handling-errors.
type ErrorValue struct {
	// Err contains a general error string provided by the server.
	Err string `json:"error"`
	// Message is a detailed, human-readable message specific to the failure.
	Message string `json:"message"`
	// Stacktrace may contain the server-side stacktrace message where the error occurred.
	Stacktrace string `json:"stacktrace,omitempty"`
	// StackTrace may contain the server-side stacktrace where the error occurred.
	StackTrace []ErrorStackTrace `json:"stackTrace,omitempty"`
}

type ErrorStackTrace struct {
	FileName   string `json:"fileName"`
	MethodName string `json:"methodName"`
	ClassName  string `json:"className"`
	LineNumber int    `json:"lineNumber"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s\n", e.Value.Err, e.Value.Message)
}

//
// UTILS
//

// safeclose is a convenient function for defer closing io.Closer c types.
func safeclose(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}

//
// RESPONSES
//

// value is a simplified response from the server.
type value struct {
	Value string `json:"value"`
}
