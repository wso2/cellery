package main

import (
	"github.com/spf13/cobra"
	"os"
)

func newCompletionCommand(root *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long: `To load completion run

. <(cellery completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(cellery completion)
`,
		Run: func(cmd *cobra.Command, args []string) {
			root.GenBashCompletion(os.Stdout);
		},
	}
	return cmd
}

