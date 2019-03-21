/*
 * Copyright (c) 2019 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
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
	"fmt"
	"github.com/cellery-io/sdk/components/cli/pkg/constants"
	"github.com/cellery-io/sdk/components/cli/pkg/util"
	"github.com/manifoldco/promptui"
	"os"
	"os/exec"
	"path/filepath"
)

func RunSetupCreateLocal(isCompleteSelected bool) {
	vmName := ""
	vmLocation := filepath.Join(util.UserHomeDir(), constants.CELLERY_HOME, constants.VM)
	repoCreateErr := util.CreateDir(vmLocation)
	if repoCreateErr != nil {
		fmt.Println("Error while creating VM location: " + repoCreateErr.Error())
		os.Exit(1)
	}
	spinner := util.StartNewSpinner("Downloading Cellery Runtime")
	defer func() {
		spinner.Stop(true)
	}()
	if isCompleteSelected {
		vmName = constants.AWS_S3_ITEM_VM_COMPLETE
	} else {
		vmName = constants.AWS_S3_ITEM_VM_MINIMAL
	}
	util.DownloadFromS3Bucket(constants.AWS_S3_BUCKET, vmName, vmLocation)
	util.ExtractTarGzFile(vmLocation, filepath.Join(util.UserHomeDir(), constants.CELLERY_HOME, constants.VM, vmName))
	util.DownloadFromS3Bucket(constants.AWS_S3_BUCKET, constants.AWS_S3_ITEM_CONFIG, vmLocation)
	util.ReplaceFile(filepath.Join(util.UserHomeDir(), ".kube", "config"), filepath.Join(util.UserHomeDir(), constants.CELLERY_HOME, constants.VM, constants.AWS_S3_ITEM_CONFIG))
	installVM(isCompleteSelected)
}

func createLocal() error {
	var isCompleteSelected = false
	cellTemplate := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U000027A4 {{ .| bold }}",
		Inactive: "  {{ . | faint }}",
		Help:     util.Faint("[Use arrow keys]"),
	}

	cellPrompt := promptui.Select{
		Label:     util.YellowBold("?") + " Select the type of runtime",
		Items:     []string{"Basic", "Complete"},
		Templates: cellTemplate,
	}
	_, value, err := cellPrompt.Run()
	if err != nil {
		return fmt.Errorf("Failed to select an option: %v", err)
	}
	if value == "Basic" {
		isCompleteSelected = true
	}
	RunSetupCreateLocal(isCompleteSelected)

	return nil
}

func installVM(isCompleteSelected bool) error {
	vmName := constants.VM_NAME
	if isCompleteSelected {
		vmName = constants.VM_NAME
	}
	spinner := util.StartNewSpinner("Installing Cellery Runtime")
	defer func() {
		spinner.Stop(true)
	}()
	util.ExecuteCommand(exec.Command(constants.VBOX_MANAGE, "createvm", "--name", vmName, "--ostype", "Ubuntu_64", "--register"), "Error Installing VM")
	util.ExecuteCommand(exec.Command(constants.VBOX_MANAGE, "modifyvm", vmName, "--ostype", "Ubuntu_64", "--cpus", "4", "--memory", "8000", "--natpf1", "guestkube,tcp,,6443,,6443", "--natpf1", "guestssh,tcp,,2222,,22", "--natpf1", "guesthttps,tcp,,443,,443", "--natpf1", "guesthttp,tcp,,80,,80"), "Error Installing VM")
	util.ExecuteCommand(exec.Command(constants.VBOX_MANAGE, "storagectl", vmName, "--name", "hd1", "--add", "sata", "--portcount", "2"), "Error Installing VM")

	vmLocation := filepath.Join(util.UserHomeDir(), constants.CELLERY_HOME, constants.VM, constants.VM_FILE_NAME)
	util.ExecuteCommand(exec.Command(constants.VBOX_MANAGE, "storageattach", vmName, "--storagectl", "hd1", "--port", "1", "--type", "hdd", "--medium", vmLocation), "Error Installing VM")
	util.ExecuteCommand(exec.Command(constants.VBOX_MANAGE, "startvm", vmName, "--type", "headless"), "Error Installing VM")

	return nil
}
