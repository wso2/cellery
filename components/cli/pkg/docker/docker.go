/*
 * Copyright (c) 2019 WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package docker

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const docker = "docker"
const push = "push"

type Docker interface {
	ServerVersion() (string, error)
	ClientVersion() (string, error)
	PushImages(dockerImages []string) error
}

type CelleryDockerCli struct {
}

// NewCelleryDockerCli returns a CelleryDockerCli instance.
func NewCelleryDockerCli() *CelleryDockerCli {
	cli := &CelleryDockerCli{}
	return cli
}

// ServerVersion returns the docker server version.
func (cli *CelleryDockerCli) ServerVersion() (string, error) {
	cmd := exec.Command(docker,
		"version",
		"--format",
		"{{.Server.Version}}",
	)
	dockerServerResult, err := cmd.CombinedOutput()
	if err != nil {
		return string(dockerServerResult), fmt.Errorf("error while getting Docker Server version, %v", strings.TrimSpace(string(dockerServerResult)))
	}
	return string(dockerServerResult), nil
}

// ClientVersion returns the docker client version.
func (cli *CelleryDockerCli) ClientVersion() (string, error) {
	cmd := exec.Command(docker,
		"version",
		"--format",
		"{{.Client.Version}}",
	)
	dockerClientResult, err := cmd.CombinedOutput()
	if err != nil {
		return string(dockerClientResult), fmt.Errorf("error while getting Docker Client version, %v", strings.TrimSpace(string(dockerClientResult)))
	}
	return string(dockerClientResult), nil
}

// PushImages pushes docker images.
func (cli *CelleryDockerCli) PushImages(dockerImages []string) error {
	// Todo: Update method signature as PushImage(dockerImage string)
	log.Printf("Pushing docker images [%s]", strings.Join(dockerImages, ", "))
	for _, elem := range dockerImages {
		cmd := exec.Command(docker,
			push,
			elem,
		)
		execError := ""
		stderrReader, _ := cmd.StderrPipe()
		stderrScanner := bufio.NewScanner(stderrReader)
		go func() {
			for stderrScanner.Scan() {
				execError += stderrScanner.Text()
			}
		}()
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Start()
		if err != nil {
			errStr := string(stderr.Bytes())
			fmt.Printf("%s\n", errStr)
			return fmt.Errorf("error occurred while pushing Docker image, %v", err)
		}
		err = cmd.Wait()
		if err != nil {
			fmt.Println()
			fmt.Printf("\x1b[31;1m\nPush Failed.\x1b[0m %v \n", execError)
			fmt.Println("\x1b[31;1m======================\x1b[0m")
			errStr := string(stderr.Bytes())
			fmt.Printf("\x1b[31;1m%s\x1b[0m", errStr)
			return fmt.Errorf("error occurred while pushing cell image, %v", err)
		}
	}
	return nil
}
