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
	"github.com/spf13/viper"
	"github.com/tj/go-spin"
	"github.com/wso2/cellery/cli/util"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"
)

var isSpinning = true
var isFirstPrint = true
var tag string
var fileName string

var cyan = color.New(color.FgCyan)
var cyanBold = cyan.Add(color.Bold).SprintFunc()
var bold = color.New(color.Bold).SprintFunc()
var faint = color.New(color.Faint).SprintFunc()
var green = color.New(color.FgGreen)
var greenBold = green.Add(color.Bold).SprintFunc()

func newBuildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build [OPTIONS]",
		Short: "Build an immutable cell image with required dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
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
		Example: "  cellery build my-project.bal -t myproject:1.0.0",
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
		if isSpinning {
			fmt.Printf("\r\033[36m%s\033[m Building %s %s", s.Next(), "image", bold(tag))
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

	viper.SetConfigName("Cellery") // name of config file (without extension)
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")        // optionally look for config in the working directory
	confErr := viper.ReadInConfig() // Find and read the config file

	if confErr != nil { // Handle errors reading the config file
		fmt.Printf("\x1b[31;1m\nError while readng toml file: %s \x1b[0m\n", confErr)
		os.Exit(1)
	}

	organization := viper.GetString("project.organization")
	projectVersion := viper.GetString("project.version")

	if tag == "" {
		tag = fileNameSuffix + ":" + projectVersion
	}
	go spinner(tag)

	//first clean target directory if exists
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Error in getting current directory location: " + err.Error())
		os.Exit(1)
	}
	_ = os.RemoveAll(filepath.Join(dir, "target"))

	cmd := exec.Command("ballerina", "run", fileName+":celleryBuild")
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
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error in executing cell build: %v \n", err)
		errStr := string(stderr.Bytes())
		fmt.Printf("%s\n", errStr)
		os.Exit(1)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Println()
		fmt.Printf("\x1b[31;1m\nBuild Failed.\x1b[0m %v \n", execError)
		fmt.Println("\x1b[31;1m======================\x1b[0m")
		errStr := string(stderr.Bytes())
		fmt.Printf("\x1b[31;1m%s\x1b[0m", errStr)
		os.Exit(1)
	}

	outStr := string(stdout.Bytes())
	fmt.Printf("\n\033[36m%s\033[m\n", outStr)

	folderCopyError := util.CopyDir(filepath.Join(dir, "target"), filepath.Join(dir, "artifacts"))
	if folderCopyError != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	folders := []string{"artifacts"}
	files := []string{fileName}
	output := fileNameSuffix + ".zip"
	err = util.RecursiveZip(files, folders, output)
	if err != nil {
		fmt.Printf("\x1b[31;1mCell build finished with error: \x1b[0m %v \n", err)
		os.Exit(1)
	}

	_ = os.RemoveAll(filepath.Join(dir, "artifacts"))

	repoLocation := filepath.Join(util.UserHomeDir(), ".cellery", "repo", organization, fileNameSuffix, projectVersion)
	repoCreateErr := util.CreateDir(repoLocation)
	if repoCreateErr != nil {
		fmt.Println("Error while creating image location: " + repoCreateErr.Error())
		os.Exit(1)
	}

	zipSrc := filepath.Join(dir, output)
	zipDst := filepath.Join(repoLocation, output)
	zipCopyError := util.CopyFile(zipSrc, zipDst)
	if zipCopyError != nil {
		fmt.Println("Error while saving image: " + zipCopyError.Error())
		os.Exit(1)
	}

	_ = os.Remove(zipSrc)

	userVar, err := user.Current()
	if err != nil {
		panic(err)
	}

	if organization != userVar.Username {
		tag = organization + "/" + tag
	}

	fmt.Printf("Successfully built cell image: %s\n", bold(tag))
	fmt.Println()
	fmt.Println(bold("Whats next ?"))
	fmt.Println("======================")
	fmt.Println("Execute the following command to run the project: ")
	fmt.Println("  $ cellery run " + tag)
	return nil
}
