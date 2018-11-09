package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tj/go-spin"
	"os/exec"
	"bufio"
	"os"
	"strings"
	"time"
)

func newBuildCommand() *cobra.Command {
	var tag string
	var descriptorName string
	cmd := &cobra.Command{
		Use:   "build [OPTIONS]",
		Short: "Build an immutable cell image with required dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			if (len(args) == 0) {
				cmd.Help()
				return nil
			}
			descriptorName = args[0]
			err := runBuild(tag, descriptorName)
			if err != nil{
				cmd.Help()
				return err
			}
			return nil
		},
		Example: "  cellery build my-project-1.0.0.cell -t myproject",
	}
	cmd.Flags().StringVarP(&tag, "tag", "t", "", "Name and optionally a tag in the 'name:tag' format")
	return cmd
}

func runBuild(tag string, descriptorName string) error {
	if tag == "" {
		tag = strings.Split(descriptorName, ".")[0]
	}
	if descriptorName == "" {
		return fmt.Errorf("no descriptor name specified")
	}

	s := spin.New()
	for i := 0; i < 40; i++ {
		fmt.Printf("\r\033[36m%s\033[m Building %s %q", s.Next(), "image", tag)
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Printf("\n")

	cmd := exec.Command("ballerina", "build", descriptorName)
	stdoutReader, _ := cmd.StdoutPipe()
	stdoutScanner := bufio.NewScanner(stdoutReader)
	go func() {
		for stdoutScanner.Scan() {
			fmt.Println(stdoutScanner.Text())
		}
	}()
	stderrReader, _ := cmd.StderrPipe()
	stderrScanner := bufio.NewScanner(stderrReader)
	go func() {
		for stderrScanner.Scan() {
			fmt.Println(stderrScanner.Text())
		}
	}()
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error in executing cell build: %v \n", err)
		os.Exit(1)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("\x1b[31;1m Cell build finished with error: \x1b[0m %v \n", err)
		os.Exit(1)
	}

	fmt.Printf("\r\033[32m Successfully built cell image \033[m %q \n", tag)

	return nil
}
