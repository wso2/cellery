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

public type DockerSource record{
    string Dockerfile;
    string tag;
    !...
};

public type ImageSource record{
    string image;
    !...
};

public type GitSource record{
    string gitRepo;
    string tag;
    !...
};

public type Definition record{
    string path;
    string method;
    !...
};

public type API record{
    string name?;
    string targetComponent;
    string|Ingress? context;
    boolean global;
    Ingress|Definition[] definitions?;
    !...
};

public type Egress record{
    string targetComponent?;
    string targetCell?;
    Ingress ingress?;
    string envVar?;
    Resiliency resiliency?;
    !...
};

public type Resiliency record{
    RetryConfig retryConfig?;
    FailoverConfig failoverConfig?;
    !...
};

public type RetryConfig record {
    int interval;
    int count;
    float backOffFactor;
    int maxWaitInterval;
    !...
};

public type FailoverConfig record {
    int timeOut;
    !...
};

public type Component record{
    string name;
    ImageSource source;
    int replicas = 1;
    map<string> env?;
    map<Ingress> ingresses?;
    Egress[] egresses?;
    boolean isStub = false;
    !...
};

public type CellStub object {
    public string name;
    map<Ingress> ingresses = {};
    public function __init(string name) {
        self.name = name;
    }
};

public type Ingress record{
    string port;
    string context?;
    Definition[] definitions?;
    !...
};

public type CellImage object {
    public string name;
    public Component[] components = [];
    public API?[] apis = [];
    public Egress?[] egresses = [];

    public function addComponent(Component component) {
        self.components[self.components.length()] = component;
    }

    public function exposeAPIsFrom(Component component) {
        foreach var (name, ingress) in component.ingresses {
            self.apis[self.apis.length()] = {
                targetComponent: component.name,
                context: ingress,
                global: false
            };
        }
    }

    public function exposeAPIFrom(Component component, string ingressName) {
        self.apis[self.apis.length()] = {
            targetComponent: component.name,
            context: component.ingresses[ingressName],
            global: false
        };
    }

    public function exposeGlobalAPI(Component component) {
        foreach var (name, ingress) in component.ingresses {
            self.apis[self.apis.length()] = {
                targetComponent: component.name,
                context: ingress,
                global: true
            };
        }
    }

    public function declareEgress(string cellName, string envVarName) {
        self.egresses[self.egresses.length()] = {
            targetCell: cellName,
            envVar: envVarName
        };
    }

    public function addStub(CellStub stub) {
        Component temp = {
            name: stub.name,
            source: {
                image: ""
            },
            ingresses: stub.ingresses,
            isStub: true
        };
        self.addComponent(temp);
    }

    public function __init(string name) {
        self.name = name;
    }
};

# Build the cell aritifacts
#
# + cellImage - The cell image definition
# + return - true/false
public extern function createImage(CellImage cellImage) returns (boolean|error);