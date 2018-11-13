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
# + backend - Name of the docker image
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
    string backend;
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


public type DockerSource record{
    string Dockerfile;
    string tag;
    !...
};

public type ImageSource record{
    string dockerImage;
    !...
};

public type GitSource record{
    string gitRepo;
    string tag;
    !...
};

public type Port record{
    string? backend;
    int port;
    string protocol;
    !...
};

public type Definition record{
    string path;
    string method;
};

public type API record{
    string context;
    Definition[] definitions;
};

public type InternalEgress record{
    string envVar;
    string value;
};

public type ExternalEgress record{
    string host;
    string protocol;
};

public type Dependencies record{
    InternalEgress[] inernalEgresses;
    ExternalEgress[] externalEgresses;
};

public type Security record{
    string ^"type";
    string issuer;
    string jwksURL;
};

public type Component record{
    string backend;
    DockerSource|ImageSource|GitSource source;
    int replicas;
    Port[] container;
    map env;
    API apis;
    Dependencies dependencies;
    Security security;
    !...
};

public type Cell object {
    Component[] components;

    public function addComponent(Component component) {
        components[lengthof components] = component;
    }
};