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

package test

type MockDockerCli struct {
	serverVersion string
	clientVersion string
}

// NewMockDockerCli returns a MockDockerCli instance.
func NewMockDockerCli(opts ...func(*MockDockerCli)) *MockDockerCli {
	cli := &MockDockerCli{}
	for _, opt := range opts {
		opt(cli)
	}
	return cli
}

// SetServerVersion sets server version of mock docker cli.
func SetServerVersion(version string) func(*MockDockerCli) {
	return func(cli *MockDockerCli) {
		cli.serverVersion = version
	}
}

// SetClientVersion sets client version of mock docker cli.
func SetClientVersion(version string) func(*MockDockerCli) {
	return func(cli *MockDockerCli) {
		cli.clientVersion = version
	}
}

// ServerVersion returns the docker server version.
func (cli *MockDockerCli) ServerVersion() (string, error) {
	return cli.serverVersion, nil
}

// ClientVersion returns the docker client version.
func (cli *MockDockerCli) ClientVersion() (string, error) {
	return cli.clientVersion, nil
}

// PushImages pushes docker images.
func (cli *MockDockerCli) PushImages(dockerImages []string) error {
	return nil
}
