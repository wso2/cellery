/*
 * Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package credentials

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/cellery-io/sdk/components/cli/pkg/config"
	"github.com/cellery-io/sdk/components/cli/pkg/util"
)

// FromBrowser requests the credentials from the user
func FromBrowser(username string) (string, string, error) {
	auth := config.LoadConfig().AuthConf
	timeout := make(chan bool)
	ch := make(chan string)
	var code string
	httpPortString := ":" + strconv.Itoa(auth.CallBackDefaultPort)
	var codeReceiverPort = auth.CallBackDefaultPort
	// This is to start the CLI auth in a different port is the default port is already occupied
	for {
		_, err := net.Dial("tcp", httpPortString)
		if err != nil {
			break
		}
		codeReceiverPort++
		httpPortString = ":" + strconv.Itoa(codeReceiverPort)
	}
	redirectUrl := url.QueryEscape("http://" + auth.CallBackHost + ":" + strconv.Itoa(codeReceiverPort) + "/auth")
	var gitHubAuthUrl = "https://" + auth.IsHost + ":" + strconv.Itoa(auth.IsPort) +
		"/oauth2/authorize?scope=openid&response_type=" + "code&redirect_uri=" + redirectUrl + "&client_id=" +
		auth.SpClientId + "&fidp=google"

	fmt.Println("\n", gitHubAuthUrl)
	go func() {
		mux := http.NewServeMux()
		server := http.Server{Addr: httpPortString, Handler: mux}
		//var timer *time.Timer
		mux.HandleFunc(auth.CallBackContextPath, func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				util.ExitWithErrorMessage("Error parsing the code", err)
			}
			code = r.Form.Get(auth.CallBackParameter)
			ch <- code
			if len(code) != 0 {
				http.Redirect(w, r, "https://"+auth.IsHost+":"+strconv.Itoa(auth.IsPort)+
					"/authenticationendpoint/auth_success.html", http.StatusSeeOther)
			} else {
				util.ExitWithErrorMessage("Did not receive any code", err)
			}
			flusher, ok := w.(http.Flusher)
			if !ok {
				util.ExitWithErrorMessage("Error in casting the flusher", err)
			}
			flusher.Flush()
			err = server.Shutdown(context.Background())
			if err != nil {
				util.ExitWithErrorMessage("Error while shutting down the server\n", err)
			}
		})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			util.ExitWithErrorMessage("Error while establishing the service", err)
		}
	}()

	err := util.OpenBrowser(gitHubAuthUrl)
	if err != nil {
		fmt.Printf("Could not resolve the given url %s. Started to operate in the headless "+
			"mode\n", gitHubAuthUrl)
		return FromTerminal(username)
	}
	// Setting up a timeout
	go func() {
		time.Sleep(5 * time.Minute)
		timeout <- true
	}()
	// Wait for a code, or timeout
	select {
	case <-ch:
	case <-timeout:
		util.ExitWithErrorMessage("Failed to authenticate", errors.
			New("time out. Did not receive any code"))
	}
	token := getTokenFromCode(code, codeReceiverPort, auth)
	username, accessToken := getUsernameAndTokenFromJWT(token)
	return username, accessToken, nil
}

// FromTerminal is to allow this login flow to work in headless mode
func FromTerminal(username string) (string, string, error) {
	var password string
	if username == "" {
		fmt.Print("Enter username: ")
		_, err := fmt.Scanln(&username)
		if err != nil {
			util.ExitWithErrorMessage("Error reading the input username", err)
		}
	}
	fmt.Print("Enter password/token: ")
	bytePassword, err := terminal.ReadPassword(0)
	if err != nil {
		util.ExitWithErrorMessage("Error reading the input token", err)
	}
	password = strings.TrimSpace(string(bytePassword))
	fmt.Println()
	return username, password, nil
}

// getUsernameAndTokenFromJWT returns the extracted subject from the JWT
func getUsernameAndTokenFromJWT(token string) (string, string) {
	fmt.Println(token)
	var result map[string]interface{}
	err := json.Unmarshal([]byte(token), &result)
	if err != nil {
		util.ExitWithErrorMessage("Error while unmarshal the id_token", err)
	}
	idToken, ok := (result["id_token"]).(string)
	accessToken, ok := (result["access_token"]).(string)
	if !ok {
		util.ExitWithErrorMessage("Error while retrieving the access token", err)
	}
	jwtToken, _ := jwt.Parse(idToken, nil)
	claims := jwtToken.Claims.(jwt.MapClaims)
	sub, ok := claims["sub"].(string)
	if !ok {
		util.ExitWithErrorMessage("Error in casting the subject", err)
	}
	return sub, accessToken
}

// getTokenFromCode returns the JWT from the auth code provided
func getTokenFromCode(code string, port int, auth config.AuthConf) string {
	tokenUrl := "https://" + auth.IsHost + ":" + strconv.Itoa(auth.IsPort) + "/oauth2/token"
	responseBody := "client_id=" + auth.SpClientId + "&grant_type=authorization_code&code=" + code +
		"&redirect_uri=http://localhost:" + strconv.Itoa(port) + auth.CallBackContextPath
	body := strings.NewReader(responseBody)
	// Token request
	req, err := http.NewRequest("POST", tokenUrl, body)
	if err != nil {
		util.ExitWithErrorMessage("Error while creating the code receiving request", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Could not connect to the client at %s", tokenUrl)
		util.ExitWithErrorMessage("Error occurred while connecting to the client", err)
	}
	defer res.Body.Close()
	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		util.ExitWithErrorMessage("Error occurred while reading the response body", err)
	}
	return string(respBody)
}
