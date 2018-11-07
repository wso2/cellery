/*
 * Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 */

package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"istio.io/istio/pkg/log"
	"os"
)

func newCliCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "cellery [OPTIONS] COMMAND [ARG...]",
		Short:         "Manage immutable cell based applications",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version: fmt.Sprintf("%s, build %s", "0.1.0", "c69f31c"),
	}

	cmd.AddCommand(
		newConfigureCommand(),
		newCompletionCommand(cmd),
		newBuildCommand(),
		newImageCommand(),
		newVersionCommand(),
		newInitCommand(),
	)
	return cmd
}

func main() {

	cmd := newCliCommand()
	if err := cmd.Execute(); err != nil {
		log.Error(fmt.Sprintf("%s: %s", "cellery", err))
		os.Exit(1)
	}
}
