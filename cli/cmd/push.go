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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wso2/cellery/cli/constants"
	"github.com/wso2/cellery/cli/util"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var cellImage string
var cellImageTag string

type Response struct {
	Message string
	Image   ResponseImage
}
type ResponseImage struct {
	Organization  string
	Name          string
	ImageVersion  string
	ImageRevision string
}

func newPushCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push [CELL IMAGE]",
		Short: "push cell image to the remote repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			if (len(args) == 0) {
				cmd.Help()
				return nil
			}
			cellImageTag = args[0]
			err := runPush(cellImageTag)
			if err != nil {
				cmd.Help()
				return err
			}
			return nil
		},
		Example: "  cellery push mycellery.org/hello:v1",
	}
	return cmd
}

func runPush(cellImageTag string) error {
	tags := []string{}
	if cellImageTag == "" {
		return fmt.Errorf("please specify the cell image")
	}

	userVar, err := user.Current()
	if err != nil {
		panic(err)
	}

	repoLocation := filepath.Join(util.UserHomeDir(), ".cellery", "repo")

	zipLocation := ""
	if !strings.Contains(cellImageTag, "/") {
		tags = strings.Split(cellImageTag, ":")
		zipLocation = filepath.Join(repoLocation, userVar.Username, tags[0], tags[1], tags[0]+".zip")
	} else {
		orgName := strings.Split(cellImageTag, "/")[0]
		tags = strings.Split(strings.Split(cellImageTag, "/")[1], ":")
		zipLocation = filepath.Join(repoLocation, orgName, tags[0], tags[1], tags[0]+".zip")
	}
	cellImage = tags[0]+".zip"

	var url = constants.REGISTRY_URL + "/" + constants.REGISTRY_ORGANIZATION + "/" + cellImage + "/2.0.0-m1"
	request, err := util.FileUploadRequest(url, nil, "file", zipLocation, false)
	if err != nil {
		fmt.Printf("\x1b[31;1m Error occurred while pushing the cell image: \x1b[0m %v \n", err)
		os.Exit(1)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	var responseBody string
	if err != nil {
		fmt.Printf("\x1b[31;1m Error occurred while pushing the cell image: \x1b[0m %v \n", err)
		os.Exit(1)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		responseBody = body.String()
		resp.Body.Close()
	}

	if (resp.StatusCode == 200) {
		var response Response
		json.Unmarshal([]byte(responseBody), &response)
		data := map[string]string{
			"ImageRevision": response.Image.ImageRevision,
			"ImageVersion":  response.Image.ImageVersion,
			"Name":          response.Image.Name,
			"Organization":  response.Image.Organization,
		}

		json, err := json.MarshalIndent(data, " ", " ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(json))
		fmt.Printf("\r\033[32mSuccessfully pushed cell image \033[m\n")
	} else {
		fmt.Printf("\x1b[31;1mError occurred while pushing the cell image: \x1b[0m %v \n", err)
		os.Exit(1)
	}
	return nil
}
