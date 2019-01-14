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
import ballerina/log;

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
    boolean global;
    TCP|HTTP ingress;
    !...
};

public type Egress record{
    string targetComponent?;
    string targetCell?;
    TCP|HTTP ingress?;
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
    map<TCP|HTTP> ingresses;
    Egress[] egresses?;
    boolean isStub = false;
    map<KeyValueInput|SecretInput> parameters?;
    !...
};

public type TCP record{
    int port;
    string host;
    !...
};

public type HTTP record{
    int port;
    string context;
    Definition[] definitions;
    !...
};

public type KeyValueInput record{
    string paramType;
    any value?;
    boolean required;
    !...
};

public type SecretInput record{
    *KeyValueInput;
    string path;
    !...
};

public type CellStub object {
    public string name;
    map<TCP|HTTP> ingresses = {};
    map<string> context = {};

    public function __init(string name) {
        self.name = name;
        self.ingresses = {};
        self.context = {};
    }

    public function addIngress(string context) {
        self.context = { name: context };
    }

    public function getIngress(string ingressName) returns string? {
        return self.context[ingressName];
    }
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
        foreach var (name, ingressTemp) in component.ingresses {
            self.apis[self.apis.length()] = {
                targetComponent: component.name,
                ingress: ingressTemp,
                global: false
            };
        }
    }

    public function exposeAPIFrom(Component component, string ingressName) {
        TCP|HTTP? ingress = component.ingresses[ingressName];
        if (ingress is (TCP|HTTP)) {
            self.apis[self.apis.length()] = {
                targetComponent: component.name,
                ingress: ingress,
                global: false
            };
        } else {
            error err = error("Ingress " + ingressName + " not found in the component " + component.name + ".");
            panic err;
        }
    }

    public function exposeGlobalAPI(Component component) {
        foreach var (name, ingressTemp) in component.ingresses {
            self.apis[self.apis.length()] = {
                targetComponent: component.name,
                ingress: ingressTemp,
                global: true
            };
        }
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

public function getHost(CellImage cellImage, Component component) returns (string) {
    return getValidName(cellImage.name) + "--" + getValidName(component.name) + "-service";
}

public function getContext(TCP|HTTP? httpIngress) returns (string) {
    if (httpIngress is HTTP) {
        return httpIngress.context;
    }
    error err = error("Unable to extract context from ingress");
    panic err;
}

function getValidName(string name) returns string {
    return name.toLower().replace("_", "-").replace(".", "-");
}

public function addParameter(KeyValueInput|SecretInput? param, any value) {
    if (param is (KeyValueInput)) {
        param.value = value;
    } else if (param is SecretInput) {
        param.value = value;
    } else {
        error err = error("Parameter not declared in the component.");
        panic err;
    }
}