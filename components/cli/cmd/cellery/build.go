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

	"github.com/cellery-io/sdk/components/cli/pkg/commands"
	"github.com/cellery-io/sdk/components/cli/pkg/util"
)

// newBuildCommand creates a cobra command which can be invoked to build a cell image from a cell file
func newBuildCommand() *cobra.Command {
	var tag string
	cmd := &cobra.Command{
		Use:   "build <cell-file>",
		Short: "Build an immutable cell image with the required dependencies",
		Args: func(cmd *cobra.Command, args []string) error {
			err := cobra.ExactArgs(1)(cmd, args)
			if err != nil {
				return err
			}
			isProperFile, err := util.FileExists(args[0])
			if err != nil || !isProperFile {
				return fmt.Errorf("expects a proper file as the cell-file, received %s", args[0])
			}
			err = util.ValidateImageTag(tag)
			if err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			commands.RunBuild(tag, args[0])
		},
		Example: "  cellery build employee.bal -t cellery-samples/employee:1.0.0",
	}
	cmd.Flags().StringVarP(&tag, "tag", "t", "",
		"Cell image in the format: <organization>/<cell-image>:<version>")
	_ = cmd.MarkFlagRequired("tag")
	return cmd
}
