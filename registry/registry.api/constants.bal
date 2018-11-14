// Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
//
// WSO2 Inc. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

import ballerina/config;
import ballerina/system;

// Configuration for Web API.
@final int API_SERVICE_PORT = config:getAsInt("API_SERVICE_PORT", default = 9090);
@final string API_VERSION = config:getAsString("API_VERSION", default = "0.0.1");

// Multiple CORS can be given with separated commas
@final string CORS_ORIGINS = config:getAsString("CORS_ORIGINS", default = "*");

// Keystore to secure Web API Service
@final string KEYSTORE_FILE = config:getAsString("KEYSTORE_FILE",
    default = "${ballerina.home}/bre/security/ballerinaKeystore.p12");
@final string KEYSTORE_PASSWORD = config:getAsString("KEYSTORE_PASSWORD", default = "ballerina");
@final string TRUSTSTORE_FILE = config:getAsString("TRUSTSTORE_FILE",
    default = "${ballerina.home}/bre/security/ballerinaTruststore.p12");
@final string TRUSTSTORE_PASSWORD = config:getAsString("TRUSTSTORE_PASSWORD", default = "ballerina");

@final string ORG_NAME_REGEX = "^[a-z0-9_]*$";
@final string PACKAGE_NAME_REGEX = "^[a-zA-Z0-9_.]*$";
@final string VERSION_REGEX = "^(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(-(0|[1-9]\\d*|\\d*[a-zA-Z-]" +
    "[0-9a-zA-Z-]*)(\\.(0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(\\+[0-9a-zA-Z-]+(\\.[0-9a-zA-Z-]+)*)?$";

@final string ARTIFACT_EXTENSION = config:getAsString("ARTIFACT_EXTENSION", default = ".zip");

@final string FILE_SEPARATOR = "/";

@final string REGISTRY_ROOT_DIRECTORY = config:getAsString("REGISTRY_ROOT_DIRECTORY",
    default = system:getUserHome() + FILE_SEPARATOR + "cellery-registry-data");

