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
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/nokia/docker-registry-client/registry"

	"github.com/cellery-io/sdk/components/cli/pkg/util"
)

// RunLogin requests the user for credentials and logs into a Cellery Registry
func RunLogin(registryURL string) {
	fmt.Println("Logging into Registry: " + util.Bold(registryURL))

	// Reading the existing credentials
	config := util.ReadUserConfig()
	existingEncodedCredentials := config.Credentials[registryURL]
	existingCredentials, err := base64.StdEncoding.DecodeString(existingEncodedCredentials)
	isCredentialsAlreadyPresent := err == nil && strings.Contains(string(existingCredentials), ":")

	var username string
	var password string
	if isCredentialsAlreadyPresent {
		existingCredentialsSplit := strings.Split(string(existingCredentials), ":")
		username = existingCredentialsSplit[0]
		password = existingCredentialsSplit[1]

		fmt.Println("Logging in with existing Credentials")
	} else {
		username, password, err = util.RequestCredentials()
		if err != nil {
			util.ExitWithErrorMessage("Error occurred while reading Credentials", err)
		}
	}

	// Initiating a connection to Cellery Registry (to validate credentials)
	fmt.Println()
	spinner := util.StartNewSpinner("Logging into Cellery Registry " + registryURL)
	_, err = registry.New("https://"+registryURL, username, password)
	if err != nil {
		spinner.Stop(false)
		if strings.Contains(err.Error(), "401") {
			util.ExitWithErrorMessage("Invalid Credentials", err)
		} else {
			util.ExitWithErrorMessage("Error occurred while initializing connection to the Cellery Registry", err)
		}
	}

	if !isCredentialsAlreadyPresent {
		// Saving the credentials
		spinner.SetNewAction("Saving credentials")
		encodedCredentials := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
		config.Credentials[registryURL] = encodedCredentials
		err = util.SaveUserConfig(config)
		if err != nil {
			spinner.Stop(false)
			util.ExitWithErrorMessage("Error occurred while saving Credentials", err)
		}
	}

	spinner.Stop(true)
	util.PrintSuccessMessage(fmt.Sprintf("Successfully logged into Registry: %s", util.Bold(registryURL)))
}
