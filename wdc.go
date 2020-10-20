// Package webdriver provides simple REST client for a remote web driver server.
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

// Session data to make requests to the server.
type Session struct {
	// ID of a remote session
	ID string
	// URL of a web driver server
	URL string
}

// client for a server API.
type client struct {
	session *Session
	client  *http.Client
	url     *url.URL
}

// New returns a new web driver REST client instance.
func New(s *Session) (*client, error) {
	if s == nil {
		return nil, errors.New("session is empty")
	}
	if s.URL == "" {
		return nil, errors.New("base URL is empty")
	}

	s.URL = strings.TrimSuffix(s.URL, "/") + "/"

	httpc := http.DefaultClient

	u, err := url.Parse(s.URL)
	if err != nil {
		return nil, err
	}

	c := &client{
		session: s,
		client:  httpc,
		url:     u,
	}

	return c, nil
}

// prepare creates a server request.
func (c *client) prepare(method string, path string, body io.Reader) (*http.Request, error) {
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
func (c *client) do(ctx context.Context, req *http.Request, v interface{}) error {
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
