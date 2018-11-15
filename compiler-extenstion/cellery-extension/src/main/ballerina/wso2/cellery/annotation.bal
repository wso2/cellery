// Copyright (c) 2018 WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
//
// WSO2 Inc. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

# Cellery annotation configuration.
#
# + name - Name of the docker image
# + registry - Docker registry url
# + tag - Docker image tag
# + username - Docker registry username
# + password - Docker registry password
# + baseImage - Base image for Dockerfile
# + push - Enable pushing docker image to registry
# + buildImage - Enable docker image build
# + enableDebug - Enable ballerina debug
# + debugPort - Ballerina debug port
# + dockerHost - Docker host IP and docker PORT. ( e.g minikube IP and docker PORT)
# + dockerCertPath - Docker certificate path
public type CelleryConfiguration record {
    string name;
    string registry;
    string tag;
    string username;
    string password;
    string baseImage;
    boolean push;
    boolean buildImage;
    boolean enableDebug;
    int debugPort;
    string dockerHost;
    string dockerCertPath;
};

# @cellery:App annotation to configure cellery artifact generation.
public annotation<type> App CelleryConfiguration;