package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tj/go-spin"
	"time"
	
)

func newBuildCommand() *cobra.Command {
	var t string
	cmd := &cobra.Command{
		Use:   "build [OPTIONS]",
		Short: "Build an immutable cell image with required dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runBuild(t)
			if err != nil{
				cmd.Help()
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&t, "tag", "t", "", "Name and optionally a tag in the 'name:tag' format")
	return cmd
}

func runBuild(tag string) error {
	if tag == "" {
		return fmt.Errorf("no build tag specified")
	}

	fmt.Printf("Generating %q\n", tag)
	s := spin.New()
	for i := 0; i < 40; i++ {
		fmt.Printf("\r\033[36m%s\033[m Building %q", s.Next(), "Gateway")
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Printf("\r\033[32mBuild success %q\033[m\n", "Gateway")
	for i := 0; i < 80; i++ {
		fmt.Printf("\r\033[36m%s\033[m Building %q", s.Next(), "STS")
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Printf("\r\033[32mBuild success %q\033[m\n", "STS")
	for i := 0; i < 100; i++ {
		fmt.Printf("\r\033[36m%s\033[m Building %q", s.Next(), "employee-service")
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Printf("\r\033[32mBuild success %q\033[m\n", "employee-service")
	for i := 0; i < 70; i++ {
		fmt.Printf("\r\033[36m%s\033[m Building %q", s.Next(), "salary-service")
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Printf("\r\033[32mBuild success %q\033[m\n", "salary-service")
	fmt.Printf("Generated cell image %q a70ad572a50f\n",tag)
	return nil
}
