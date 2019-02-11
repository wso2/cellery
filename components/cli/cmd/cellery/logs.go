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
	"github.com/spf13/cobra"

	"github.com/cellery-io/sdk/components/cli/pkg/commands"
)

var cellName, component string

func newLogsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs [OPTIONS]",
		Short: "Displays logs for either the cell instance, or a component of a running cell instance.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
				return nil
			}
			cellName := args[0]
			err := commands.RunLogs(cellName, component)
			if err != nil {
				cmd.Help()
				return err
			}
			return nil
		},
		Example: "  cellery logs my-cell\n  cellery logs my-cell -c component_name",
	}
	cmd.Flags().StringVarP(&component, "component", "c", "", "component of a cell")
	return cmd
}
