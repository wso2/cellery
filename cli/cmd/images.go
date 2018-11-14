package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

type Component struct {
	name, dockerImage string
	ports []int
	deployment, service string
}

func newImageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "images [OPTIONS]",
		Short: "List cell images",
		RunE: func(cmd *cobra.Command, args []string) error {
			if (len(args) == 0) {
				return runImage()
			} else {
				return runImageInformations(args[0])
			}
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
		tablewriter.Colors{tablewriter.FgHiBlueColor},
		tablewriter.Colors{})

	table.AppendBulk(data)
	table.Render()

	return nil
}

func runImageInformations(input string) error {
	componentArrays := [][]Component{
		{
			{"API-GW", "cellery.io/micro-gw:v1.1", []int{443, 80}, "API_GW_deployment", "API_GW_service"},
			{"STS", "cellery.io/sts:v1.0", []int{443}, "STS_deployment", "STS_service"},
			{"hrApp", "hrApp:v1.0", []int{9443}, "hrApp_deployment", "hrApp_service"},
			{"employeeApp", "employee:v.10", []int{8080}, "employeeApp_deployment", "employeeApp_service"},
			{"stockApp", "stockApp:v2.0", []int{8085}, "stockApp_deployment", "stockApp_service"},
		},
		{
			{"API-GW", "cellery.io/micro-gw:v1.1", []int{443, 80}, "API_GW_deployment", "API_GW_service"},
			{"STS", "cellery.io/sts:v1.0", []int{443}, "STS_deployment", "STS_service"},
			{"helloApp", "helloApp:v1.0", []int{9444}, "helloApp_deployment", "helloApp_service"},
		},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"COMPONENT", "DOCKER-IMAGE", "PORTS", "DEPLOYMENT", "SERVICE"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetAlignment(3)
	table.SetRowSeparator("-")
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold})
	table.SetColumnColor(
		tablewriter.Colors{},
		tablewriter.Colors{tablewriter.FgHiBlueColor},
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{})

	if strings.Contains(input, ":") {
		table.AppendBulk(componentArrayToStringArray(getComponentByNameAndVersion(input, componentArrays)))
	} else {
		table.AppendBulk(componentArrayToStringArray(getComponentByUniqueId(input, componentArrays)))
	}
	table.Render()

	return nil
}

func getComponentByNameAndVersion (nameAndVersion string, componentArrays [][]Component) []Component {
	components := map[string][]Component{
		"abc/hr_app_cell:v1.0.0": componentArrays[0],
		"xyz/Hello_cell:v1.0.0": componentArrays[1],
	}

	return components[nameAndVersion]
}

func getComponentByUniqueId (uniqueId string, componentArrays [][]Component) []Component {
	components := map[string][]Component{
		"a70ad572a50f": componentArrays[0],
		"6efa497099d9": componentArrays[1],
	}

	return components[uniqueId]
}

func componentArrayToStringArray(components []Component) [][]string {
	convertedArray := [][]string{}

	for i := 0; i < len(components); i++  {
		component := make([]string, 5)
		component[0] = components[i].name
		component[1] = components[i].dockerImage
		component[2] = intArrayToString(components[i].ports)
		component[3] = components[i].deployment
		component[4] = components[i].service

		convertedArray = append(convertedArray, component)
	}

	return convertedArray
}

func intArrayToString(intArray []int) string {
	return strings.Trim(strings.Replace(fmt.Sprint(intArray), " ", ", ", -1), "[]")
}
