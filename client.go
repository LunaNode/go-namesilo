// Copyright (c) 2017 LunaNode Hosting Inc. All right reserved.
// Use of this source code is governed by the MIT License. See LICENSE file.

package namesilo

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const API_BASE_URL string = "https://www.namesilo.com/api/"
const API_VERSION = "1"

// NameSilo API client.
type Client struct {
	// The NameSilo API key.
	APIKey string

	// An HTTP client to perform API requests.
	HTTPClient *http.Client

	apiBaseURL string
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey: apiKey,
		HTTPClient: &http.Client{},
	}
}

func (client *Client) request(operation string, params map[string]string, responseTarget interface{}) error {
	// build URL
	queryValues := url.Values{}
	queryValues.Set("version", API_VERSION)
	queryValues.Set("type", "xml")
	queryValues.Set("key", client.APIKey)
	if params != nil {
		for k, v := range params {
			queryValues.Set(k, v)
		}
	}

	var baseURL string
	if client.apiBaseURL != "" {
		baseURL = client.apiBaseURL
	} else {
		baseURL = API_BASE_URL
	}
	apiURL := baseURL + operation + "?" + queryValues.Encode()

	// perform request
	response, err := client.HTTPClient.Get(apiURL)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading HTTP response: %v", err)
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("error %d: %s", response.StatusCode, string(contents))
	}

	if err := decodeResponse(contents, responseTarget); err != nil {
		return err
	}

	return nil
}

type genericResponse struct {
	XMLName xml.Name `xml:"namesilo"`
	Reply struct {
		XMLName xml.Name `xml:"reply"`
		Code int `xml:"code"`
		Detail string `xml:"detail"`
		InnerXML []byte `xml:",innerxml"`
	}
}

func (generic genericResponse) InnerXML() []byte {
	return []byte("<reply>" + string(generic.Reply.InnerXML) + "</reply>")
}

func decodeResponse(bytes []byte, responseTarget interface{}) error {
	var generic genericResponse
	if err := xml.Unmarshal(bytes, &generic); err != nil {
		return fmt.Errorf("invalid response: %v", err)
	}

	// reply codes 300 to 302 indicate success
	if generic.Reply.Code < 300 || generic.Reply.Code > 302 {
		return fmt.Errorf("error %d: %s", generic.Reply.Code, generic.Reply.Detail)
	}

	if responseTarget != nil {
		if err := xml.Unmarshal(generic.InnerXML(), responseTarget); err != nil {
			return fmt.Errorf("error decoding reply: %v", err)
		}
	}

	return nil
}
