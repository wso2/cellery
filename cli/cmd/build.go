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
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tj/go-spin"
	"github.com/wso2/cellery/cli/util"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var isSpinning = true
var isFirstPrint = true
var white = color.New(color.FgWhite)
var boldWhite = white.Add(color.Bold).SprintFunc()
var tag string
var fileName string

func newBuildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build [OPTIONS]",
		Short: "Build an immutable cell image with required dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			if (len(args) == 0) {
				cmd.Help()
				return nil
			}
			fileName = args[0]
			err := runBuild(tag, fileName)
			if err != nil {
				cmd.Help()
				return err
			}
			return nil
		},
		Example: "  cellery build my-project-v1.0.bal -t myproject",
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

func runBuild(tag string, fileName string) error {
	if fileName == "" {
		return fmt.Errorf("no file name specified")
	}
	var extension = filepath.Ext(fileName)
	var fileNameSuffix = fileName[0 : len(fileName)-len(extension)]

	if tag == "" {
		tag = fileNameSuffix
	}
	go spinner(tag)

	cmd := exec.Command("ballerina", "run", fileName+":lifeCycleBuild")
	execError := ""
	stderrReader, _ := cmd.StderrPipe()
	stderrScanner := bufio.NewScanner(stderrReader)
	go func() {
		for stderrScanner.Scan() {
			execError += stderrScanner.Text()
		}
	}()
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error in executing cell build: %v \n", err)
		errStr := string(stderr.Bytes())
		fmt.Printf("%s\n", errStr)
		os.Exit(1)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("\x1b[31;1m\n Error occurred while building cell image:\x1b[0m %v \n", execError)
		errStr := string(stderr.Bytes())
		fmt.Printf("\x1b[31;1m\n  %s\x1b[0m\n", errStr)
		os.Exit(1)
	}

	outStr := string(stdout.Bytes())
	fmt.Printf("\n\033[36m%s\033[m\n", outStr)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Error in getting current directory location: " + err.Error());
		os.Exit(1)
	}
	folderRenameError := os.Rename(dir+"/target", dir+"/k8s") // rename directory
	if folderRenameError != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	folders := []string{"./k8s"}
	files := []string{fileName}
	output := fileNameSuffix + ".zip"
	err = util.RecursiveZip(files, folders, output);
	os.RemoveAll(dir + "/k8s")
	if err != nil {
		fmt.Printf("\x1b[31;1m Cell build finished with error: \x1b[0m %v \n", err)
		os.Exit(1)
	}

	fmt.Printf("\r\033[32m Successfully built cell image \033[m %s \n", boldWhite(tag+".zip"))
	return nil
}
