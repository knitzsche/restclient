/*
 * Copyright (C) 2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package main

import (
	"fmt"
	"path/filepath"
	"github.com/knitzsche/restclient/restclient"
)

func uri(endpoint string) string {
        //return fmt.Sprintf("https://unix/%s", filepath.Join("v2", endpoint))
        return fmt.Sprintf("http://unix/%s", filepath.Join("v2", endpoint))
}

func main() {

	c := restclient.DefaultRestClient()

	u := uri("assertions/model")	
	fmt.Println(u)
	resp, err := c.SendHTTPRequest(u, "GET", nil)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	fmt.Printf("RESULT: %v\n", resp)

}
