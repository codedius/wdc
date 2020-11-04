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

// ElementFind command is used to find an element.
//
// https://www.w3.org/TR/webdriver/#find-element
func (c *Client) ElementFind(ctx context.Context, by LocatorStrategy, v string) (ElementID, error) {
	if by == "" {
		return "", ErrorLocatorStrategyIsRequired
	}
	if v == "" {
		return "", ErrorValueIsRequired
	}

	r := &elementRequest{
		Using: by,
		Value: v,
	}

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

	return "", ErrorElementIsNotFound
}

// ElementsFind command is used to find elements.
//
// https://www.w3.org/TR/webdriver/#find-elements
func (c *Client) ElementsFind(ctx context.Context, by LocatorStrategy, v string) ([]ElementID, error) {
	if by == "" {
		return nil, ErrorLocatorStrategyIsRequired
	}
	if v == "" {
		return nil, ErrorValueIsRequired
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
		return nil, ErrorElementsAreNotFound
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

// ElementFindFrom command is used to find an element from element.
//
// https://www.w3.org/TR/webdriver/#find-element-from-element
func (c *Client) ElementFindFrom(ctx context.Context, eid ElementID, by LocatorStrategy, v string) (ElementID, error) {
	if eid == "" {
		return "", ErrorElementIDIsRequired
	}
	if by == "" {
		return "", ErrorLocatorStrategyIsRequired
	}
	if v == "" {
		return "", ErrorValueIsRequired
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

	return "", ErrorElementIsNotFound
}

// ElementsFindFrom command is used to find elements from an element.
//
// https://www.w3.org/TR/webdriver/#find-elements-from-element
func (c *Client) ElementsFindFrom(ctx context.Context, eid ElementID, by LocatorStrategy, v string) ([]ElementID, error) {
	if eid == "" {
		return nil, ErrorElementIDIsRequired
	}
	if by == "" {
		return nil, ErrorLocatorStrategyIsRequired
	}
	if v == "" {
		return nil, ErrorValueIsRequired
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
		return nil, ErrorElementsAreNotFound
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

// ElementClick command is used to click on an element.
//
// https://www.w3.org/TR/webdriver/#element-click
func (c *Client) ElementClick(ctx context.Context, eid ElementID) error {
	if eid == "" {
		return ErrorElementIDIsRequired
	}

	route := fmt.Sprintf("session/%s/element/%s/click", c.session.ID, eid)

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

// ElementClear command is used to clear an input or textarea element.
//
// https://www.w3.org/TR/webdriver/#element-clear
func (c *Client) ElementClear(ctx context.Context, eid ElementID) error {
	if eid == "" {
		return ErrorElementIDIsRequired
	}

	route := fmt.Sprintf("session/%s/element/%s/clear", c.session.ID, eid)

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

// ElementSendKeys command is used to send provided keys to an element.
//
// https://www.w3.org/TR/webdriver/#element-send-keys
func (c *Client) ElementSendKeys(ctx context.Context, eid ElementID, keys string) error {
	if eid == "" {
		return ErrorElementIDIsRequired
	}
	if keys == "" {
		return ErrorKeysAreRequired
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

	err = c.do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}

// ElementAttribute command is used to get the value of an element's attribute.
//
// https://www.w3.org/TR/webdriver/#get-element-attribute
func (c *Client) ElementAttribute(ctx context.Context, eid ElementID, attr string) (string, error) {
	if eid == "" {
		return "", ErrorElementIDIsRequired
	}
	if attr == "" {
		return "", ErrorAttributeIsRequired
	}

	route := fmt.Sprintf("session/%s/element/%s/attribute/%s", c.session.ID, eid, attr)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return "", err
	}

	res := new(value)

	err = c.do(ctx, req, res)
	if err != nil {
		return "", err
	}

	return res.Value, nil
}

// ElementProperty command is used to get the value of an element's property.
//
// https://www.w3.org/TR/webdriver/#get-element-property
func (c *Client) ElementProperty(ctx context.Context, eid ElementID, prop string) (string, error) {
	if eid == "" {
		return "", ErrorElementIDIsRequired
	}
	if prop == "" {
		return "", ErrorPropertyIsRequired
	}

	route := fmt.Sprintf("session/%s/element/%s/property/%s", c.session.ID, eid, prop)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return "", err
	}

	res := new(value)

	err = c.do(ctx, req, res)
	if err != nil {
		return "", err
	}

	return res.Value, nil
}

// ElementCSSValue command is used to get the value of an element's CSS property.
//
// https://www.w3.org/TR/webdriver/#get-element-css-value
func (c *Client) ElementCSSValue(ctx context.Context, eid ElementID, prop string) (string, error) {
	if eid == "" {
		return "", ErrorElementIDIsRequired
	}
	if prop == "" {
		return "", ErrorCSSPropertyIsRequired
	}

	route := fmt.Sprintf("session/%s/element/%s/css/%s", c.session.ID, eid, prop)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return "", err
	}

	res := new(value)

	err = c.do(ctx, req, res)
	if err != nil {
		return "", err
	}

	return res.Value, nil
}

// ElementText command is used to get the text of an element.
//
// https://www.w3.org/TR/webdriver/#get-element-text
func (c *Client) ElementText(ctx context.Context, eid ElementID) (string, error) {
	if eid == "" {
		return "", ErrorElementIDIsRequired
	}

	route := fmt.Sprintf("session/%s/element/%s/text", c.session.ID, eid)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return "", err
	}

	res := new(value)

	err = c.do(ctx, req, res)
	if err != nil {
		return "", err
	}

	return res.Value, nil
}

// ElementTagName command is used to get a tag name of an element.
//
// https://www.w3.org/TR/webdriver/#get-element-tag-name
func (c *Client) ElementTagName(ctx context.Context, eid ElementID) (string, error) {
	if eid == "" {
		return "", ErrorElementIDIsRequired
	}

	route := fmt.Sprintf("session/%s/element/%s/name", c.session.ID, eid)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return "", err
	}

	res := new(value)

	err = c.do(ctx, req, res)
	if err != nil {
		return "", err
	}

	return res.Value, nil
}

// ElementScreenshot command is used to take a screenshot of an element.
//
// https://www.w3.org/TR/webdriver/#take-element-screenshot
func (c *Client) ElementScreenshot(ctx context.Context, eid ElementID) (string, error) {
	if eid == "" {
		return "", ErrorElementIDIsRequired
	}

	route := fmt.Sprintf("session/%s/element/%s/screenshot", c.session.ID, eid)

	req, err := c.prepare(http.MethodGet, route, nil)
	if err != nil {
		return "", err
	}

	res := new(value)

	err = c.do(ctx, req, res)
	if err != nil {
		return "", err
	}

	return res.Value, nil
}
