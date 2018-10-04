// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2017 Canonical Ltd
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3 as
// published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package restclient

import (
        "fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
)

const (
	versionURI       = "/v1"
	configurationURI = "/configuration"
)

var socketPath = "/run/snapd.socket"

// TransportClient operations executed by any client requesting server.
type TransportClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// RestClient defines client for rest api exposed by a unix socket
type RestClient struct {
	transportClient TransportClient
}

func newRestClient(client TransportClient) *RestClient {
	return &RestClient{transportClient: client}
}

func unixDialer(_, _ string) (net.Conn, error) {
	return net.Dial("unix", socketPath)
}

// DefaultRestClient created a RestClient object pointing to default socket path
func DefaultRestClient() *RestClient {
	return newRestClient(&http.Client{
		Transport: &http.Transport{
			Dial: unixDialer,
		},
	})
}

func (restClient *RestClient) SendHTTPRequest(uri string, method string, body io.Reader) (string, error) {
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return "", err
	}
	resp, err := restClient.transportClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	b, _ := ioutil.ReadAll(resp.Body)

	return string(b), nil
}

func (restclient *RestClient) Yeah(s string) {
	fmt.Println(s)
}

func (restClient *RestClient) SendHTTPRequestHeaders(uri string, method string, body io.Reader, headers map[string]string) (string, error) {
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return "", err
	}

	for key, value := range headers {
		req.Header.Set(key,value)
	}
	resp, err := restClient.transportClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)

	return string(b), nil
}
