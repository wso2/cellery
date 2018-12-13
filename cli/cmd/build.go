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
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tj/go-spin"
	"github.com/wso2/cellery/cli/util"
	"os"
	"os/exec"
	"strings"
	"time"
)

var isSpinning = true
var isFirstPrint = true
var white = color.New(color.FgWhite)
var boldWhite = white.Add(color.Bold).SprintFunc()
var tag string
var descriptorName string
var balFileName string

func newBuildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build [OPTIONS]",
		Short: "Build an immutable cell image with required dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			if (len(args) == 0) {
				cmd.Help()
				return nil
			}
			descriptorName = args[0]
			err := runBuild(tag, descriptorName)
			if err != nil {
				cmd.Help()
				return err
			}
			return nil
		},
		Example: "  cellery build my-project-1.0.0.cell -t myproject",
	}
	cmd.Flags().StringVarP(&tag, "tag", "t", "", "Name and optionally a tag in the 'name:tag' format")
	return cmd
}

/**
 Spinner
 */
func spinner(tag string) {
	s := spin.New()
	for {
		if (isSpinning) {
			fmt.Printf("\r\033[36m%s\033[m Building %s %s", s.Next(), "image", boldWhite(tag))
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func runBuild(tag string, descriptorName string) error {
	if tag == "" {
		tag = strings.Split(descriptorName, ".")[0]
	}
	if descriptorName == "" {
		return fmt.Errorf("no descriptor name specified")
	}
	go spinner(tag)

	//ballerinaFile := util.FindInDirectory(descriptorName, ".bal")[0]
	//ballerinaFileName := strings.Split(ballerinaFile, ".")[0]
	//
	//dir, errPath := filepath.Abs(filepath.Dir(os.Args[0]))
	//if errPath != nil {
	//	fmt.Println("Error in getting current directory location: " + errPath.Error());
	//	os.Exit(1)
	//}
	//os.Chdir(dir + "/" + descriptorName)

	cmd := exec.Command("ballerina", "run", descriptorName + ":lifeCycleBuild")
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error in executing cell build: %v \n", err)
		os.Exit(1)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("\x1b[31;1m Cell build finished with error: \x1b[0m %v \n", err)
		os.Exit(1)
	}

	folders := []string{"./target"}
	output := strings.Split(descriptorName, ".")[0] + ".zip"
	err = util.RecursiveZip(nil, folders, output);
	if err != nil {
		fmt.Printf("\x1b[31;1m Cell build finished with error: \x1b[0m %v \n", err)
		os.Exit(1)
	}

	fmt.Printf("\r\033[32m Successfully built cell image \033[m %s \n", boldWhite(tag))

	return nil
}
