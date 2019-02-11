/*
 * Copyright (c) 2018 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cellery-io/sdk/components/cli/pkg/constants"
	"github.com/cellery-io/sdk/components/cli/pkg/util"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func RunSetup() error {
	bold := color.New(color.Bold).SprintFunc()
	selectTemplate := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U000027A4 {{ .| bold }}",
		Inactive: "  {{ . | faint }}",
		Selected: bold("Selected cluster: ") + "{{ . }}",
		Help:     util.Faint("[Use arrow keys]"),
	}

	cellPrompt := promptui.Select{
		Label:     util.YellowBold("?") + " Setup Cellery runtime",
		Items:     []string{constants.CELLERY_SETUP_MANAGE, constants.CELLERY_SETUP_CREATE, constants.CELLERY_SETUP_SWITCH},
		Templates: selectTemplate,
	}
	_, value, err := cellPrompt.Run()
	if err != nil {
		return fmt.Errorf("Failed to select an option: %v", err)
	}

	switch value {
	case constants.CELLERY_SETUP_MANAGE:
		{
			manageEnvironment()
		}
	case constants.CELLERY_SETUP_CREATE:
		{
			createEnvironment()
		}
	default:
		{
			selectEnvironment()
		}
	}
	return nil
}

func selectEnvironment() error {
	bold := color.New(color.Bold).SprintFunc()
	cellTemplate := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U000027A4 {{ .| bold }}",
		Inactive: "  {{ . | faint }}",
		Selected: bold("Selected cluster: ") + "{{ . }}",
		Help:     util.Faint("[Use arrow keys]"),
	}

	cellPrompt := promptui.Select{
		Label:     util.YellowBold("?") + " Select a Cellery Installed Kubernetes Cluster",
		Items:     getContexts(),
		Templates: cellTemplate,
	}
	_, value, err := cellPrompt.Run()
	if err != nil {
		return fmt.Errorf("failed to select cluster: %v", err)
	}

	setContext(value)
	fmt.Printf(util.GreenBold("\n\U00002714") + " Successfully configured Cellery.\n")
	fmt.Println()
	fmt.Println(bold("What's next ?"))
	fmt.Println("======================")
	fmt.Println("To create your first project, execute the command: ")
	fmt.Println("  $ cellery init ")
	return nil
}

func getContexts() []string {
	contexts := []string{}
	cmd := exec.Command("kubectl", "config", "view", "-o", "json")
	stdoutReader, _ := cmd.StdoutPipe()
	stdoutScanner := bufio.NewScanner(stdoutReader)
	output := ""
	go func() {
		for stdoutScanner.Scan() {
			output = output + stdoutScanner.Text()
		}
	}()
	stderrReader, _ := cmd.StderrPipe()
	stderrScanner := bufio.NewScanner(stderrReader)

	execError := ""
	go func() {
		for stderrScanner.Scan() {
			execError += stderrScanner.Text()
		}
	}()
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error in executing cellery setup: %v \n", err)
		os.Exit(1)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("\x1b[31;1m Error occurred while configuring cellery: \x1b[0m %v \n", execError)
		os.Exit(1)
	}
	jsonOutput := &Config{}
	errJson := json.Unmarshal([]byte(output), jsonOutput)
	if errJson != nil {
		fmt.Println(errJson)
	}

	for i := 0; i < len(jsonOutput.Contexts); i++ {
		contexts = append(contexts, jsonOutput.Contexts[i].Name)
	}
	return contexts
}

func setContext(context string) error {
	cmd := exec.Command("kubectl", "config", "use-context", context)
	stderrReader, _ := cmd.StderrPipe()
	stderrScanner := bufio.NewScanner(stderrReader)

	execError := ""
	go func() {
		for stderrScanner.Scan() {
			execError += stderrScanner.Text()
		}
	}()
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error in executing cellery setup: %v \n", err)
		os.Exit(1)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("\x1b[31;1m Error occurred while configuring cellery: \x1b[0m %v \n", execError)
		os.Exit(1)
	}
	return nil
}

func createEnvironment() error {
	bold := color.New(color.Bold).SprintFunc()
	cellTemplate := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U000027A4 {{ .| bold }}",
		Inactive: "  {{ . | faint }}",
		Selected: bold("Selected cluster: ") + "{{ . }}",
		Help:     util.Faint("[Use arrow keys]"),
	}

	cellPrompt := promptui.Select{
		Label:     util.YellowBold("?") + " Select an environment to be installed",
		Items:     []string{constants.CELLERY_CREATE_LOCAL, constants.CELLERY_CREATE_KUBEADM, constants.CELLERY_CREATE_GCP, constants.CELLERY_SETUP_BACK},
		Templates: cellTemplate,
	}
	_, value, err := cellPrompt.Run()
	if err != nil {
		return fmt.Errorf("Failed to install environment: %v", err)
	}

	switch value {
	case constants.CELLERY_CREATE_LOCAL:
		{
			vmLocation := filepath.Join(util.UserHomeDir(), ".cellery", "vm")
			repoCreateErr := util.CreateDir(vmLocation)
			if repoCreateErr != nil {
				fmt.Println("Error while creating VM location: " + repoCreateErr.Error())
				os.Exit(1)
			}
			util.DownloadFromS3Bucket(constants.AWS_S3_BUCKET, constants.AWS_S3_ITEM_VM, vmLocation)
			util.ExtractTarGzFile(vmLocation, filepath.Join(util.UserHomeDir(), ".cellery", "vm", constants.AWS_S3_ITEM_VM))
			util.DownloadFromS3Bucket(constants.AWS_S3_BUCKET, constants.AWS_S3_ITEM_CONFIG, vmLocation)
			installVM()
		}
	case constants.CELLERY_CREATE_KUBEADM:
		{
			fmt.Println("Installing kubeadm")
		}
	case constants.CELLERY_CREATE_GCP:
		{
			fmt.Println("Installing GCP")
		}
	default:
		{
			RunSetup()
		}
	}

	fmt.Printf(util.GreenBold("\n\U00002714") + " Successfully installed Cellery runtime.\n")
	fmt.Println()
	fmt.Println(bold("What's next ?"))
	fmt.Println("======================")
	fmt.Println("To create your first project, execute the command: ")
	fmt.Println("  $ cellery init ")
	return nil
}

func installVM() error {
	util.ExecuteCommand(exec.Command("VBoxManage", "createvm", "--name", constants.VM_NAME, "--ostype", "Ubuntu_64", "--register"), "Error Installing VM")
	util.ExecuteCommand(exec.Command("VBoxManage", "modifyvm", constants.VM_NAME, "--ostype", "Ubuntu_64", "--cpus", "4", "--memory", "8000", "--natpf1", "guestkube,tcp,,6443,,6443", "--natpf1", "guestssh,tcp,,2222,,22"), "Error Installing VM")
	util.ExecuteCommand(exec.Command("VBoxManage", "storagectl", constants.VM_NAME, "--name", "hd1", "--add", "sata", "--portcount", "2"), "Error Installing VM")

	vmLocation := filepath.Join(util.UserHomeDir(), ".cellery", "vm", constants.VM_FILE_NAME)
	util.ExecuteCommand(exec.Command("VBoxManage", "storageattach", constants.VM_NAME, "--storagectl", "hd1", "--port", "1", "--type", "hdd", "--medium", vmLocation), "Error Installing VM")
	util.ExecuteCommand(exec.Command("VBoxManage", "startvm", constants.VM_NAME, "--type", "headless"), "Error Installing VM")

	return nil
}

func manageEnvironment() error {
	bold := color.New(color.Bold).SprintFunc()
	cellTemplate := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U000027A4 {{ .| bold }}",
		Inactive: "  {{ . | faint }}",
		Selected: bold("Selected cluster: ") + "{{ . }}",
		Help:     util.Faint("[Use arrow keys]"),
	}

	cellPrompt := promptui.Select{
		Label:     util.YellowBold("?") + " Select an environment to be installed",
		Items:     getManageEnvOptions(),
		Templates: cellTemplate,
	}
	_, value, err := cellPrompt.Run()
	if err != nil {
		return fmt.Errorf("Failed to install environment: %v", err)
	}

	switch value {
	case constants.CELLERY_MANAGE_STOP:
		{
			util.ExecuteCommand(exec.Command("VBoxManage", "controlvm", constants.VM_NAME, "poweroff"), "Error stopping VM")
		}
	case constants.CELLERY_MANAGE_START:
		{
			util.ExecuteCommand(exec.Command("VBoxManage", "startvm", constants.VM_NAME, "--type", "headless"), "Error starting VM")
		}
	case constants.CELLERY_MANAGE_CLEANUP:
		{
			fmt.Println("Cleaning up")
		}
	default:
		{
			RunSetup()
		}
	}
	return nil
}

func isVmRuning() bool {
	cmd := exec.Command("vboxmanage", "showvminfo", constants.VM_NAME)
	stdoutReader, _ := cmd.StdoutPipe()
	stdoutScanner := bufio.NewScanner(stdoutReader)
	output := ""
	go func() {
		for stdoutScanner.Scan() {
			output = output + stdoutScanner.Text()
		}
	}()
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error occurred while checking VM status: %v \n", err)
		os.Exit(1)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("\x1b[31;1m Error occurred while checking VM status \x1b[0m %v \n", err)
		os.Exit(1)
	}

	if strings.Contains(output, "running (since") {
		return true
	}
	return false
}

func getManageEnvOptions() []string {
	if isVmRuning() {
		return []string{constants.CELLERY_MANAGE_STOP, constants.CELLERY_MANAGE_CLEANUP, constants.CELLERY_SETUP_BACK}
	}
	return []string{constants.CELLERY_MANAGE_START, constants.CELLERY_MANAGE_CLEANUP, constants.CELLERY_SETUP_BACK}
}
