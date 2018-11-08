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
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetAlignment(3)
	table.SetRowSeparator("-")
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold})
	table.SetColumnColor(
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{tablewriter.FgHiCyanColor},
		tablewriter.Colors{})

	table.AppendBulk(data)
	table.Render()

	return nil
}
