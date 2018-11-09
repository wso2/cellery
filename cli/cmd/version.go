package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"strconv"
	"github.com/fatih/color"
)

const BASE_URL = "http://localhost:8080"


func newVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Get cellery runtime version",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runVersion()
			if err != nil{
				cmd.Help()
				return err
			}
			return nil
		},
	}
	return cmd
}

func runVersion() error {

	type Version struct {
		CelleryTool struct {
			Version          string `json:"version"`
			APIVersion       string `json:"apiVersion"`
			BallerinaVersion string `json:"ballerinaVersion"`
			GitCommit        string `json:"gitCommit"`
			Built            string `json:"built"`
			OsArch           string `json:"osArch"`
			Experimental     bool   `json:"experimental"`
		} `json:"celleryTool"`
		CelleryRepository struct {
			Server        string `json:"server"`
			APIVersion    string `json:"apiVersion"`
			Authenticated bool   `json:"authenticated"`
		} `json:"celleryRepository"`
		Kubernetes struct {
			Version string `json:"version"`
			Crd     string `json:"crd"`
		} `json:"kubernetes"`
		Docker struct {
			Registry string `json:"registry"`
		} `json:"docker"`
	}

	// Call API server to get the version details.
	resp, err := http.Get(BASE_URL + "/version")
	if err != nil {
		fmt.Println("Error connecting to the cellery api-server.")
	}
	defer resp.Body.Close()

	// Read the response body and get to an bute array.
	body, err := ioutil.ReadAll(resp.Body)
	response := Version{}

	// set the resopose byte array to the version struct. 
    json.Unmarshal(body, &response)

	//Print the version details.
	white := color.New(color.FgWhite)
	boldWhite := white.Add(color.Bold)

	boldWhite.Println("Cellery:")
	fmt.Println(" Version: \t    " + response.CelleryTool.Version);
	fmt.Println(" API version: \t    " + response.CelleryTool.APIVersion);
	fmt.Println(" Ballerina version: " + response.CelleryTool.BallerinaVersion);
	fmt.Println(" Git commit: \t    " + response.CelleryTool.GitCommit);
	fmt.Println(" Built: \t    " + response.CelleryTool.Built);
	fmt.Println(" OS/Arch: \t    " + response.CelleryTool.OsArch);
	fmt.Println(" Experimental: \t    " + strconv.FormatBool(response.CelleryTool.Experimental));

	boldWhite.Println("\nCellery Repository:")
	fmt.Println(" Server: \t    " + response.CelleryRepository.Server);
	fmt.Println(" API version: \t    " + response.CelleryRepository.APIVersion);
	fmt.Println(" Authenticated:     " + strconv.FormatBool(response.CelleryRepository.Authenticated));

	boldWhite.Println("\nKubernetes")
	fmt.Println(" Version: \t    " + response.Kubernetes.Version);
	fmt.Println(" CRD: \t\t    " + response.Kubernetes.Crd);

	boldWhite.Println("\nDocker:")
	fmt.Println(" Registry: \t    " + response.Docker.Registry);

	return nil
}
