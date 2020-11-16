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

// WebElementID is the string constant defined by the W3C.
//
// https://www.w3.org/TR/webdriver/#elements
const WebElementID = "element-6066-11e4-a52e-4f735466cecf"

// WebElementLegacyID is the legacy constant.
const WebElementLegacyID = "ELEMENT"

// ElementID represents an id of a web element.
type ElementID string

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
	Value []rune `json:"value"`
}

//
// RESPONSES
//

type elementResponse struct {
	Value map[string]ElementID `json:"value"`
}

type elementsResponse struct {
	Value []map[string]ElementID `json:"value"`
}

//
// METHODS
//

// ElementFind command is used to find an element by locator strategy with value v.
//
// https://www.w3.org/TR/webdriver/#find-element
func (c *Client) ElementFind(ctx context.Context, by LocatorStrategy, v string) (ElementID, error) {
	if by == "" {
		return "", errors.New("locator strategy is empty")
	}
	if v == "" {
		return "", errors.New("value is empty")
	}

	r := &elementRequest{Using: by, Value: v}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return "", err
	}

	route := fmt.Sprintf("session/%s/element", c.session.ID)

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return "", err
	}

	res := new(elementResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return "", err
	}

	if id, ok := res.Value[WebElementID]; ok {
		return id, nil
	}
	if id, ok := res.Value[WebElementLegacyID]; ok {
		return id, nil
	}

	return "", ErrorNoSuchElement
}

// ElementsFind command is used to find elements by locator strategy with value v.
//
// https://www.w3.org/TR/webdriver/#find-elements
func (c *Client) ElementsFind(ctx context.Context, by LocatorStrategy, v string) ([]ElementID, error) {
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

	ids := make([]ElementID, len(res.Value))

	for i, e := range res.Value {
		if id, ok := e[WebElementID]; ok {
			ids[i] = id
			continue
		}
		if id, ok := e[WebElementLegacyID]; ok {
			ids[i] = id
		}
	}

	return ids, nil
}

// ElementFindFrom command is used to find an element by locator strategy with value v from element with ID eid.
//
// https://www.w3.org/TR/webdriver/#find-element-from-element
func (c *Client) ElementFindFrom(ctx context.Context, eid ElementID, by LocatorStrategy, v string) (ElementID, error) {
	if eid == "" {
		return "", errors.New("element ID is empty")
	}
	if by == "" {
		return "", errors.New("locator strategy is empty")
	}
	if v == "" {
		return "", errors.New("value is empty")
	}

	r := &elementRequest{Using: by, Value: v}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return "", err
	}

	route := fmt.Sprintf("session/%s/element/%s/element", c.session.ID, eid)

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return "", err
	}

	res := new(elementResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return "", err
	}

	if id, ok := res.Value[WebElementID]; ok {
		return id, nil
	}
	if id, ok := res.Value[WebElementLegacyID]; ok {
		return id, nil
	}

	return "", ErrorNoSuchElement
}

// ElementsFindFrom command is used to find elements by locator strategy with value v from an element with ID eid.
//
// https://www.w3.org/TR/webdriver/#find-elements-from-element
func (c *Client) ElementsFindFrom(ctx context.Context, eid ElementID, by LocatorStrategy, v string) ([]ElementID, error) {
	if eid == "" {
		return nil, errors.New("element ID is empty")
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

	route := fmt.Sprintf("session/%s/element/%s/elements", c.session.ID, eid)

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

	ids := make([]ElementID, len(res.Value))

	for i, e := range res.Value {
		if id, ok := e[WebElementID]; ok {
			ids[i] = id
			continue
		}
		if id, ok := e[WebElementLegacyID]; ok {
			ids[i] = id
		}
	}

	return ids, nil
}

// ElementClick command is used to click on an element with ID eid.
//
// https://www.w3.org/TR/webdriver/#element-click
func (c *Client) ElementClick(ctx context.Context, eid ElementID) error {
	if eid == "" {
		return errors.New("element ID is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/click", c.session.ID, eid)

	req, err := c.prepare(http.MethodPost, route, nil)
	if err != nil {
		return err
	}

	return c.do(ctx, req, nil)
}

// ElementClear command is used to clear an input or textarea element with ID eid.
//
// https://www.w3.org/TR/webdriver/#element-clear
func (c *Client) ElementClear(ctx context.Context, eid ElementID) error {
	if eid == "" {
		return errors.New("element ID is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/clear", c.session.ID, eid)

	req, err := c.prepare(http.MethodPost, route, nil)
	if err != nil {
		return err
	}

	return c.do(ctx, req, nil)
}

// ElementSendKeys command is used to send provided keys to an element with ID eid.
//
// https://www.w3.org/TR/webdriver/#element-send-keys
func (c *Client) ElementSendKeys(ctx context.Context, eid ElementID, keys string) error {
	if eid == "" {
		return errors.New("element ID is empty")
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

	route := fmt.Sprintf("session/%s/element/%s/value", c.session.ID, eid)

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return err
	}

	res := new(elementsResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return err
	}

	return c.do(ctx, req, nil)
}

// ElementSendKeys command is used to send provided keys to an element with ID eid.
//
// https://www.w3.org/TR/webdriver/#element-send-keys
func (c *Client) ElementSendKeysLegacy(ctx context.Context, eid ElementID, keys []rune) error {
	if eid == "" {
		return errors.New("element ID is empty")
	}
	if len(keys) == 0 {
		return errors.New("keys are empty")
	}

	r := &elementSendKeysLegacyRequest{Value: keys}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return err
	}

	route := fmt.Sprintf("session/%s/element/%s/value", c.session.ID, eid)

	req, err := c.prepare(http.MethodPost, route, b)
	if err != nil {
		return err
	}

	res := new(elementsResponse)

	err = c.do(ctx, req, res)
	if err != nil {
		return err
	}

	return c.do(ctx, req, nil)
}

// ElementAttribute command is used to get the attribute attr value of an element with ID eid.
//
// https://www.w3.org/TR/webdriver/#get-element-attribute
func (c *Client) ElementAttribute(ctx context.Context, eid ElementID, attr string) (string, error) {
	if eid == "" {
		return "", errors.New("element ID is empty")
	}
	if attr == "" {
		return "", errors.New("attribute is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/attribute/%s", c.session.ID, eid, attr)

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

// ElementProperty command is used to get the property prop value of an element with ID eid.
//
// https://www.w3.org/TR/webdriver/#get-element-property
func (c *Client) ElementProperty(ctx context.Context, eid ElementID, prop string) (string, error) {
	if eid == "" {
		return "", errors.New("element ID is empty")
	}
	if prop == "" {
		return "", errors.New("property is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/property/%s", c.session.ID, eid, prop)

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

// ElementCSSValue command is used to get the CSS property prop value of an element with ID eid.
//
// https://www.w3.org/TR/webdriver/#get-element-css-value
func (c *Client) ElementCSSValue(ctx context.Context, eid ElementID, prop string) (string, error) {
	if eid == "" {
		return "", errors.New("element ID is empty")
	}
	if prop == "" {
		return "", errors.New("CSS property is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/css/%s", c.session.ID, eid, prop)

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

// ElementText command is used to get the text of an element with ID eid.
//
// https://www.w3.org/TR/webdriver/#get-element-text
func (c *Client) ElementText(ctx context.Context, eid ElementID) (string, error) {
	if eid == "" {
		return "", errors.New("element ID is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/text", c.session.ID, eid)

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

// ElementTagName command is used to get a tag name of an element with ID eid.
//
// https://www.w3.org/TR/webdriver/#get-element-tag-name
func (c *Client) ElementTagName(ctx context.Context, eid ElementID) (string, error) {
	if eid == "" {
		return "", errors.New("element ID is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/name", c.session.ID, eid)

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

// ElementScreenshot command is used to take a screenshot of an element with ID eid.
//
// https://www.w3.org/TR/webdriver/#take-element-screenshot
func (c *Client) ElementScreenshot(ctx context.Context, eid ElementID) (string, error) {
	if eid == "" {
		return "", errors.New("element ID is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/screenshot", c.session.ID, eid)

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

// ElementIsSelected command is used to determine if option/input/checkbox/radiobutton element with ID eid is currently selected.
//
// https://www.w3.org/TR/webdriver/#is-element-selected
func (c *Client) ElementIsSelected(ctx context.Context, eid ElementID) (bool, error) {
	if eid == "" {
		return false, errors.New("element ID is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/selected", c.session.ID, eid)

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

// ElementIsEnabled command is used to determine if an element with ID eid is currently enabled.
//
// https://www.w3.org/TR/webdriver/#is-element-enabled
func (c *Client) ElementIsEnabled(ctx context.Context, eid ElementID) (bool, error) {
	if eid == "" {
		return false, errors.New("element ID is empty")
	}

	route := fmt.Sprintf("session/%s/element/%s/enabled", c.session.ID, eid)

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
