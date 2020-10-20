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

// LocatorStrategy is an enumerated attribute deciding what technique should be used to search for elements in the current browsing context.
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

// ElementFind command is used to find an element in the current browsing context.
//
// https://www.w3.org/TR/webdriver/#find-element
func (c *client) ElementFind(ctx context.Context, by LocatorStrategy, v string) (ElementID, error) {
	if by == "" {
		return "", errors.New("locator strategy is required")
	}
	if v == "" {
		return "", errors.New("value is required")
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

	return "", errors.New("element is not found")
}

// ElementsFind command is used to find elements in the current browsing context.
//
// https://www.w3.org/TR/webdriver/#find-elements
func (c *client) ElementsFind(ctx context.Context, by LocatorStrategy, v string) ([]ElementID, error) {
	if by == "" {
		return nil, errors.New("locator strategy is required")
	}
	if v == "" {
		return nil, errors.New("value is required")
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
		return nil, errors.New("elements are not found")
	}

	ids := make([]ElementID, len(res.Value))

	for i, e := range res.Value {
		if id, ok := e[WebElementID]; ok {
			ids[i] = id
		}
	}

	return ids, nil
}

// ElementFindFrom command is used to find an element from element in the current browsing context.
//
// https://www.w3.org/TR/webdriver/#find-element-from-element
func (c *client) ElementFindFrom(ctx context.Context, eid ElementID, by LocatorStrategy, v string) (ElementID, error) {
	if eid == "" {
		return "", errors.New("element ID is required")
	}
	if by == "" {
		return "", errors.New("locator strategy is required")
	}
	if v == "" {
		return "", errors.New("value is required")
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

	return "", errors.New("element is not found")
}

// ElementsFindFrom command is used to find elements from element in the current browsing context.
//
// https://www.w3.org/TR/webdriver/#find-elements-from-element
func (c *client) ElementsFindFrom(ctx context.Context, eid ElementID, by LocatorStrategy, v string) ([]ElementID, error) {
	if eid == "" {
		return nil, errors.New("element ID is required")
	}
	if by == "" {
		return nil, errors.New("locator strategy is required")
	}
	if v == "" {
		return nil, errors.New("value is required")
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
		return nil, errors.New("elements are not found")
	}

	ids := make([]ElementID, len(res.Value))

	for i, e := range res.Value {
		if id, ok := e[WebElementID]; ok {
			ids[i] = id
		}
	}

	return ids, nil
}

// ElementClick command is used to click element.
//
// https://www.w3.org/TR/webdriver/#element-click
func (c *client) ElementClick(ctx context.Context, eid ElementID) error {
	if eid == "" {
		return errors.New("element ID is required")
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

// ElementClear command is used to click element.
//
// https://www.w3.org/TR/webdriver/#element-clear
func (c *client) ElementClear(ctx context.Context, eid ElementID) error {
	if eid == "" {
		return errors.New("element ID is required")
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

// ElementSendKeys command is used to send the provided keys to the element.
//
// https://www.w3.org/TR/webdriver/#element-send-keys
func (c *client) ElementSendKeys(ctx context.Context, eid ElementID, keys string) error {
	if eid == "" {
		return errors.New("element ID is required")
	}
	if keys == "" {
		return errors.New("keys is required")
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

// ElementAttribute command is used to get element attribute.
//
// https://www.w3.org/TR/webdriver/#get-element-attribute
func (c *client) ElementAttribute(ctx context.Context, eid ElementID, attr string) (string, error) {
	if eid == "" {
		return "", errors.New("element ID is required")
	}
	if attr == "" {
		return "", errors.New("attribute is required")
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

// ElementProperty command is used to get element property.
//
// https://www.w3.org/TR/webdriver/#get-element-property
func (c *client) ElementProperty(ctx context.Context, eid ElementID, prop string) (string, error) {
	if eid == "" {
		return "", errors.New("element ID is required")
	}
	if prop == "" {
		return "", errors.New("property is required")
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

// ElementCSSValue command is used to get element CSS property value.
//
// https://www.w3.org/TR/webdriver/#get-element-css-value
func (c *client) ElementCSSValue(ctx context.Context, eid ElementID, prop string) (string, error) {
	if eid == "" {
		return "", errors.New("element ID is required")
	}
	if prop == "" {
		return "", errors.New("CSS property is required")
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

// ElementText command is used to get element property.
//
// https://www.w3.org/TR/webdriver/#get-element-text
func (c *client) ElementText(ctx context.Context, eid ElementID) (string, error) {
	if eid == "" {
		return "", errors.New("element ID is required")
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

// ElementTagName command is used to get element tag name.
//
// https://www.w3.org/TR/webdriver/#get-element-tag-name
func (c *client) ElementTagName(ctx context.Context, eid ElementID) (string, error) {
	if eid == "" {
		return "", errors.New("element ID is required")
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
