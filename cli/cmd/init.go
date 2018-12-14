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

	cellTemplate :=
	 	"import ballerina/io; \n" +
		"import celleryio/cellery; \n" +
		"cellery:Component componentB = { \n" +
		"    name: \"ComponentB\", \n" +
		"    source: { \n" +
		"        image: \"docker.io/wso2vick/component-b:v1\" \n" +
		"    }, \n" +
		"    env: { ENV1: \"\", ENV2: \"\" }, \n" +
		"    ingresses: [ \n" +
		"        { \n" +
		"            name: \"bar\", \n" +
		"            port: \"8080:80\", \n" +
		"            context: \"bar\", \n" +
		"            definitions: [ \n" +
		"                { \n" +
		"                    path: \"*\", \n" +
		"                    method: \"GET,POST,PUT,DELETE\" \n" +
		"                } \n" +
		"            ] \n" +
		"        } \n" +
		"    ] \n" +
		"}; \n" +
		"//Components \n" +
		"cellery:Component componentA = { \n" +
		"   name: \"ComponentA\", \n" +
		"   source: { \n" +
		"       image: \"docker.io/wso2vick/component-a:v1\" \n" +
		"   }, \n" +
		"   env: { ENV1: \"\", ENV2: \"\" }, \n" +
		"   ingresses: [ \n" +
		"       { \n" +
		"           name: \"foo\", \n" +
		"           port: \"8080:80\", \n" +
		"           context: \"foo\", \n" +
		"           definitions: [ \n" +
		"               { \n" +
		"                   path: \"*\", \n" +
		"                   method: \"GET,POST,PUT,DELETE\" \n" +
		"               } \n" +
		"           ] \n" +
		"       } \n" +
		"   ], \n" +
		"   egresses: [ \n" +
		"       { \n" +
		"           parent: componentB.name, \n" +
		"           ingress: componentB.ingresses[0], \n" +
		"           envVar: \"ENV1\", \n" +
		"           resiliency: { \n" +
		"               retryConfig: { \n" +
		"                   interval: 100, \n" +
		"                   count: 10, \n" +
		"                   backOffFactor: 0.5, \n" +
		"                   maxWaitInterval: 20000 \n" +
		"               } \n" +
		"           } \n" +
		"       } \n" +
		"   ] \n" +
		"}; \n" +
		"//Cell \n" +
		"cellery:Cell cellA = new (\"CellA\"); \n" +
		"//Build Function \n" +
		"public function lifeCycleBuild() { \n" +
		"   componentA.env[\"ENV2\"] = \"1VALUE2\"; \n" +
		"   cellA.addComponent(componentA); \n" +
		"   componentB.env[\"ENV1\"] = \"2VALUE1\"; \n" +
		"   componentB.env[\"ENV2\"] = \"2VALUE2\"; \n" +
		"   cellA.addComponent(componentB); \n" +
		"   cellA.egresses = [ \n" +
		"       { \n" +
		"           parent: componentB.name, \n" +
		"           ingress: componentB.ingresses[0], \n" +
		"           envVar: \"ENV1\" \n" +
		"       } \n" +
		"   ]; \n" +
		"   cellA.apis = [ \n" +
		"       { \n" +
		"           parent: componentB.name, \n" +
		"           context:componentB.ingresses[0], \n" +
		"           global: true \n" +
		"       } \n" +
		"   ]; \n" +
		"   _ = cellery:build(cellA); \n" +
		"}";


	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Error in getting current directory location: " + err.Error());
		os.Exit(1)
	}

	if _, err := os.Stat(dir + "/" + projectName + "-" + projectVersion); os.IsNotExist(err) {
		projectPath = dir + "/" + projectName + "-" + projectVersion
		os.Mkdir(projectPath, os.ModePerm)
	}

	file, err := os.Create(projectPath + "/" + projectName+ "-" + projectVersion + ".bal")
	if err != nil {
		fmt.Println("Error in creating file: " + err.Error());
		os.Exit(1)
	}
	defer file.Close();
	w := bufio.NewWriter(file)
	w.WriteString(fmt.Sprintf("%s", cellTemplate))
	w.Flush()
	fmt.Println("Initialized cell project in dir: " + faintWhite(dir))

	return nil
}
