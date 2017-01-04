// Copyright (c) 2017 LunaNode Hosting Inc. All right reserved.
// Use of this source code is governed by the MIT License. See LICENSE file.

package namesilo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func GetTestClient(t *testing.T, path string, response string) *Client {
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			t.Errorf("expected client to request path %s, but got path %s", path, r.URL.Path)
		}
		w.Write([]byte(response))
		go server.Close()
	}))
	return &Client{
		HTTPClient: &http.Client{},
		apiBaseURL: server.URL + "/",
	}
}

func TestDecodeResponse(t *testing.T) {
	var testResponse struct {
		Response string `xml:"response"`
	}
	rawResponse := `<namesilo>
	<reply>
		<code>300</code>
		<detail>success</detail>
		<response>test</response>
	</reply>
</namesilo>`
	if err := decodeResponse([]byte(rawResponse), &testResponse); err != nil {
		t.Fatalf("unexpected decode response error: %v", err)
	}
	if testResponse.Response != "test" {
		t.Fatalf("response is '%s' but expected 'test'", testResponse.Response)
	}
}
