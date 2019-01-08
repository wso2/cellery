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
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wso2/cellery/cli/util"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

func newRunCommand() *cobra.Command {
	var cellImage string
	cmd := &cobra.Command{
		Use:   "run [OPTIONS]",
		Short: "Use a cell image to create a running instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
				return nil
			}
			cellImage = args[0]
			err := run(cellImage)
			if err != nil {
				cmd.Help()
				return err
			}
			return nil
		},
		Example: "  cellery run my-project:1.0.0 -n myproject-v1.0.0",
	}
	return cmd
}

func run(cellImageTag string) error {
	if cellImageTag == "" {
		return fmt.Errorf("please specify the cell image")
	}

	fmt.Printf("Running cell image: %s ...\n", util.Bold(cellImageTag))

	userVar, err := user.Current()
	if err != nil {
		panic(err)
	}

	repoLocation := filepath.Join(util.UserHomeDir(), ".cellery", "repo")

	zipLocation := ""
	if !strings.Contains(cellImageTag, "/") {
		tags := strings.Split(cellImageTag, ":")
		zipLocation = filepath.Join(repoLocation, userVar.Username, tags[0], tags[1], tags[0]+".zip")
	} else {
		orgName := strings.Split(cellImageTag, "/")[0]
		tags := strings.Split(strings.Split(cellImageTag, "/")[1], ":")
		zipLocation = filepath.Join(repoLocation, orgName, tags[0], tags[1], tags[0]+".zip")
	}

	if _, err := os.Stat(zipLocation); os.IsNotExist(err) {
		// TODO need to pull from registry if not present, consider for next iteration
		fmt.Printf("\x1b[31;1m\n Error occurred while running cell image:\x1b[0m%v. Image does not exist\n",
			cellImageTag)
		os.Exit(1)
	}

	//create tmp directory
	currentTIme := time.Now()
	timstamp := currentTIme.Format("20060102150405")
	tmpPath := filepath.Join(util.UserHomeDir(), ".cellery", "tmp", timstamp)
	err = util.CreateDir(tmpPath)
	if err != nil {
		panic(err)
	}

	err = util.Unzip(zipLocation, tmpPath)
	if err != nil {
		panic(err)
	}

	kubeYamlDir := filepath.Join(tmpPath, "artifacts", "kubernetes")

	cmd := exec.Command("kubectl", "apply", "-f", kubeYamlDir)
	execError := ""
	stdoutReader, _ := cmd.StdoutPipe()
	stdoutScanner := bufio.NewScanner(stdoutReader)
	go func() {
		for stdoutScanner.Scan() {
			fmt.Printf("\033[36m%s\033[m\n", stdoutScanner.Text())
		}
	}()
	stderrReader, _ := cmd.StderrPipe()
	stderrScanner := bufio.NewScanner(stderrReader)
	go func() {
		for stderrScanner.Scan() {
			execError += stderrScanner.Text()
		}
	}()
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error in executing cell run: %v \n", err)
		os.Exit(1)
	}
	err = cmd.Wait()

	//_ = os.RemoveAll(tmpPath)

	if err != nil {
		fmt.Printf("\x1b[31;1m\n Error occurred while running cell image:\x1b[0m %v \n", execError)
		os.Exit(1)
	}

	fmt.Printf("\nSuccessfully deployed cell image: %s\n", util.Bold(cellImageTag))
	fmt.Println()
	fmt.Println(util.Bold("Whats next ?"))
	fmt.Println("======================")
	fmt.Println("Execute the following command to view the running cells: ")
	fmt.Println("  $ cellery ps ")

	return nil
}
