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

import ballerina/http;
import ballerina/io;
import ballerina/log;
import ballerina/internal as file;
import ballerina/crypto;
import ballerina/mime;

@final boolean initialized = initServer();

endpoint http:Listener apiEndpoint {
    port: API_SERVICE_PORT,
    secureSocket: {
        keyStore: {
            path: KEYSTORE_FILE,
            password: KEYSTORE_PASSWORD
        }
    }
};

@http:ServiceConfig {
    basePath: "/registry/" + API_VERSION,
    cors: {
        allowOrigins: CORS_ORIGINS.split(","),
        allowCredentials: true
    }
}

service<http:Service> registry bind apiEndpoint {

    @http:ResourceConfig {
        methods: ["POST"],
        path: "/images/{orgName}/{cellImageName}/{cellImageVersion}"
    }
    push (endpoint outboundEP, http:Request request, string orgName, string cellImageName, string cellImageVersion) {
        http:Response res = new;
        json responseJson;
        if (!(check orgName.matches(ORG_NAME_REGEX)) || lengthof orgName > 255) {
            json errorJson = {
                message: string `Invalid organization name provided '{{orgName}}'. Only lowercase alphanumerics and
                underscores are allowed in an organization name and the maximum length is 256 characters.`
            };
            res.statusCode = 400;
            res.setJsonPayload(untaint errorJson);
            log:printError(errorJson.message.toString() + " Possible tainted information.");
        } else if (!(check cellImageName.matches(PACKAGE_NAME_REGEX)) || lengthof cellImageName > 255) {
            json errorJson = {
                message: string `Invalid Cell Image name provided '{{cellImageName}}'. only alphanumerics,
                underscores and periods are allowed in a module name and the maximum length is 256 characters.`
            };
            res.statusCode = 400;
            res.setJsonPayload(untaint errorJson);
            log:printError(errorJson.message.toString() + " Possible tainted information.");
        } else if (!(check cellImageVersion.matches(VERSION_REGEX)) || lengthof cellImageVersion > 50) {
            json errorJson = {
                message: string `Invalid module version provided '{{cellImageVersion}}'. version should be
                semver and less than 50 characters.`
            };
            res.statusCode = 400;
            res.setJsonPayload(untaint errorJson);
            log:printError(errorJson.message.toString() + " Possible tainted information.");
        } else {
            orgName = untaint orgName;
            cellImageName = untaint cellImageName;
            cellImageVersion = untaint cellImageVersion;
            log:printInfo(string `Persisting image : {{orgName}}/{{cellImageName}}:{{cellImageVersion}}`);

            try {
                string orgPath = string `{{REPOSITORY_DIRECTORY_PATH}}{{FILE_SEPARATOR}}{{orgName}}`;
                string cellImagePath = string `{{orgPath}}{{FILE_SEPARATOR}}{{cellImageName}}`;
                string cellVersionPath = string `{{cellImagePath}}{{FILE_SEPARATOR}}{{cellImageVersion}}`;

                _ = createDirectory(orgPath);
                _ = createDirectory(cellImagePath);
                _ = createDirectory(cellVersionPath);

                string dstFilePath = string `{{cellVersionPath}}{{FILE_SEPARATOR}}{{cellImageName}}.zip`;

                match request.getBodyParts() {
                    mime:Entity[] entities => {
                        match getArtifactContent(entities) {
                            (io:ReadableByteChannel) sourceChannel => {
                                var archiveFileCopyResult = copyArtifactFile(sourceChannel, dstFilePath);
                                match archiveFileCopyResult {
                                    string => {
                                        responseJson = {
                                            message: string `Image {{orgName}}/{{cellImageName}}:{{cellImageVersion}} push successful`
                                        };
                                        res.setJsonPayload(responseJson);
                                    }
                                    RegistryError err => {
                                        json errorJson = { message: string `Image {{orgName}}/{{cellImageName}}:{{cellImageVersion}} push failed.
                                        Unexpected error occurred.` };
                                        res.statusCode = 500;
                                        res.setJsonPayload(errorJson);
                                        log:printError(err.message);
                                    }
                                }

                            }
                            error mimeErr => {
                                json mimeErrJson = { message: "Invalid request received." };
                                res.setJsonPayload(mimeErrJson);
                                res.statusCode = 400;
                                log:printError(mimeErr.message);
                            }
                        }
                    }
                    error mimeErr => {
                        json mimeErrJson = { message: "Invalid request received." };
                        res.setJsonPayload(mimeErrJson);
                        res.statusCode = 400;
                        log:printError(mimeErr.message);
                    }
                }
            } catch (error err) {
                json errorJson = { message: string `Image {{orgName}}/{{cellImageName}}:{{cellImageVersion}} push failed.
                                        Unexpected error occurred.` };
                res.statusCode = 500;
                res.setJsonPayload(errorJson);
                log:printError(err.message, err = err);
            }
        }
        //Send the response back
        _ = outboundEP -> respond(res);
    }


    @http:ResourceConfig {
        methods: ["GET"],
        path: "/images/{orgName}/{cellImageName}/{cellImageVersion}"
    }
    pull (endpoint outboundEP, http:Request request, string orgName, string cellImageName, string cellImageVersion) {
        http:Response res = new;
        json responseJson;
        if (!(check orgName.matches(ORG_NAME_REGEX)) || lengthof orgName > 255) {
            json errorJson = {
                message: string `Invalid organization name provided '{{orgName}}'. Only lowercase alphanumerics and
                underscores are allowed in an organization name and the maximum length is 256 characters.`
            };
            res.statusCode = 400;
            res.setJsonPayload(untaint errorJson);
            log:printError(errorJson.message.toString() + " Possible tainted information.");
        } else if (!(check cellImageName.matches(PACKAGE_NAME_REGEX)) || lengthof cellImageName > 255) {
            json errorJson = {
                message: string `Invalid Cell Image name provided '{{cellImageName}}'. only alphanumerics,
                underscores and periods are allowed in a module name and the maximum length is 256 characters.`
            };
            res.statusCode = 400;
            res.setJsonPayload(untaint errorJson);
            log:printError(errorJson.message.toString() + " Possible tainted information.");
        } else if (!(check cellImageVersion.matches(VERSION_REGEX)) || lengthof cellImageVersion > 50) {
            json errorJson = {
                message: string `Invalid module version provided '{{cellImageVersion}}'. version should be
                semver and less than 50 characters.`
            };
            res.statusCode = 400;
            res.setJsonPayload(untaint errorJson);
            log:printError(errorJson.message.toString() + " Possible tainted information.");
        } else {
            orgName = untaint orgName;
            cellImageName = untaint cellImageName;
            cellImageVersion = untaint cellImageVersion;
            log:printInfo(string `Image : {{orgName}}/{{cellImageName}}:{{cellImageVersion}} is being pulled.`);

            try {
                res = getArtifact(request, orgName, cellImageName, cellImageVersion);
            } catch (error err) {
                json errorJson = { message: string `Image {{orgName}}/{{cellImageName}}:{{cellImageVersion}} pull failed.
                                        Unexpected error occurred.` };
                res.statusCode = 500;
                res.setJsonPayload(errorJson);
                log:printError(err.message, err = err);
            }
        }
        //Send the response back
        _ = outboundEP -> respond(res);
    }
}

# Get the artifact information from the entities in the multipart request.
# + entities - The entities of the multipart request.
# + return - Map with artifact metadata and byte channel of the artifact. Else error.
function getArtifactContent (mime:Entity[] entities) returns (io:ReadableByteChannel | error) {
    io:ReadableByteChannel? artifactByteChannel;
    foreach mimeEntity in entities {
        if (mimeEntity.getContentDisposition() != null) {
            if (mimeEntity.getContentDisposition().name == "file") {
                match mimeEntity.getByteChannel() {
                    error mimeErr => {
                        return mimeErr;
                    }
                    io:ReadableByteChannel byteChannel => {
                        artifactByteChannel = byteChannel;
                    }
                }
            }
        }
    }
    match artifactByteChannel {
        () => {
            error err = {
                message: "Artifact is missing in the payload. artifact is a required parameter."
            };
            return err;
        }
        io:ReadableByteChannel byteChannel => {
            return untaint (byteChannel);
        }
    }
}

# HTTP response to be sent back for pull request
# + orgName - Organization name
# + cellImageName - Cell Image name
# + cellImageVersion - Cell Image Version
# + return - The HTTP response with cell image location
function getArtifact (http:Request req, string orgName, string cellImageName, string cellImageVersion) returns
                                                                                                           (http:Response) {
    http:Response res = new;
    try {
        string FS = FILE_SEPARATOR;
        string artifactPath = REPOSITORY_DIRECTORY_PATH + FILE_SEPARATOR +
                                orgName + FILE_SEPARATOR +
                                cellImageName + FILE_SEPARATOR +
                                cellImageVersion + FILE_SEPARATOR +
                                cellImageName + ARTIFACT_EXTENSION;
        log:printDebug(string `Requested file path : {{artifactPath}}`);
        file:Path artifactFile = new (artifactPath);
        if (artifactPath.contains("..")) {
            json errorJson = { message: "Invalid path requested." };
            res.statusCode = 400;
            res.setJsonPayload(errorJson);
            log:printError("Invalid path requested. possible tainted information");
        } else if (!artifactFile.exists()){
            json errorJson = { message: "Image not found." };
            res.statusCode = 404;
            res.setJsonPayload(errorJson);
            log:printError("Image not found. Possible tainted information");
        } else {
            res.setFileAsPayload(artifactPath, contentType = mime:APPLICATION_OCTET_STREAM);
        }
    } catch (error err) {
        error regErrRuntime = {
            message: string `Unable to create the HTTP response with cell image image :
                {{orgName}}/{{cellImageName}}:{{cellImageVersion}} is being pulled.`,
            cause: err
        };
        log:printError(err.message, err = err);
        throw regErrRuntime;
    }
    return res;
}
