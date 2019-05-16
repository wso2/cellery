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

const CONFIG_FILE = "config.json"

type Conf struct {
	AuthConf AuthConf `json:"auth"`
}

type AuthConf struct {
	CallBackDefaultPort int    `json:"callBackDefaultPort"`
	CallBackContextPath string `json:"callBackContextPath"`
	CallBackHost        string `json:"callBackHost"`
	CallBackParameter   string `json:"callBackParameter"`
	RedirectSuccessUrl  string `json:"redirectSuccessUrl"`
	IsHost              string `json:"idpHost"`
	IsPort              int    `json:"idpPort"`
	SpClientId          string `json:"spClientId"`
}

// LoadConfig reads the config file from the Cellery home and returns the Config struct
func LoadConfig() Conf {
	var conf = Conf{}
	//auth = AuthConf{}
	configFilePath := filepath.Join(util.UserHomeDir(), constants.CELLERY_HOME, CONFIG_FILE)
	configFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Printf("Could not read from the file %s \n", configFilePath)
		util.ExitWithErrorMessage("Error reading the file", err)
	}
	err = json.Unmarshal(configFile, &conf)
	if err != nil {
		util.ExitWithErrorMessage("Error while unmarshal the json", err)
	}
	return conf
}
