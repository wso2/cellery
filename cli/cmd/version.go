/*
 * Copyright (c) 2018 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"strconv"
	"github.com/fatih/color"
	"github.com/wso2/cellery/cli/constants"
)




func newVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Get cellery runtime version",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runVersion()
			if err != nil{
				cmd.Help()
				return err
			}
			return nil
		},
	}
	return cmd
}

func runVersion() error {

	type Version struct {
		CelleryTool struct {
			Version          string `json:"version"`
			APIVersion       string `json:"apiVersion"`
			BallerinaVersion string `json:"ballerinaVersion"`
			GitCommit        string `json:"gitCommit"`
			Built            string `json:"built"`
			OsArch           string `json:"osArch"`
			Experimental     bool   `json:"experimental"`
		} `json:"celleryTool"`
		CelleryRepository struct {
			Server        string `json:"server"`
			APIVersion    string `json:"apiVersion"`
			Authenticated bool   `json:"authenticated"`
		} `json:"celleryRepository"`
		Kubernetes struct {
			Version string `json:"version"`
			Crd     string `json:"crd"`
		} `json:"kubernetes"`
		Docker struct {
			Registry string `json:"registry"`
		} `json:"docker"`
	}

	// Call API server to get the version details.
	resp, err := http.Get(constants.BASE_API_URL + "/version")
	if err != nil {
		fmt.Println("Error connecting to the cellery api-server.")
	}
	defer resp.Body.Close()

	// Read the response body and get to an bute array.
	body, err := ioutil.ReadAll(resp.Body)
	response := Version{}

	// set the resopose byte array to the version struct. 
    json.Unmarshal(body, &response)

	//Print the version details.
	white := color.New(color.FgWhite)
	boldWhite := white.Add(color.Bold)

	boldWhite.Println("Cellery:")
	fmt.Println(" Version: \t    " + response.CelleryTool.Version);
	fmt.Println(" API version: \t    " + response.CelleryTool.APIVersion);
	fmt.Println(" Ballerina version: " + response.CelleryTool.BallerinaVersion);
	fmt.Println(" Git commit: \t    " + response.CelleryTool.GitCommit);
	fmt.Println(" Built: \t    " + response.CelleryTool.Built);
	fmt.Println(" OS/Arch: \t    " + response.CelleryTool.OsArch);
	fmt.Println(" Experimental: \t    " + strconv.FormatBool(response.CelleryTool.Experimental));

	boldWhite.Println("\nCellery Repository:")
	fmt.Println(" Server: \t    " + response.CelleryRepository.Server);
	fmt.Println(" API version: \t    " + response.CelleryRepository.APIVersion);
	fmt.Println(" Authenticated:     " + strconv.FormatBool(response.CelleryRepository.Authenticated));

	boldWhite.Println("\nKubernetes")
	fmt.Println(" Version: \t    " + response.Kubernetes.Version);
	fmt.Println(" CRD: \t\t    " + response.Kubernetes.Crd);

	boldWhite.Println("\nDocker:")
	fmt.Println(" Registry: \t    " + response.Docker.Registry);

	return nil
}
