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
	"github.com/spf13/viper"
	"github.com/wso2/cellery/cli/constants"
	"github.com/wso2/cellery/cli/util"
	"os"
	"path/filepath"
	"strings"
)

func newPullCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull [CELL IMAGE]",
		Short: "pull cell image from the remote repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
				return nil
			}
			cellImage = args[0]
			err := runPull(cellImage)
			if err != nil {
				cmd.Help()
				return err
			}
			return nil
		},
		Example: "  cellery pull mycellery.org/hello:v1",
	}
	return cmd
}

func runPull(cellImage string) error {
	url := constants.REGISTRY_URL + "/" + constants.REGISTRY_ORGANIZATION + "/" + cellImage + "/2.0.0-m1"
	if cellImage == "" {
		return fmt.Errorf("no cell image specified")
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Error in getting current directory location: " + err.Error())
		os.Exit(1)
	}
	response, downloadError := util.DownloadFile(dir+"/"+cellImage, url)

	if downloadError != nil {
		fmt.Printf("\x1b[31;1m Error occurred while pulling the cell image: \x1b[0m %v \n", err)
		os.Exit(1)
	}

	if response.StatusCode == 200 {
		fmt.Printf("\r\033[32m Successfully pulled cell image \033[m %s \n", util.Bold(cellImage))
	}
	if response.StatusCode == 404 {
		fmt.Printf("\x1b[31;1m Error occurred while running cell image:\x1b[0m %v not found in registry\n", cellImage)
	}

	err = util.Unzip(filepath.Join(dir, cellImage), dir)
	if err != nil {
		panic(err)
	}

	viper.SetConfigName("Cellery") // name of config file (without extension)
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")        // optionally look for config in the working directory
	confErr := viper.ReadInConfig() // Find and read the config file

	if confErr != nil { // Handle errors reading the config file
		fmt.Printf("\x1b[31;1m\nError while readng toml file: %s \x1b[0m\n", confErr)
		os.Exit(1)
	}

	organization := viper.GetString("project.organization")
	projectName := strings.Split(cellImage, ".")[0]
	projectVersion := viper.GetString("project.version")

	repoLocation := filepath.Join(util.UserHomeDir(), ".cellery", "repo", organization, projectName, projectVersion)
	repoCreateErr := util.CreateDir(repoLocation)
	if repoCreateErr != nil {
		fmt.Println("Error while creating image location: " + repoCreateErr.Error())
		os.Exit(1)
	}

	zipSrc := filepath.Join(dir, cellImage)
	zipDst := filepath.Join(repoLocation, cellImage)
	zipCopyError := util.CopyFile(zipSrc, zipDst)
	if zipCopyError != nil {
		fmt.Println("Error while saving image: " + zipCopyError.Error())
		os.Exit(1)
	}

	_ = os.Remove(zipSrc)
	return nil
}
