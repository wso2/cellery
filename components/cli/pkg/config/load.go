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

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/cellery-io/sdk/components/cli/pkg/constants"
	"github.com/cellery-io/sdk/components/cli/pkg/util"
)

const configFile = "config.json"
const callBackDefaultPort = 8888
const callBackContextPath = "/auth"
const callBackHost = "localhost"
const callBackParameter = "code"
const redirectSuccessUrl = "https://cloud.google.com/sdk/auth_success"
const idpHost = "localhost"
const idpPort = 9443
const spClientId = "htl0MoApITB1j0a7HkqDjc_1REIa"

type Conf struct {
	AuthConf AuthConf `json:"auth"`
}

type AuthConf struct {
	CallBackDefaultPort int    `json:"callBackDefaultPort"`
	CallBackContextPath string `json:"callBackContextPath"`
	CallBackHost        string `json:"callBackHost"`
	CallBackParameter   string `json:"callBackParameter"`
	RedirectSuccessUrl  string `json:"redirectSuccessUrl"`
	IdpHost             string `json:"idpHost"`
	IdpPort             int    `json:"idpPort"`
	SpClientId          string `json:"spClientId"`
}

// LoadConfig reads the config file from the Cellery home and returns the Config struct
func LoadConfig() *Conf {
	var conf = Conf{}
	configFilePath := filepath.Join(util.UserHomeDir(), constants.CELLERY_HOME, configFile)
	configFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Printf("Could not read from the file %s. Reading from default values\n", configFilePath)
		loadFromConstants(&conf)
		return &conf
	}
	err = json.Unmarshal(configFile, &conf)
	if err != nil {
		util.ExitWithErrorMessage("Error while unmarshal the json", err)
	}
	return &conf
}

func loadFromConstants(conf *Conf) {
	conf.AuthConf.CallBackDefaultPort = callBackDefaultPort
	conf.AuthConf.CallBackContextPath = callBackContextPath
	conf.AuthConf.CallBackHost = callBackHost
	conf.AuthConf.CallBackParameter = callBackParameter
	conf.AuthConf.RedirectSuccessUrl = redirectSuccessUrl
	conf.AuthConf.IdpHost = idpHost
	conf.AuthConf.IdpPort = idpPort
	conf.AuthConf.SpClientId = spClientId
}
