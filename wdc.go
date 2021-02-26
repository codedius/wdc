// Package wdc provides simple REST Client for a remote web driver server.
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

//
// TYPES
//

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

//
// MAIN
//

// New returns a new web driver REST Client instance.
func New(s *Session) (*Client, error) {
	if s == nil {
		return nil, errors.New("session is empty")
	}
	if s.URL == "" {
		return nil, errors.New("URL is empty")
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

	if errResp.Value.Err != "" {
		if err, ok := errs[errResp.Value.Err]; ok {
			return fmt.Errorf("%w: %s", err, errResp)
		}
		return errResp
	}

	// Support legacy status code to define error
	if errResp.Status != 0 {
		if err, ok := legacyErrs[errResp.Status]; ok {
			return fmt.Errorf("%w: %s", err, errResp)
		}
	}

	return errResp
}

//
// RESPONSES
//

// stringValue is a simplified string response from the server.
type stringValue struct {
	Value string `json:"value"`
}

// UnmarshalJSON first tries to unmarshal the value as string, then defaults to json.RawMessage
func (s *stringValue) UnmarshalJSON(bytes []byte) error {
	type alias struct {
		Value string
	}
	al := alias{}
	// try with string first (we need alias to avoid stack overflow)
	if err := json.Unmarshal(bytes, &al); err == nil {
		s.Value = al.Value
		return nil
	}

	// now try with raw json, is more forgiving
	type priv struct {
		Value json.RawMessage
	}
	p := priv{}
	err := json.Unmarshal(bytes, &p)
	if err != nil {
		return err
	}
	s.Value = string(p.Value)
	return nil
}

// boolValue is a simplified bool response from the server.
type boolValue struct {
	Value bool `json:"value"`
}

//
// ERRORS
//

var (
	ErrorElementClickIntercepted = errors.New("element click intercepted")
	ErrorElementNotInteractable  = errors.New("element not interactable")
	ErrorInsecureCertificate     = errors.New("insecure certificate")
	ErrorInvalidArgument         = errors.New("invalid argument")
	ErrorInvalidCookieDomain     = errors.New("invalid cookie domain")
	ErrorInvalidElementState     = errors.New("invalid element state")
	ErrorInvalidSelector         = errors.New("invalid selector")
	ErrorInvalidSessionID        = errors.New("invalid session id")
	ErrorJavaScriptError         = errors.New("javascript error")
	ErrorMoveTargetOutOfBounds   = errors.New("move target out of bounds")
	ErrorNoSuchAlert             = errors.New("no such alert")
	ErrorNoSuchCookie            = errors.New("no such cookie")
	ErrorNoSuchElement           = errors.New("no such element")
	ErrorNoSuchFrame             = errors.New("no such frame")
	ErrorNoSuchWindow            = errors.New("no such window")
	ErrorScriptTimeout           = errors.New("script timeout")
	ErrorSessionNotCreated       = errors.New("session not created")
	ErrorStaleElementReference   = errors.New("stale element reference")
	ErrorTimeout                 = errors.New("timeout")
	ErrorUnableToSetCookie       = errors.New("unable to set cookie")
	ErrorUnableToCaptureScreen   = errors.New("unable to capture screen")
	ErrorUnexpectedAlertOpen     = errors.New("unexpected alert open")
	ErrorUnknownCommand          = errors.New("unknown command")
	ErrorUnknownError            = errors.New("unknown error")
	ErrorUnknownMethod           = errors.New("unknown method")
	ErrorUnsupportedOperation    = errors.New("unsupported operation")
)

var errs = map[string]error{
	"element click intercepted": ErrorElementClickIntercepted,
	"element not interactable":  ErrorElementNotInteractable,
	"insecure certificate":      ErrorInsecureCertificate,
	"invalid argument":          ErrorInvalidArgument,
	"invalid cookie domain":     ErrorInvalidCookieDomain,
	"invalid element state":     ErrorInvalidElementState,
	"invalid selector":          ErrorInvalidSelector,
	"invalid session id":        ErrorInvalidSessionID,
	"javascript error":          ErrorJavaScriptError,
	"move target out of bounds": ErrorMoveTargetOutOfBounds,
	"no such alert":             ErrorNoSuchAlert,
	"no such cookie":            ErrorNoSuchCookie,
	"no such element":           ErrorNoSuchElement,
	"no such frame":             ErrorNoSuchFrame,
	"no such window":            ErrorNoSuchWindow,
	"script timeout":            ErrorScriptTimeout,
	"session not created":       ErrorSessionNotCreated,
	"stale element reference":   ErrorStaleElementReference,
	"timeout":                   ErrorTimeout,
	"unable to set cookie":      ErrorUnableToSetCookie,
	"unable to capture screen":  ErrorUnableToCaptureScreen,
	"unexpected alert open":     ErrorUnexpectedAlertOpen,
	"unknown command":           ErrorUnknownCommand,
	"unknown error":             ErrorUnknownError,
	"unknown method":            ErrorUnknownMethod,
	"unsupported operation":     ErrorUnsupportedOperation,
}

var legacyErrs = map[int]error{
	7:  ErrorNoSuchElement,
	8:  ErrorNoSuchFrame,
	9:  ErrorUnknownCommand,
	10: ErrorStaleElementReference,
	12: ErrorInvalidElementState,
	13: ErrorUnknownError,
	17: ErrorJavaScriptError,
	21: ErrorTimeout,
	23: ErrorNoSuchWindow,
	24: ErrorInvalidCookieDomain,
	25: ErrorUnableToSetCookie,
	26: ErrorUnexpectedAlertOpen,
	28: ErrorScriptTimeout,
	32: ErrorInvalidSelector,
	33: ErrorSessionNotCreated,
	34: ErrorMoveTargetOutOfBounds,
}

type ErrorResponse struct {
	// SessionID is an ID of the WebDriver session.
	SessionID string `json:"sessionId"`
	// Value is an WebDriver error.
	Value ErrorValue `json:"value"`
	// Status is a legacy response status code.
	Status int `json:"status"`
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
	return e.Value.Message
}

//
// UTILS
//

// safeclose is a convenient function for defer closing io.Closer c types.
func safeclose(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Panic(err)
	}
}
