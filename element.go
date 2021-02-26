package wdc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

//
// TYPES
//

// LocatorStrategy is an enumerated attribute deciding what technique should be used to search for elements.
//
// https://www.w3.org/TR/webdriver/#locator-strategies
type LocatorStrategy string

const (
	BySelector        LocatorStrategy = "css selector"
	ByLinkText        LocatorStrategy = "link text"
	ByPartialLinkText LocatorStrategy = "partial link text"
	ByTagName         LocatorStrategy = "tag name"
	ByXPath           LocatorStrategy = "xpath"
)

type WebElement struct {
	ID        WebElementID
	Reference WebElementReference
}

// WebElementID is the string constant defined by the W3C.
//
// https://www.w3.org/TR/webdriver/#elements
type WebElementID string

const (
	WebElementIDW3C    WebElementID = "element-6066-11e4-a52e-4f735466cecf"
	WebElementIDLegacy WebElementID = "ELEMENT"
)

// WebElementReference represents a reference of a web element.
type WebElementReference string

//
// REQUESTS
//

type elementRequest struct {
	Using LocatorStrategy `json:"using"`
	Value string          `json:"value"`
}

type elementSendKeysRequest struct {
	Text string `json:"text"`
}

type elementSendKeysLegacyRequest struct {
	Value []string `json:"value"`
}

//
// RESPONSES
//

type elementResponse struct {
	Value map[WebElementID]WebElementReference `json:"value"`
}

type elementsResponse struct {
	Value []map[WebElementID]WebElementReference `json:"value"`
}

//
// METHODS
//

// ElementFind command is used to find an element by locator strategy with value v.
//
// https://www.w3.org/TR/webdriver/#find-element
func (c *Client) ElementFind(ctx context.Context, by LocatorStrategy, v string) (WebElement, error) {
	if by == "" {
		return WebElement{}, errors.New("locator strategy is empty")
	}
	if v == "" {
		return WebElement{}, errors.New("value is empty")
	}

	r := &elementRequest{Using: by, Value: v}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return WebElement{}, err
	}

	route := fmt.Sprintf("session/%s/element", c.session.ID)

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return WebElement{}, err
	}

	res := new(elementResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return WebElement{}, err
	}

	if ref, ok := res.Value[WebElementIDW3C]; ok {
		return WebElement{ID: WebElementIDW3C, Reference: ref}, nil
	}
	if ref, ok := res.Value[WebElementIDLegacy]; ok {
		return WebElement{ID: WebElementIDLegacy, Reference: ref}, nil
	}

	return WebElement{}, ErrorNoSuchElement
}

// ElementsWaitForUnload sets interval i and amount of time t the driver should wait to determine if an element e is undefined.
func (c *Client) ElementWaitForUndefined(ctx context.Context, by LocatorStrategy, v string, i time.Duration, t time.Duration) error {
	r := &elementRequest{Using: by, Value: v}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return err
	}

	route := fmt.Sprintf("session/%s/element", c.session.ID)

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return err
	}

	res := new(elementResponse)

	start := time.Now()

	for {
		err = c.do(ctx, req, res)
		if err != nil && errors.Is(err, ErrorNoSuchElement) {
			return nil
		}
		if err != nil {
			return err
		}
		if elapsed := time.Since(start); elapsed > t {
			return fmt.Errorf("timeout after %v", elapsed)
		}
		time.Sleep(i)
	}
}

// ElementFindShadowDOM command is used to find a shadow root of element e.
func (c *Client) ElementFindShadowDOM(ctx context.Context, e WebElement) (WebElement, error) {
	if e.ID == "" {
		return WebElement{}, errors.New("element web ID is empty")
	}
	if e.Reference == "" {
		return WebElement{}, errors.New("element ID is empty")
	}

	var args []interface{}
	args = append(args, map[WebElementID]WebElementReference{
		e.ID: e.Reference,
	})

	r := &scriptRequest{Script: "return arguments[0].shadowRoot", Args: args}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return WebElement{}, err
	}

	route := fmt.Sprintf("session/%s/execute/sync", c.session.ID)

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return WebElement{}, err
	}

	res := new(elementResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return WebElement{}, err
	}

	if ref, ok := res.Value[WebElementIDW3C]; ok {
		return WebElement{ID: WebElementIDW3C, Reference: ref}, nil
	}
	if ref, ok := res.Value[WebElementIDLegacy]; ok {
		return WebElement{ID: WebElementIDLegacy, Reference: ref}, nil
	}

	return WebElement{}, ErrorNoSuchElement
}

// ElementFindShadowDOMLegacy command is used to find a shadow root of element e.
func (c *Client) ElementFindShadowDOMLegacy(ctx context.Context, e WebElement) (WebElement, error) {
	if e.ID == "" {
		return WebElement{}, errors.New("element web ID is empty")
	}
	if e.Reference == "" {
		return WebElement{}, errors.New("element ID is empty")
	}

	var args []interface{}
	args = append(args, map[WebElementID]WebElementReference{
		e.ID: e.Reference,
	})

	r := &scriptRequest{Script: "return arguments[0].shadowRoot", Args: args}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return WebElement{}, err
	}

	route := fmt.Sprintf("session/%s/execute", c.session.ID)

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return WebElement{}, err
	}

	res := new(elementResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return WebElement{}, err
	}

	if ref, ok := res.Value[WebElementIDW3C]; ok {
		return WebElement{ID: WebElementIDW3C, Reference: ref}, nil
	}
	if ref, ok := res.Value[WebElementIDLegacy]; ok {
		return WebElement{ID: WebElementIDLegacy, Reference: ref}, nil
	}

	return WebElement{}, ErrorNoSuchElement
}

// ElementsFind command is used to find elements by locator strategy with value v.
//
// https://www.w3.org/TR/webdriver/#find-elements
func (c *Client) ElementsFind(ctx context.Context, by LocatorStrategy, v string) ([]WebElement, error) {
	if by == "" {
		return nil, errors.New("locator strategy is empty")
	}
	if v == "" {
		return nil, errors.New("value is empty")
	}

	r := &elementRequest{Using: by, Value: v}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return nil, err
	}

	route := fmt.Sprintf("session/%s/elements", c.session.ID)

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return nil, err
	}

	res := new(elementsResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	if len(res.Value) == 0 {
		return nil, ErrorNoSuchElement
	}

	elems := make([]WebElement, len(res.Value))

	for i, e := range res.Value {
		if ref, ok := e[WebElementIDW3C]; ok {
			elems[i] = WebElement{ID: WebElementIDW3C, Reference: ref}
			continue
		}
		if ref, ok := e[WebElementIDLegacy]; ok {
			elems[i] = WebElement{ID: WebElementIDLegacy, Reference: ref}
		}
	}

	return elems, nil
}

// ElementFindFrom command is used to find an element by locator strategy with value v from element e.
//
// https://www.w3.org/TR/webdriver/#find-element-from-element
func (c *Client) ElementFindFrom(ctx context.Context, e WebElement, by LocatorStrategy, v string) (WebElement, error) {
	if e.Reference == "" {
		return WebElement{}, errors.New("element is empty")
	}
	if by == "" {
		return WebElement{}, errors.New("locator strategy is empty")
	}
	if v == "" {
		return WebElement{}, errors.New("value is empty")
	}

	r := &elementRequest{Using: by, Value: v}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return WebElement{}, err
	}

	route := fmt.Sprintf("session/%s/element/%s/element", c.session.ID, e.Reference)

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return WebElement{}, err
	}

	res := new(elementResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return WebElement{}, err
	}

	if ref, ok := res.Value[WebElementIDW3C]; ok {
		return WebElement{ID: WebElementIDW3C, Reference: ref}, nil
	}
	if ref, ok := res.Value[WebElementIDLegacy]; ok {
		return WebElement{ID: WebElementIDLegacy, Reference: ref}, nil
	}

	return WebElement{}, ErrorNoSuchElement
}

// ElementsFindFrom command is used to find elements by locator strategy with value v from an element e.
//
// https://www.w3.org/TR/webdriver/#find-elements-from-element
func (c *Client) ElementsFindFrom(ctx context.Context, e WebElement, by LocatorStrategy, v string) ([]WebElement, error) {
	if e.Reference == "" {
		return nil, errors.New("element is empty")
	}
	if by == "" {
		return nil, errors.New("locator strategy is empty")
	}
	if v == "" {
		return nil, errors.New("value is empty")
	}

	r := &elementRequest{Using: by, Value: v}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return nil, err
	}

	route := fmt.Sprintf("session/%s/element/%s/elements", c.session.ID, e.Reference)

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return nil, err
	}

	res := new(elementsResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	if len(res.Value) == 0 {
		return nil, ErrorNoSuchElement
	}

	elems := make([]WebElement, len(res.Value))

	for i, e := range res.Value {
		if ref, ok := e[WebElementIDW3C]; ok {
			elems[i] = WebElement{ID: WebElementIDW3C, Reference: ref}
			continue
		}
		if ref, ok := e[WebElementIDLegacy]; ok {
			elems[i] = WebElement{ID: WebElementIDLegacy, Reference: ref}
		}
	}

	return elems, nil
}

// ElementClick command is used to click on an element e.
//
// https://www.w3.org/TR/webdriver/#element-click
func (c *Client) ElementClick(ctx context.Context, e WebElement) error {
	if e.Reference == "" {
		return errors.New("element is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/click", c.session.ID, e.Reference)

	req, err := c.prepare(http.MethodPost, route, nil)
	if err != nil {
		return err
	}

	return c.do(ctx, req, nil)
}

// ElementClear command is used to clear an input or textarea element e.
//
// https://www.w3.org/TR/webdriver/#element-clear
func (c *Client) ElementClear(ctx context.Context, e WebElement) error {
	if e.Reference == "" {
		return errors.New("element is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/clear", c.session.ID, e.Reference)

	req, err := c.prepare(http.MethodPost, route, nil)
	if err != nil {
		return err
	}

	return c.do(ctx, req, nil)
}

// ElementSendKeys command is used to send provided keys to an element e.
//
// https://www.w3.org/TR/webdriver/#element-send-keys
func (c *Client) ElementSendKeys(ctx context.Context, e WebElement, keys string) error {
	if e.Reference == "" {
		return errors.New("element is empty")
	}
	if keys == "" {
		return errors.New("keys are empty")
	}

	r := &elementSendKeysRequest{Text: keys}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return err
	}

	route := fmt.Sprintf("session/%s/element/%s/value", c.session.ID, e.Reference)

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return err
	}

	res := new(elementResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return err
	}

	return nil
}

// ElementSendKeysLegacy command is used to send provided keys to an element e.
//
// https://github.com/SeleniumHQ/selenium/wiki/JsonWireProtocol#sessionsessionidelementidvalue
func (c *Client) ElementSendKeysLegacy(ctx context.Context, e WebElement, keys string) error {
	if e.Reference == "" {
		return errors.New("element is empty")
	}
	if len(keys) == 0 {
		return errors.New("keys are empty")
	}

	var v []string
	for _, k := range keys {
		v = append(v, string(k))
	}

	r := &elementSendKeysLegacyRequest{Value: v}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return err
	}

	route := fmt.Sprintf("session/%s/element/%s/value", c.session.ID, e.Reference)

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return err
	}

	res := new(elementsResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return err
	}

	return nil
}

// ElementAttribute command is used to get the attribute attr value of an element e.
//
// https://www.w3.org/TR/webdriver/#get-element-attribute
func (c *Client) ElementAttribute(ctx context.Context, e WebElement, attr string) (string, error) {
	if e.Reference == "" {
		return "", errors.New("element is empty")
	}
	if attr == "" {
		return "", errors.New("attribute is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/attribute/%s", c.session.ID, e.Reference, attr)

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

// ElementProperty command is used to get the property prop value of an element e.
//
// https://www.w3.org/TR/webdriver/#get-element-property
func (c *Client) ElementProperty(ctx context.Context, e WebElement, prop string) (string, error) {
	if e.Reference == "" {
		return "", errors.New("element is empty")
	}
	if prop == "" {
		return "", errors.New("property is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/property/%s", c.session.ID, e.Reference, prop)

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

// ElementCSSValue command is used to get the CSS property prop value of an element e.
//
// https://www.w3.org/TR/webdriver/#get-element-css-value
func (c *Client) ElementCSSValue(ctx context.Context, e WebElement, prop string) (string, error) {
	if e.Reference == "" {
		return "", errors.New("element is empty")
	}
	if prop == "" {
		return "", errors.New("CSS property is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/css/%s", c.session.ID, e.Reference, prop)

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

// ElementText command is used to get the text of an element e.
//
// https://www.w3.org/TR/webdriver/#get-element-text
func (c *Client) ElementText(ctx context.Context, e WebElement) (string, error) {
	if e.Reference == "" {
		return "", errors.New("element is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/text", c.session.ID, e.Reference)

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

// ElementWaitForEnabled sets interval i and amount of time t the driver should wait to determine if an element e is filled.
func (c *Client) ElementWaitForText(ctx context.Context, e WebElement, i time.Duration, t time.Duration) (string, error) {
	if e.Reference == "" {
		return "", errors.New("element is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/text", c.session.ID, e.Reference)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return "", err
	}

	res := new(stringValue)

	start := time.Now()

	for {
		err = c.do(ctx, req, res)
		if err != nil {
			return "", err
		}
		if res.Value != "" {
			return res.Value, nil
		}

		if elapsed := time.Since(start); elapsed > t {
			return "", fmt.Errorf("timeout after %v", elapsed)
		}
		time.Sleep(i)
	}
}

// ElementTagName command is used to get a tag name of an element e.
//
// https://www.w3.org/TR/webdriver/#get-element-tag-name
func (c *Client) ElementTagName(ctx context.Context, e WebElement) (string, error) {
	if e.Reference == "" {
		return "", errors.New("element is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/name", c.session.ID, e.Reference)

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

// ElementScreenshot command is used to take a screenshot of an element e.
//
// https://www.w3.org/TR/webdriver/#take-element-screenshot
func (c *Client) ElementScreenshot(ctx context.Context, e WebElement) (string, error) {
	if e.Reference == "" {
		return "", errors.New("element ID is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/screenshot", c.session.ID, e.Reference)

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

// ElementIsSelected command is used to determine if option/input/checkbox/radiobutton element e is currently selected.
//
// https://www.w3.org/TR/webdriver/#is-element-selected
func (c *Client) ElementIsSelected(ctx context.Context, e WebElement) (bool, error) {
	if e.Reference == "" {
		return false, errors.New("element is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/selected", c.session.ID, e.Reference)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return false, err
	}

	res := new(boolValue)

	err = c.do(ctx, req, res)
	if err != nil {
		return false, err
	}

	return res.Value, nil
}

// ElementIsEnabled command is used to determine if an element e is currently enabled.
//
// https://www.w3.org/TR/webdriver/#is-element-enabled
func (c *Client) ElementIsEnabled(ctx context.Context, e WebElement) (bool, error) {
	if e.Reference == "" {
		return false, errors.New("element is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/enabled", c.session.ID, e.Reference)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return false, err
	}

	res := new(boolValue)

	err = c.do(ctx, req, res)
	if err != nil {
		return false, err
	}

	return res.Value, nil
}

// ElementWaitForEnabled sets interval i and amount of time t the driver should wait to determine if an element e is enabled.
func (c *Client) ElementWaitForEnabled(ctx context.Context, e WebElement, i time.Duration, t time.Duration) error {
	if e.Reference == "" {
		return errors.New("element is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/enabled", c.session.ID, e.Reference)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return err
	}

	res := new(boolValue)

	start := time.Now()

	for {
		err = c.do(ctx, req, res)
		if err != nil {
			return err
		}
		if res.Value == true {
			return nil
		}

		if elapsed := time.Since(start); elapsed > t {
			return fmt.Errorf("timeout after %v", elapsed)
		}
		time.Sleep(i)
	}
}

// ElementIsDisplayed command is used to determine if an element e is currently displayed.
//
// https://www.w3.org/TR/webdriver/#element-displayedness
func (c *Client) ElementIsDisplayed(ctx context.Context, e WebElement) (bool, error) {
	if e.Reference == "" {
		return false, errors.New("element is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/displayed", c.session.ID, e.Reference)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return false, err
	}

	res := new(boolValue)

	err = c.do(ctx, req, res)
	if err != nil {
		return false, err
	}

	return res.Value, nil
}

// ElementWaitForEnabled sets interval i and amount of time t the driver should wait to determine if an element e is displayed.
func (c *Client) ElementWaitForDisplayed(ctx context.Context, e WebElement, i time.Duration, t time.Duration) error {
	if e.Reference == "" {
		return errors.New("element is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/displayed", c.session.ID, e.Reference)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return err
	}

	res := new(boolValue)

	start := time.Now()

	for {
		err = c.do(ctx, req, res)
		if err != nil {
			return err
		}
		if res.Value == true {
			return nil
		}

		if elapsed := time.Since(start); elapsed > t {
			return fmt.Errorf("timeout after %v", elapsed)
		}
		time.Sleep(i)
	}
}
