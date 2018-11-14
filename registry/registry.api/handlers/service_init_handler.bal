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

import ballerina/internal as file;
import ballerina/system;
import ballerina/config;
import ballerina/log;

@final string REGISTRY_ROOT_DIRECTORY_PATH = createRegistryRootDirectory();
@final string ARCHIVES_DIRECTORY_PATH = string `{{REGISTRY_ROOT_DIRECTORY_PATH}}{{FILE_SEPARATOR}}archives`;
@final string REPOSITORY_DIRECTORY_PATH = string `{{REGISTRY_ROOT_DIRECTORY_PATH}}{{FILE_SEPARATOR}}repository`;

# Initializes the registry
# + return - True if server started successfully.
public function initServer() returns (boolean) {
    var initResult = initializeFileStorage();
    match initResult {
        boolean isInitialized => {
            if (isInitialized) {
                log:printInfo("Registry initialized successfully.");
                return isInitialized;
            } else {
                error initErr = {
                    message: "Error occurred while initalizing registry server."
                };
                throw initErr;
            }
        }
        RegistryError regErr => {
            error initErr = {
                message: "Error occurred while initalizing registry server.",
                cause: regErr.cause
            };
            throw initErr;
        }
    }
}

# Initializes the folder structure for artifact storage.
# + return - True if artifacts and archive folders are created.
function initializeFileStorage() returns (boolean|RegistryError) {
    boolean intialized = false;

    try {
        // Creating directory structure.
        file:Path archivesDirectory = new(ARCHIVES_DIRECTORY_PATH);
        file:Path repositoryDirectory = new(REPOSITORY_DIRECTORY_PATH);
        if (!archivesDirectory.exists()) {
            log:printDebug(string `Creating archives directory: {{archivesDirectory.getPathValue()}}`);
            match archivesDirectory.createDirectory() {
                () => {
                    intialized = true;
                }
                error err => {
                    RegistryError regErr = new (string `Error occurred while initializing registry:
                        {{archivesDirectory.getPathValue()}}`, "R2000", cause = err
                    );
                    return regErr;
                }
            }
        } else {
            intialized = true;
            log:printDebug(string `Using existing archive directory: {{archivesDirectory.getPathValue()}}`);
        }

        if (!repositoryDirectory.exists()) {
            log:printDebug(string `Creating repository directory: {{repositoryDirectory.getPathValue()}}`);
            match repositoryDirectory.createDirectory() {
                () => {
                    intialized = intialized && true;
                }
                error err => {
                    RegistryError regErr = new (string `Error occurred while initializing remote repository:
                        {{repositoryDirectory.getPathValue()}}`, "R2000", cause = err
                    );
                    return regErr;
                }
            }
        } else {
            intialized = intialized && true;
            log:printDebug(string `Using existing repository directory: {{repositoryDirectory.getPathValue()}}`);
        }

        return intialized;

    } catch (error err) {
        RegistryError regErr = new ("Error occurred while initializing Cellery Registry.", "r2000");
        return regErr;
    }
}

# Creates or identify existing registry root location.
# + return - Location of the repository.
function createRegistryRootDirectory() returns (string) {
    file:Path rootDir = new (REGISTRY_ROOT_DIRECTORY);
    try {
        if (!rootDir.exists()) {
            log:printDebug(string `Creating Registry root directory: {{rootDir.getPathValue()}}`);
            match rootDir.createDirectory() {
                () => {}
                error err => {
                    error regErrRuntime = {
                        message: string `Unable to create Registry root directory at location: {{rootDir.getPathValue()}}.`,
                        cause: err
                    };
                    throw regErrRuntime;
                }
            }
        }
        return rootDir.getPathValue();
    } catch (error err) {
        error regErrRuntime = {
            message: string `Unable to create or identify Registry root directtory at location: {{rootDir.getPathValue()}}.`,
            cause: err
        };
        throw regErrRuntime;
    }
}
