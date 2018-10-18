// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 Canonical Ltd
//
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


package main

import (
    "bytes"
    "encoding/json"
    "io"
    "fmt"
    "os"
    "mime/multipart"
    "path/filepath"

    "github.com/knitzsche/restclient/restclient"
)

func help() string {

    text :=
        `Usage: [sudo] snapd-restclient.install COMMAND [VALUE]

Commands:
    ack SNAP-ASSERTION-FILE (requires sudo: note: cp FILE to SNAP_DATA first)
    sideload SNAP-FILE (requires sudo. cp FILE to SNAP_DATA first)
    system-info (returns system info)
    snaps (returns installed snap info)
    login UBUNTU-ACCOUNT-EMAIL PASSWORD [TWO_FACTOR] (requires internet connection. returns auth macaroon)
    test (outputs "yeah")
`
    return text
}

func uri(endpoint string) string {
        //return fmt.Sprintf("https://unix/%s", filepath.Join("v2", endpoint))
        return fmt.Sprintf("http://unix/%s", filepath.Join("v2", endpoint))
}


func main() {

    if len(os.Args) < 2 { 
        fmt.Println("Error: no command arguments provided")
        return
    }   
    args := os.Args[1:]

    switch args[0] {
    case "help":
        fmt.Printf("%s\n", help())
    case "-help":
        fmt.Printf("%s\n", help())
    case "-h":
        fmt.Printf("%s\n", help())
    case "--help":
        fmt.Printf("%s\n", help())
    case "system-info":
        c := restclient.DefaultRestClient()
	u := uri("system-info")
        //fmt.Println(u)
        resp, err := c.SendHTTPRequest(u, "GET", nil)
        if err != nil {
                fmt.Println("error", err)
                return
        }
        fmt.Printf("%v\n", resp)
    case "snaps":
        c := restclient.DefaultRestClient()
	u := uri("snaps")
        //fmt.Println(u)
        resp, err := c.SendHTTPRequest(u, "GET", nil)
        if err != nil {
                fmt.Println("error", err)
                return
        }
        fmt.Printf("%v\n", resp)
    case "ack": //assert must be in SNAP_DATA
        if len(args) < 2 { 
            fmt.Println("Error: no arg for assertion filename.")
            return
        }
	f := filepath.Join(os.Getenv("SNAP_DATA"), args[1])
        _, err := os.Stat(f)
        if os.IsNotExist(err) {
		fmt.Printf("Error: arg %v does not exist.\n", f)
		return
	}
	fmt.Printf("%v exists.\n", f)
	file, err2 := os.Open(f)
	defer file.Close()
	if err2 != nil {
		fmt.Printf("Error:  cannot read %v due to %v.\n", f, err2)
		return
	}
        c := restclient.DefaultRestClient()
	u := uri("assertions")
        resp, err3 := c.SendHTTPRequest(u, "POST", file)
        if err3 != nil {
                fmt.Println("error", err3)
                return
        }
        fmt.Printf("%v\n", resp)
    case "test":
        c := restclient.DefaultRestClient()
	c.Yeah("yeah")
    case "sideload": //assert and snap file must be in SNAP_DATA
	var err error
        if len(args) < 2 { 
            fmt.Println("Error: no arg for sideload snap filename.")
            return
        }
	filePath := filepath.Join(os.Getenv("SNAP_DATA"), args[1])
        _, err = os.Stat(filePath)
        if os.IsNotExist(err) {
		fmt.Printf("Error: arg %v does not exist.\n", args[1])
		return
	}
	file, err1 := os.Open(filePath)
	if err1 != nil {
		fmt.Printf("Error: cannot read %v due to %v.\n", file, err1)
		return
	}
	defer file.Close()

	body := new(bytes.Buffer)
	body_writer := multipart.NewWriter(body)
	part, err2 := body_writer.CreateFormFile("snap", filePath)
	if err2 != nil {
		fmt.Printf("Error: cannot CreateFormFile %v due to %v.\n", file, err2)
		return
	}
	_, err = io.Copy(part, file)
	if err  != nil {
		fmt.Printf("Error: cannot copy snap to part due to %v.\n", file, err)
		return
	}
	u := uri("snaps")
        c := restclient.DefaultRestClient()
	headers := map[string]string{
		"Content-Type": body_writer.FormDataContentType(),
	}
	err = body_writer.Close()
        if err != nil {
                fmt.Println("body_writer close error", err)
                return
        }
        resp, err2 := c.SendHTTPRequestHeaders(u, "POST", body, headers)
        if err2 != nil {
                fmt.Println("error", err2)
                return
        }
        fmt.Printf("%v\n", resp)
    case "login":
	var otp = false // otp means two-factor auth
        if len(args) < 3 { 
            fmt.Println("Error: need two args: Ubuntu account email and password.")
            return
        }
	if len(args) == 4 {
		otp = true
	}
	type Login struct {
		Email string `json:"email"`
		Password string `json:"password"`
		OTP string `json:"otp"`
	}

	l := &Login{}
	l.Email = args[1]
	l.Password = args[2]
	if otp {
		l.OTP = args[3]
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(l)
        c := restclient.DefaultRestClient()
	u := uri("login")
        resp, err := c.SendHTTPRequest(u, "POST", b)
        if err != nil {
                fmt.Println("sendHTTPRequest error", err)
                return
        }
        fmt.Printf("%v\n", resp)

	type Response struct {
		Result     map[string]interface{} `json:"result"`
		Status     string                 `json:"status"`
		StatusCode int                    `json:"status-code"`
		Type       string                 `json:"type"`
	}

	r:= Response{}
	err = json.Unmarshal([]byte(resp), &r)
	if err != nil {
                fmt.Println("Unmarshal error", err)
                return
	}
	fmt.Println(r.Result["macaroon"])
    }
}
