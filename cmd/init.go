package main

import (
	"bufio"
	//"bytes"
	//"encoding/json"
	"fmt"
	"github.com/fatih/color"
	i "github.com/oxequa/interact"
	"github.com/spf13/cobra"
	//"io/ioutil"
	//"net/http"
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

	// Define colors
	cyan := color.New(color.FgCyan).SprintFunc()
	white := color.New(color.FgWhite)
	boldWhite := white.Add(color.Bold).SprintFunc()
	whitef := color.New(color.FgWhite)
	faintWhite := whitef.Add(color.Faint).SprintFunc()


	prefix := cyan("?")
	projectName := ""
	projectVersion := ""

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

	// hard coding the cell format for now
	cellTemplate :=
	"import wso2/cellery; \n" +
	"// my-project \n" +
	"cellery:Component my-project = { \n" +
		"   name: 'my-project', \n" +
		"   source: {\n" +
		"      dockerImage: 'docker.io/my-project:v1' \n" +
		"   }, \n" +
		"   replicas: { \n" +
		"      min: 1, \n" +
		"      max: 1 \n" +
		"   }, \n" +
		"   container:{}, \n" +
		"   env: {}, \n" +
		"   apis: {}, \n" +
		"   dependencies: {}, \n" +
		"   security: {} \n" +
	"};"
	//type Payload struct {
	//	Name    string `json:"name"`
	//	Version string `json:"version"`
	//}
	//payload := Payload{}
	//payload.Name = projectName
	//payload.Version = projectVersion
	//
	//bytesPayload, err := json.Marshal(payload)
	//
	//// Call API server to get the version details.
	//resp, err := http.Post(BASE_URL + "/cell/init", "application/json", bytes.NewBuffer(bytesPayload))
	//if err != nil {
	//	fmt.Println("Error connecting to the cellery api server.")
	//}
	//defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)

	// Get current directory
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Error in getting current directory location: " + err.Error());
		os.Exit(1)
	}

	file, err := os.Create(dir + "/" + projectName+ "-" + projectVersion + ".cell")
	if err != nil {
		fmt.Println("Error in creating file: " + err.Error());
		os.Exit(1)
	}
	defer file.Close();
	w := bufio.NewWriter(file)
	//w.WriteString(string(body))
	w.WriteString(fmt.Sprintf("%s", cellTemplate))
	w.Flush()

	fmt.Println("Initialized cell project in dir: " + faintWhite(dir))

	return nil
}
