package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
	"bufio"
	"os"
)

func newBuildCommand() *cobra.Command {
	var t string
	var descriptorName string
	cmd := &cobra.Command{
		Use:   "build [OPTIONS]",
		Short: "Build an immutable cell image with required dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runBuild(t, descriptorName)
			if err != nil{
				cmd.Help()
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&t, "tag", "t", "", "Name and optionally a tag in the 'name:tag' format")
	cmd.Flags().StringVarP(&descriptorName, "descriptor", "d", "", "Cell descriptor name")
	return cmd
}

func runBuild(tag string, descriptorName string) error {
	if tag == "" {
		return fmt.Errorf("no build tag specified")
	}
	if descriptorName == "" {
		return fmt.Errorf("no descriptor name specified")
	}

	fmt.Printf("Building %q\n", tag)
	//s := spin.New()
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
	//for i := 0; i < 40; i++ {
	//	fmt.Printf("\r\033[36m%s\033[m Building %s %q", s.Next(), "image", tag)
	//	time.Sleep(100 * time.Millisecond)
	//}
	fmt.Printf("\r\033[32m Successfully built cell image \033[m %q \n", tag)
	//for i := 0; i < 80; i++ {
	//	fmt.Printf("\r\033[36m%s\033[m Building %q", s.Next(), "STS")
	//	time.Sleep(100 * time.Millisecond)
	//}
	//fmt.Printf("\r\033[32mBuild success %q\033[m\n", "STS")
	//for i := 0; i < 100; i++ {
	//	fmt.Printf("\r\033[36m%s\033[m Building %q", s.Next(), "employee-service")
	//	time.Sleep(100 * time.Millisecond)
	//}
	//fmt.Printf("\r\033[32mBuild success %q\033[m\n", "employee-service")
	//for i := 0; i < 70; i++ {
	//	fmt.Printf("\r\033[36m%s\033[m Building %q", s.Next(), "salary-service")
	//	time.Sleep(100 * time.Millisecond)
	//}
	//fmt.Printf("\r\033[32mBuild success %q\033[m\n", "salary-service")
	//fmt.Printf("Generated cell image %q a70ad572a50f\n",tag)
	return nil
}
