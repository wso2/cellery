package main

import (
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

func newImageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "images [OPTIONS]",
		Short: "List cell images",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runImage()
		},
	}
	return cmd
}

func runImage() error {
	data := [][]string{
		{"abc/hr_app_cell", "v1.0.0", "a70ad572a50f", "128MB"},
		{"xyz/Hello_cell", "v1.0.0", "6efa497099d9", "150MB"},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"CELL", "VERSION", "CELL-IMAGE-ID", "SIZE"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
	return nil
}
