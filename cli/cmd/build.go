package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tj/go-spin"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var isSpinning = true
var isFirstPrint = true
var white = color.New(color.FgWhite)
var boldWhite = white.Add(color.Bold).SprintFunc()
var tag string
var descriptorName string
var balFileName string

func newBuildCommand() *cobra.Command {
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

/**
 Spinner
 */
func spinner(tag string) {
	s := spin.New()
	for {
		if (isSpinning) {
			fmt.Printf("\r\033[36m%s\033[m Building %s %s", s.Next(), "image", boldWhite(tag))
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func createBalFile() string {
	// Create a temp directory
	tempDir, err := ioutil.TempDir("", "cell")
	if err != nil {
		log.Fatal(err)
	}

	// Copy file
	r, err := os.Open(descriptorName)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	balFileName = tempDir + "/" + descriptorName + ".bal"
	w, err := os.Create(balFileName)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	if err != nil {
		panic(err)
	}

	return tempDir
}

func trim(stream string) string {
	var trimmedString string
	if (strings.Contains(stream, ".cell.balx")) {
		trimmedString = strings.Replace(stream, ".cell.balx", ".celx", -1)
	} else if (strings.Contains(stream, ".bal")) {
		trimmedString = strings.Replace(stream, ".bal", "", -1)
	} else if (strings.Contains(stream, ".cell")) {
		trimmedString = strings.Replace(stream, ".cell", "", -1)
	} else {
		trimmedString = stream
	}
	return trimmedString
}

func runBuild(tag string, descriptorName string) error {
	if tag == "" {
		tag = strings.Split(descriptorName, ".")[0]
	}
	if descriptorName == "" {
		return fmt.Errorf("no descriptor name specified")
	}

	// Start spinner in a seperate thread
	go spinner(tag)

	// Move cell file to a ballerina file
	tempDir := createBalFile()

	// clean up after finishing the task
	defer os.RemoveAll(tempDir)

	// Run ballerina command
	cmd := exec.Command("ballerina", "build", balFileName)
	stdoutReader, _ := cmd.StdoutPipe()
	stdoutScanner := bufio.NewScanner(stdoutReader)
	go func() {
		for stdoutScanner.Scan() {
			isSpinning = false
			if (isFirstPrint) {
				isFirstPrint = false
				// At the first time, print with a new line
				fmt.Println("\n  " + stdoutScanner.Text())
			} else {
				fmt.Println("  " + trim(stdoutScanner.Text()))
			}
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

	// Move balx file to a celx
	os.Rename(descriptorName + ".balx", strings.Replace(descriptorName, ".cell", ".celx", -1))

	fmt.Printf("\r\033[32m Successfully built cell image \033[m %s \n", boldWhite(tag))

	return nil
}
