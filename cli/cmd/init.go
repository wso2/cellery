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
	"github.com/fatih/color"
	i "github.com/oxequa/interact"
	"github.com/spf13/cobra"
	"os"
	"os/user"
	"path/filepath"
)

func newInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a cell project",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runInit()
			if err != nil {
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

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	prefix := cyan("?")
	orgName := ""
	projectName := ""
	projectVersion := ""
	projectPath := ""

	i.Run(&i.Interact{
		Before: func(c i.Context) error {
			c.SetPrfx(color.Output, prefix)
			return nil
		},
		Questions: []*i.Question{
			{
				Before: func(c i.Context) error {
					c.SetPrfx(nil, cyan("?"))
					c.SetDef(user.Username, faintWhite("["+user.Username+"]"))
					return nil
				},
				Quest: i.Quest{
					Msg: boldWhite("Organization name: "),
				},
				Action: func(c i.Context) interface{} {
					orgName, _ = c.Ans().String()
					return nil
				},
			},
			{
				Before: func(c i.Context) error {
					c.SetPrfx(nil, cyan("?"))
					c.SetDef("my-project", faintWhite("[my-project]"))
					return nil
				},
				Quest: i.Quest{
					Msg: boldWhite("Project name: "),
				},
				Action: func(c i.Context) interface{} {
					projectName, _ = c.Ans().String()
					return nil
				},
			},
			{
				Before: func(c i.Context) error {
					c.SetPrfx(nil, cyan("?"))
					c.SetDef("v1.0.0", faintWhite("[v1.0.0]"))
					return nil
				},
				Quest: i.Quest{
					Msg: boldWhite("Project version: "),
				},
				Action: func(c i.Context) interface{} {
					projectVersion, _ = c.Ans().String()
					return nil
				},
			},
		},
	})

	cellTemplate :=
		"import ballerina/io;\n" +
			"import celleryio/cellery;\n" +
			"\n" +
			"cellery:Component helloWorldComp = {\n" +
			"    name: \"HelloWorld\",\n" +
			"    source: {\n" +
			"        image: \"sumedhassk/hello-world:1.0.0\" \n" +
			"    },\n" +
			"    ingresses: [\n" +
			"        {\n" +
			"            name: \"hello\",\n" +
			"            port: \"9090:80\",\n" +
			"            context: \"hello\",\n" +
			"            definitions: [\n" +
			"                {\n" +
			"                    path: \"/\",\n" +
			"                    method: \"GET\"\n" +
			"                }\n" +
			"            ]\n" +
			"        }\n" +
			"    ]\n" +
			"};\n" +
			"\n" +
			"cellery:Cell helloCell = new(\"HelloCell\");\n" +
			"\n" +
			"public function lifeCycleBuild() {\n" +
			"    helloCell.addComponent(helloWorldComp);\n" +
			"\n" +
			"    helloCell.apis = [\n" +
			"        {\n" +
			"            parent:helloWorldComp.name,\n" +
			"            context: helloWorldComp.ingresses[0],\n" +
			"            global: true\n" +
			"        }\n" +
			"    ];\n" +
			"\n" +
			"    var out = cellery:build(helloCell);\n" +
			"    if(out is boolean) {\n" +
			"        io:println(\"Hello Cell Built successfully.\");\n" +
			"    }\n" +
			"}\n"

	tomlTemplate := "[project]\n" +
		"org-name = \"" + orgName + "\"\n" +
		"version = \"" + projectVersion + "\"\n"

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Error in getting current directory location: " + err.Error())
		os.Exit(1)
	}

	if _, err := os.Stat(dir + "/" + projectName); os.IsNotExist(err) {
		projectPath = dir + "/" + projectName
		os.Mkdir(projectPath, os.ModePerm)
	}

	balFile, err := os.Create(projectPath + "/" + projectName + ".bal")
	if err != nil {
		fmt.Println("Error in creating Ballerina File: " + err.Error())
		os.Exit(1)
	}
	defer balFile.Close()
	balW := bufio.NewWriter(balFile)
	balW.WriteString(fmt.Sprintf("%s", cellTemplate))
	balW.Flush()

	tomlFile, err := os.Create(projectPath + "/" + projectName + ".toml")
	if err != nil {
		fmt.Println("Error in creating Toml File: " + err.Error())
		os.Exit(1)
	}
	defer tomlFile.Close()
	tomlW := bufio.NewWriter(tomlFile)
	tomlW.WriteString(fmt.Sprintf("%s", tomlTemplate))
	tomlW.Flush()

	fmt.Println("Initialized cell project in dir: " + faintWhite(projectPath))

	return nil
}
