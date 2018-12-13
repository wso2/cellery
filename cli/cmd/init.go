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
	i "github.com/oxequa/interact"
	"github.com/spf13/cobra"
	"github.com/wso2/cellery/cli/util"
	"os/exec"
	"os"
	"path/filepath"
)


func newInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a cell project",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runInit()
			if err != nil{
				cmd.Help()
				return err
			}
			return nil
		},
	}
	return cmd
}

func runInit() error {
	cyan := color.New(color.FgCyan).SprintFunc()
	white := color.New(color.FgWhite)
	boldWhite := white.Add(color.Bold).SprintFunc()
	whitef := color.New(color.FgWhite)
	faintWhite := whitef.Add(color.Faint).SprintFunc()


	prefix := cyan("?")
	projectName := ""
	projectVersion := ""
	projectPath := ""

	i.Run(&i.Interact{
		Before: func(c i.Context) error{
			c.SetPrfx(color.Output, prefix)
			return nil
		},
		Questions: []*i.Question{
			{
				Before: func(c i.Context) error{
					c.SetPrfx(nil, cyan("?"))
					c.SetDef("my-project", faintWhite("[my-project]"))
					return nil
				},
				Quest: i.Quest{
					Msg:     boldWhite("Enter a project name"),
				},
				Action: func(c i.Context) interface{} {
					projectName, _ = c.Ans().String()
					return nil
				},
			},
			{
				Before: func(c i.Context) error{
					c.SetPrfx(nil, cyan("?"))
					c.SetDef("v1.0.0", faintWhite("[v1.0.0]"))
					return nil
				},
				Quest: i.Quest{
					Msg:     boldWhite("Enter project version"),
				},
				Action: func(c i.Context) interface{} {
					projectVersion, _ = c.Ans().String()
					return nil
				},
			},
		},
	})

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Error in getting current directory location: " + err.Error());
		os.Exit(1)
	}

	if _, err := os.Stat(dir + "/" + projectName); os.IsNotExist(err) {
		projectPath = dir + "/" + projectName
		os.Mkdir(projectPath, os.ModePerm)
	}
	os.Chdir(projectPath)
	cmd := exec.Command("ballerina", "init")
	intitErr := cmd.Start()
	if intitErr != nil {
		fmt.Printf("Error in executing ballerina init: %v \n", err)
		os.Exit(1)
	}
	intitErr = cmd.Wait()
	if err != nil {
		fmt.Printf("\x1b[31;1m  ballerina init finished with error: \x1b[0m %v \n", err)
		os.Exit(1)
	}
	os.Chdir(dir)

	util.CopyFile("cli/resources/sample.bal", projectPath + "/" + projectName+ "-" + projectVersion + ".bal")

	fmt.Println("Initialized cell project in dir: " + faintWhite(dir))

	return nil
}
