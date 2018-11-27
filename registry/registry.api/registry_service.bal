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
    push(endpoint outboundEP, http:Request request, string orgName, string cellImageName, string cellImageVersion) {
        http:Response response = new;
        json responseJson;
        match verifyInputParams(orgName, cellImageName, imageVersion = cellImageVersion) {
            () => {
                orgName = untaint orgName;
                cellImageName = untaint cellImageName;
                cellImageVersion = untaint cellImageVersion;

                CellImage image = new CellImage();
                image.orgName = orgName;
                image.name = cellImageName;
                image.imageVersion = cellImageVersion;
                image.imageRevision = "";

                string tag = image.getImageTag();
                log:printInfo(string `Persisting image : {{tag}}`);

                match request.getBodyParts() {
                    mime:Entity[] entities => {
                        match getArtifactContent(entities) {
                            (io:ReadableByteChannel) sourceChannel => {
                                match handleImagePush(sourceChannel, image) {
                                    CellImage img => {
                                        responseJson = {
                                            message: string `Image {{tag}} push successful`,
                                            image: {
                                                organization: image.orgName,
                                                name: image.name,
                                                imageVersion: image.imageVersion,
                                                imageRevision: image.imageRevision
                                            }
                                        };
                                        response.setJsonPayload(responseJson);
                                    }
                                    RegistryError err => {
                                        json errorJson = { message: string `Image {{tag}} push failed.
                                        Unexpected error occurred.` };
                                        response.statusCode = 500;
                                        response.setJsonPayload(errorJson);
                                        log:printError(err.message);
                                    }
                                }
                                _ = sourceChannel.close();
                            }
                            error mimeErr => {
                                json mimeErrJson = { message: "Invalid request received." };
                                response.setJsonPayload(mimeErrJson);
                                response.statusCode = 400;
                                log:printError(mimeErr.message);
                            }
                        }
                    }
                    error mimeErr => {
                        json mimeErrJson = { message: "Invalid request received." };
                        response.setJsonPayload(mimeErrJson);
                        response.statusCode = 400;
                        log:printError(mimeErr.message);
                    }
                }
            }
            http:Response res => {
                response = res;
            }
        }
        //Send the response back
        _ = outboundEP->respond(response);

    }

    @http:ResourceConfig {
        methods: ["GET"],
        path: "/images/{orgName}/{cellImageName}/{cellImageVersion}"
    }
    pullVersion(endpoint outboundEP, http:Request request, string orgName, string cellImageName, string cellImageVersion) {
        http:Response response = new;
        json responseJson;
        match verifyInputParams(orgName, cellImageName, imageVersion = cellImageVersion) {
            () => {
                CellImage image = new CellImage();
                image.orgName = untaint orgName;
                image.name = untaint cellImageName;
                image.imageVersion = untaint cellImageVersion;
                image.imageRevision = "";

                string tag = untaint image.getImageTag();
                log:printInfo(string `Image : {{tag}} is being pulled.`);

                response = handleImagePull(image, false);
            }
            http:Response res => {
                response = res;
            }
        }
        //Send the response back
        _ = outboundEP->respond(response);
    }

    @http:ResourceConfig {
        methods: ["GET"],
        path: "/images/r/{orgName}/{cellImageName}/{cellImageRevision}"
    }
    pullRevision(endpoint outboundEP, http:Request request, string orgName, string cellImageName,
                 string cellImageRevision) {
        http:Response response = new;
        json responseJson;

        match verifyInputParams(orgName, cellImageName, imageRevision = cellImageRevision) {
            () => {
                CellImage image = new CellImage();
                image.orgName = untaint orgName;
                image.name = untaint cellImageName;
                image.imageVersion = "";
                image.imageRevision = untaint cellImageRevision;

                string tag = untaint image.getImageTag();
                log:printInfo(string `Revision {{cellImageRevision}} of Image : {{tag}} is being pulled.`);

                response = handleImagePull(image, true);
            }
            http:Response res => {
                response = res;
            }
        }
        //Send the response back
        _ = outboundEP->respond(response);
    }
}

# Get the artifact information from the entities in the multipart request.
# + entities - The entities of the multipart request.
# + return - Map with artifact metadata and byte channel of the artifact. Else error.
function getArtifactContent(mime:Entity[] entities) returns (io:ReadableByteChannel|error) {
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

# Validates the parametes of image Push and Pull calls
# + orgName - Organization name
# + imageName - Cell Image Name
# + imageVersion - Cell Image version
# + imageRevision - Cell Image Revision
# + return - Error HTTP Response if any of the parameters are not valid
function verifyInputParams(string orgName, string imageName, string imageVersion = "", string imageRevision = "")
                                                                                    returns (http:Response?) {
    http:Response res = new;
    json responseJson;
    if (!(check orgName.matches(ORG_NAME_REGEX)) || lengthof orgName > 255) {
        string msg = "Invalid organization name: '" + orgName +"' provided.Only lowercase alphanumerics and " +
            "underscores are allowed in an organization name and the maximum length is 256 characters.";
        json errorJson = {
            message: msg
        };
        res.statusCode = 400;
        res.setJsonPayload(untaint errorJson);
        log:printWarn(errorJson.message.toString() + " Possible tainted information.");
        return res;
    } else if (!(check imageName.matches(IMAGE_NAME_REGEX)) || lengthof imageName > 255) {
        string msg = "Invalid Cell Image name: '" + imageName +"' provided. only alphanumerics, underscores," +
            " dash and periods are allowed in a image name and the maximum length is 256 characters.";
        json errorJson = {
            message: msg
        };
        res.statusCode = 400;
        res.setJsonPayload(untaint errorJson);
        log:printWarn(errorJson.message.toString() + " Possible tainted information.");
        return res;
    }

    //Image version and revision
    if (!imageVersion.equalsIgnoreCase("")) {

    }

    //Image version
    if (!imageVersion.equalsIgnoreCase("")) {
        if (!(check imageVersion.matches(VERSION_REGEX)) || lengthof imageVersion > 50) {
            string msg = "Invalid version provided '" + imageVersion +"'. version should be semver and less than" +
                " 50 characters.";
            json errorJson = {
                message: msg
            };
            res.statusCode = 400;
            res.setJsonPayload(untaint errorJson);
            log:printWarn(errorJson.message.toString() + " Possible tainted information.");
            return res;
        }
    }

    //Image revision
    if (!imageRevision.equalsIgnoreCase("")) {
        if (lengthof imageRevision != 64 || !(check imageRevision.matches(REVISION_REGEX))) {
            string msg = "Invalid Revision provided '" + imageRevision +"'. Revision should be a SHA 256 hash and exactly" +
                " 64 characters.";
            json errorJson = {
                message: msg
            };
            res.statusCode = 400;
            res.setJsonPayload(untaint errorJson);
            log:printWarn(errorJson.message.toString() + " Possible tainted information.");
            return res;
        }
    }
    return ();
}
