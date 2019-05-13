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

package commands

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/cellery-io/sdk/components/cli/pkg/constants"
	"github.com/cellery-io/sdk/components/cli/pkg/registry/credentials"
	"github.com/cellery-io/sdk/components/cli/pkg/util"
	"github.com/jwt-go"
	"github.com/nokia/docker-registry-client/registry"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var loginConf Conf

type Conf struct {
	CallBackDefaultPort int    `json:"call-back-default-port"`
	CallBackContextPath string `json:"call-back-context-path"`
	CallBackHost        string `json:"call-back-host"`
	CallBackParameter   string `json:"call-back-parameter"`
	RedirectSuccessUrl  string `json:"redirect-success-url"`
	IsHost              string `json:"is-host"`
	IsPort              int    `json:"is-port"`
	SpClientId          string `json:"sp-client-id"`
}

// RunLogin requests the user for credentials and logs into a Cellery Registry
func RunLogin(registryURL string, username string, password string) {
	fmt.Println("Logging into Registry: " + util.Bold(registryURL))

	var registryCredentials = &credentials.RegistryCredentials{
		Registry: registryURL,
		Username: username,
		Password: password,
	}
	isCredentialsProvided := registryCredentials.Username != "" &&
		registryCredentials.Password != ""

	credManager, err := credentials.NewCredManager()
	if err != nil {
		util.ExitWithErrorMessage("Error occurred while creating Credentials Manager", err)
	}
	var isCredentialsAlreadyPresent bool
	if !isCredentialsProvided {
		// Reading the existing credentials
		registryCredentials, err := credManager.GetCredentials(registryURL)
		if registryCredentials == nil {
			registryCredentials = &credentials.RegistryCredentials{
				Registry: registryURL,
			}
		}
		// errors are ignored and considered as credentials not present
		isCredentialsAlreadyPresent = err == nil && registryCredentials.Username != "" &&
			registryCredentials.Password != ""
	}

	if isCredentialsProvided {
		fmt.Println("Logging in with provided Credentials")
	} else if isCredentialsAlreadyPresent {
		fmt.Println("Logging in with existing Credentials")
	} else {
		if password == "" {
			registryCredentials.Username, registryCredentials.Password, err = requestCredentials(username)
			if err != nil {
				util.ExitWithErrorMessage("Error occurred while reading Credentials", err)
			}
		} else {
			registryCredentials.Username = username
			registryCredentials.Password = password
		}
	}

	// Initiating a connection to Cellery Registry (to validate credentials)
	fmt.Println()
	spinner := util.StartNewSpinner("Logging into Cellery Registry " + registryURL)
	_, err = registry.New("https://"+registryURL, registryCredentials.Username, registryCredentials.Password)
	if err != nil {
		spinner.Stop(false)
		if strings.Contains(err.Error(), "401") {
			util.ExitWithErrorMessage("Invalid Credentials", err)
		} else {
			util.ExitWithErrorMessage("Error occurred while initializing connection to the Cellery Registry",
				err)
		}
	}

	if !isCredentialsAlreadyPresent {
		// Saving the credentials
		spinner.SetNewAction("Saving credentials")
		err = credManager.StoreCredentials(registryCredentials)
		if err != nil {
			spinner.Stop(false)
			util.ExitWithErrorMessage("Error occurred while saving Credentials", err)
		}
	}

	spinner.Stop(true)
	util.PrintSuccessMessage(fmt.Sprintf("Successfully logged into Registry: %s", util.Bold(registryURL)))
}

// This function is to request credentials from the user
func requestCredentials(username string) (string, string, error) {
	initConf()
	timeout := make(chan bool)
	ch := make(chan string)
	var code string
	httpPortString := ":" + strconv.Itoa(loginConf.CallBackDefaultPort)
	var codeReceiverPort = loginConf.CallBackDefaultPort
	// This is to start the CLI auth in a different port is the default port is already occupied
	for {
		_, err := net.Dial("tcp", httpPortString)
		if err != nil {
			break
		}
		codeReceiverPort++
		httpPortString = ":" + strconv.Itoa(codeReceiverPort)
	}
	var gitHubAuthUrl = "https://" + loginConf.IsHost + ":" + strconv.Itoa(loginConf.IsPort) +
		"/oauth2/authorize?scope=openid&response_type=" + "code&redirect_uri=http%3A%2F%2F" +
		loginConf.CallBackHost + "%3A" + strconv.Itoa(codeReceiverPort) +
		"%2Fauth&client_id=" + loginConf.SpClientId + "&fidp=google"
	fmt.Println("\n", gitHubAuthUrl)
	go func() {
		mux := http.NewServeMux()
		server := http.Server{Addr: httpPortString, Handler: mux}
		//var timer *time.Timer
		mux.HandleFunc(loginConf.CallBackContextPath, func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				util.ExitWithErrorMessage("Error parsing the code", err)
			}
			code = r.Form.Get(loginConf.CallBackParameter)
			ch <- code
			if err != nil {
				fmt.Printf("Error occurred while establishing the server: %s\n", err)
				return
			}
			if len(code) != 0 {
				http.Redirect(w, r, "https://"+loginConf.IsHost+":"+strconv.Itoa(loginConf.IsPort)+
					"/authenticationendpoint/auth_success.html", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "fail", http.StatusSeeOther)
			}
			flusher, ok := w.(http.Flusher)
			if ok {
				flusher.Flush()
			}
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
		fmt.Printf("Could not resolve the given url %s. Started to operate in the healess "+
			"mode\n", gitHubAuthUrl)
		return handleHeadless(username)
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
		fmt.Println("Time out. Did not receive any code")
		os.Exit(0)
	}
	token := getTokenFromCode(code, codeReceiverPort)
	username, idToken := getSubjectAndUsernameFromJWT(token)
	return username, idToken, nil
}

// This function is to allow this login flow to work in headless mode
func handleHeadless(username string) (string, string, error) {
	var password string
	if username == "" {
		var err error
		fmt.Print("Enter username: ")
		_, err = fmt.Scanln(&username)
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

// Returns the extracted subject from the JWT
func getSubjectAndUsernameFromJWT(token string) (string, string) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(token), &result)
	if err != nil {
		util.ExitWithErrorMessage("Error while unmarshal the id_token", err)
	}
	idToken := result["id_token"]
	jwtToken, _ := jwt.Parse(idToken.(string), nil)
	claims := jwtToken.Claims.(jwt.MapClaims)
	return claims["sub"].(string), idToken.(string)
}

// Returns the JWT from the auth code provided
func getTokenFromCode(code string, port int) string {
	tokenUrl := "https://" + loginConf.IsHost + ":" + strconv.Itoa(loginConf.IsPort) + "/oauth2/token"
	responseBody := "client_id=" + loginConf.SpClientId + "&grant_type=authorization_code&code=" + code +
		"&redirect_uri=http://localhost:" + strconv.Itoa(port) + loginConf.CallBackContextPath
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

// Read the config.json file from cellery home
func initConf() {
	loginConf = Conf{}
	configFilePath := filepath.Join(util.UserHomeDir(), constants.CELLERY_HOME, constants.CONFIG_FILE)
	configFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Printf("Could not read from the file %s \n", configFilePath)
		util.ExitWithErrorMessage("Error reading the file", err)
	}
	err = json.Unmarshal(configFile, &loginConf)
	if err != nil {
		util.ExitWithErrorMessage("Error while unmarshal the yaml", err)
	}
}
