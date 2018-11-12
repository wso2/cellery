package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"github.com/manifoldco/promptui"
	i "github.com/oxequa/interact"
	"github.com/spf13/cobra"
	"github.com/tj/go-spin"
	"time"
)

func newConfigureCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configure",
		Short: "configure the cellery installation",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigure()
		},
	}
	return cmd
}

func runConfigure() error {
	// Define colors
	yellow := color.New(color.FgYellow).SprintFunc()
	faint := color.New(color.Faint).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	white := color.New(color.FgWhite)
	boldWhite := white.Add(color.Bold).SprintFunc()
	whitef := color.New(color.FgWhite)
	faintWhite := whitef.Add(color.Faint).SprintFunc()

	cellTemplate := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U000027A4 {{ .| cyan }}",
		Inactive: "  {{ . | white }}",
		Selected: green("\U00002713 ") + boldWhite("cell context: ") + "{{ .  | faint }}",
		Help: faint("[Use arrow keys]"),
	}

	environmentTemplate := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U000027A4 {{ .| cyan }}",
		Inactive: "  {{ . | white }}",
		Selected: green("\U00002713 ") + boldWhite("cell environment: ") + "{{ .  | faint }}",
		Help: faint("[Use arrow keys]"),
	}


	cellPrompt := promptui.Select{
		Label: yellow("?") + " Select cellery cluster endpoint",
		Items: []string{"Local", "GKE", "AWS"},
		Templates: cellTemplate,
	}
	environmentPrompt := promptui.Select{
		Label: yellow("?") + " Select cellery environment",
		Items: []string{"Dev", "Staging", "Test", "Prod"},
		Templates: environmentTemplate,
	}
	_, _, err := cellPrompt.Run()
	if err != nil {
		return fmt.Errorf("failed to select context: %v", err)

	}
	_, _, err = environmentPrompt.Run()

	if err != nil {
		return fmt.Errorf("failed to select context: %v", err)

	}

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
					c.SetDef("", faintWhite("[press enter to use default]"))
					return nil
				},
				Quest: i.Quest{
					Msg:     boldWhite("Enter gateway URL"),
				},
				Action: func(c i.Context) interface{} {
					projectName, _ = c.Ans().String()
					return nil
				},
			},
			{
				Before: func(c i.Context) error{
					c.SetPrfx(nil, cyan("?"))
					c.SetDef("", faintWhite("[press enter to use default]"))
					return nil
				},
				Quest: i.Quest{
					Msg:     boldWhite("Enter STS URL"),
				},
				Action: func(c i.Context) interface{} {
					projectVersion, _ = c.Ans().String()
					return nil
				},
			},
			{
				Before: func(c i.Context) error{
					c.SetPrfx(nil, cyan("?"))
					c.SetDef("cellery.io", faintWhite("[cellery.io]"))
					return nil
				},
				Quest: i.Quest{
					Msg:     boldWhite("Enter docker registry"),
				},
				Action: func(c i.Context) interface{} {
					projectVersion, _ = c.Ans().String()
					return nil
				},
			},
		},
	})

	// Define Spinner
	s := spin.New()

	// Define writers
	writer := uilive.New()
	writer2 := uilive.New()

	// start listening for updates and render
	writer.Start()
	writer2.Start()

	for {
		fmt.Printf("\r\033[36m%s\033[m Configuring ...", s.Next())
		fmt.Fprintf(writer, "Downloading.. (%d/) GB\n", 100)
		fmt.Fprintf(writer, "Downloading 2.. (%d/) GB\n", 100)
		time.Sleep(time.Second)
	}

	fmt.Fprintln(writer, "Finished: Downloaded 100GB")
	writer.Stop() // flush and stop rendering

	fmt.Printf("\rCellery is configured succesfully\n")
	return nil
}
