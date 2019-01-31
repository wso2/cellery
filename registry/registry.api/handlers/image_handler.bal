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
import ballerina/http;
import ballerina/io;
import ballerina/time;
import cellery/registry;


# Handles the push operation of an image. The image is persisted into the file system and metadata store.
# + sourceFileChannel - The byte channel of the file.
# + image - Cell Image object
# + return - The newly created image, else error.
function handleImagePush(io:ReadableByteChannel sourceFileChannel, CellImage image) returns (CellImage|RegistryError) {

    //Fisrt copy the image to a tmp location and calculate the hash
    time:Time time = time:currentTime();
    int timeValue = time.time;
    string tmpFilename = string `{{image.orgName}}_{{image.name}}_{{image.imageVersion}}_{{timeValue}}`;
    string tmpFilePath = string `{{TMP_DIRECTORY_PATH}}{{FILE_SEPARATOR}}{{tmpFilename}}`;

    string hash;
    match persistFile(sourceFileChannel, tmpFilePath) {
        () => {
            match registry:hash(tmpFilePath, registry:SHA256) {
                string str => {
                    hash = str;
                }
                error err => {
                    string msg = string `Error occured while generating the hash of the file`;
                    log:printError(msg, err = err);
                    RegistryError regErr = new(msg, "R2100", cause = err);
                    return regErr;
                }
            }
        }
        RegistryError err => {
            string msg = string `Error occured while creating temp file`;
            log:printError(string `{{msg}}, {{err.toString()}}`);
            RegistryError regErr = new(msg, "R2100", regError = err);
            return regErr;
        }
    }

    //Create the directory structure according to hash value in archive directory
    string hashPath = string `{{ARCHIVES_DIRECTORY_PATH}}{{FILE_SEPARATOR}}{{hash}}`;
    match createDirectory(hashPath) {
        () => {}
        RegistryError err => {
            string msg = string `Error occured while creating the  archive directory`;
            log:printError(string `{{msg}}, {{err.toString()}}`);
            RegistryError regErr = new(msg, "R2100", regError = err);
            return regErr;
        }
    }

    //Save the file to the above hash directory
    string imagePath = string `{{hashPath}}{{FILE_SEPARATOR}}{{image.name}}{{IMAGE_EXTENSION}}`;
    file:Path imageFile = new(imagePath);
    file:Path tmpImageFile = new(tmpFilePath);
    if (!imageFile.exists()) {
        match tmpImageFile.moveTo(imageFile) {
            () => {
            }
            error err => {
                string msg = string `Error occured while persisting the image file`;
                log:printError(msg, err = err);
                RegistryError regErr = new(msg, "R2100", cause = err);
                //Delete temp file in case of an error
                _ = tmpImageFile.delete();
                return regErr;
            }
        }    
    } else {
        // The image is already existing, hence this could be just a new tagging. Therefore no need to move the tmp file
        _ = tmpImageFile.delete();
    }
    

    //Create the directory structure in repository
    string basePath = REPOSITORY_DIRECTORY_PATH + FILE_SEPARATOR + image.orgName + FILE_SEPARATOR + image.name;
    string versionPath = basePath + FILE_SEPARATOR + TAGS_DIR_NAME + FILE_SEPARATOR + image.imageVersion ;
    string currentPath = versionPath + FILE_SEPARATOR + CURRENT_DIR_NAME;
    string versionRevisionsPath = versionPath + FILE_SEPARATOR + REVISIONS_DIR_NAME + FILE_SEPARATOR + hash;
    string revisionsPath = basePath + FILE_SEPARATOR + REVISIONS_DIR_NAME + FILE_SEPARATOR + hash;

    match createDirectory(currentPath) {
        () => {
        }
        RegistryError err => {
            string msg = "Error occured while creating the the `current` directory "
                + "in repository directory structure";
            log:printError(string `{{msg}}, {{err.toString()}}`);
            RegistryError regErr = new(msg, "R2100", regError = err);
            return regErr;
        }
    }

    match createDirectory(versionRevisionsPath) {
        () => {
        }
        RegistryError err => {
            string msg = "Error occured while creating the the version`revisions` directory "
                + "in repository directory structure";
            log:printError(string `{{msg}}, {{err.toString()}}`);
            RegistryError regErr = new(msg, "R2100", regError = err);
            return regErr;
        }
    }

    match createDirectory(revisionsPath) {
        () => {
        }
        RegistryError err => {
            string msg = "Error occured while creating the the `revisions` directory "
                + "in repository directory structure";
            log:printError(string `{{msg}}, {{err.toString()}}`);
            RegistryError regErr = new(msg, "R2100", regError = err);
            return regErr;
        }
    }

    string currentLinkPath = string `{{currentPath}}{{FILE_SEPARATOR}}link`;

    //Create if the link file for the current revision is not present
    file:Path currentLinkFile = new(currentLinkPath);
    if(!currentLinkFile.exists()) {
        match currentLinkFile.createFile() {
            () => {}
            error err => {
                string msg =
                    string `Error while creating new file for the current link: {{currentLinkFile.getPathValue()}}`;
                log:printError(msg, err = err);
                RegistryError regErr = new(msg, "R2100", cause = err);
                return regErr;
            }
        }
    }

    //Write the hash of the current revision
    match writeStringToFile(hash, currentLinkPath) {
        () => {}
        error err => {
            string msg =
                string `Error while writing content to the current link file: {{currentLinkFile.getPathValue()}}`;
            log:printError(msg, err = err);
            RegistryError regErr = new(msg, "R2100", cause = err);
            return regErr;
        }
    }

    image.imageRevision = hash;

    ////Update the version revisions
    //string versionRevisionLinkPath = string `{{versionRevisionsPath}}{{FILE_SEPARATOR}}link`;
    ////Create if the link file for the revision if not present
    //file:Path versionRevisionLinkFile = new(versionRevisionLinkPath);
    //if(!versionRevisionLinkFile.exists()) {
    //    match versionRevisionLinkFile.createFile() {
    //        () => {}
    //        error err => {
    //            string msg =
    //            string `Error while creating new file for the revision link: {{versionRevisionLinkFile.getPathValue()}}`;
    //            log:printError(msg, err=err);
    //            RegistryError regErr = new(msg, "R2100", cause = err);
    //            return regErr;
    //        }
    //    }
    //}
    //
    ////Write the hash of the current revision
    //match writeStringToFile(hash, versionRevisionLinkPath) {
    //    () => {}
    //    error err => {
    //        string msg =
    //        string `Error while writing content to the revision link file: {{versionRevisionLinkFile.getPathValue()}}`;
    //        log:printError(msg, err=err);
    //        RegistryError regErr = new(msg, "R2100", cause = err);
    //        return regErr;
    //    }
    //}
    //
    ////Update the revisions
    //string revisionLinkPath = string `{{revisionsPath}}{{FILE_SEPARATOR}}link`;
    ////Create if the link file for the revision if not present
    //file:Path revisionLinkFile = new(revisionLinkPath);
    //if(!revisionLinkFile.exists()) {
    //    match revisionLinkFile.createFile() {
    //        () => {}
    //        error err => {
    //            string msg =
    //                string `Error while creating new file for the revision link: {{revisionLinkFile.getPathValue()}}`;
    //            log:printError(msg, err=err);
    //            RegistryError regErr = new(msg, "R2100", cause = err);
    //            return regErr;
    //        }
    //    }
    //}
    //
    ////Write the hash of the current revision
    //match writeStringToFile(hash, revisionLinkPath) {
    //    () => {}
    //    error err => {
    //        string msg =
    //            string `Error while writing content to the revision link file: {{revisionLinkFile.getPathValue()}}`;
    //        log:printError(msg, err=err);
    //        RegistryError regErr = new(msg, "R2100", cause = err);
    //        return regErr;
    //    }
    //}

    return image;
}


# Handles the pull operation of an image. The rtifact is inserted into the file store and metadata store.
# + image - Cell Image object
# + return - The path the image, else error.
function handleImagePull(CellImage image, boolean aRevision) returns (http:Response) {
    http:Response res = new;
    if(aRevision) {
        string revisionPath = REPOSITORY_DIRECTORY_PATH + FILE_SEPARATOR + image.orgName + FILE_SEPARATOR + image.name
            + FILE_SEPARATOR + REVISIONS_DIR_NAME + FILE_SEPARATOR + image.imageRevision;
        log:printDebug(string `Revision  path : {{revisionPath}}`);

        match verifyPath(revisionPath) {
            () => {
                string imageFileName = image.name + IMAGE_EXTENSION;
                string imageFilePath = ARCHIVES_DIRECTORY_PATH + FILE_SEPARATOR + image.imageRevision + FILE_SEPARATOR
                    + imageFileName;
                log:printDebug(string `Requested image file path : {{imageFilePath}}`);
                file:Path imageFile = new (imageFilePath);
                if (!imageFile.exists()){
                    json errorJson = { message: "Image not found." };
                    res.statusCode = 404;
                    res.setJsonPayload(errorJson);
                    log:printError("Image does not exist even though the revision direcotory is present.");
                } else {
                    res.setFileAsPayload(untaint imageFilePath, contentType = mime:APPLICATION_OCTET_STREAM);
                }
            }
            http:Response response => {
                return response;
            }
        }
    } else {
        string repoBasePath = REPOSITORY_DIRECTORY_PATH + FILE_SEPARATOR + image.orgName + FILE_SEPARATOR + image.name
            + FILE_SEPARATOR + TAGS_DIR_NAME + FILE_SEPARATOR + image.imageVersion;
        string currentLinkPath = repoBasePath + FILE_SEPARATOR + CURRENT_DIR_NAME +  FILE_SEPARATOR + LINK_FILE_NAME;
        log:printDebug(string `Current Link file path : {{currentLinkPath}}`);
        match verifyPath(currentLinkPath) {
            () => {
                match readFile(currentLinkPath) {
                    string hash => {
                        string imageFileName = image.name + IMAGE_EXTENSION;
                        string imageFilePath = ARCHIVES_DIRECTORY_PATH + FILE_SEPARATOR + hash + FILE_SEPARATOR
                            + imageFileName;
                        log:printDebug(string `Requested image file path : {{imageFilePath}}`);
                        file:Path imageFile = new (imageFilePath);
                        if (!imageFile.exists()){
                            json errorJson = { message: "Image not found." };
                            res.statusCode = 404;
                            res.setJsonPayload(errorJson);
                            log:printError("Image does not exist even though the link is present.");
                        } else {
                            res.setFileAsPayload(untaint imageFilePath, contentType = mime:APPLICATION_OCTET_STREAM);
                        }
                    }
                    error err => {
                        json errorJson = { message: string `Image {{image.getImageTag()}} pull failed.
                                        Unexpected error occurred.` };
                        res.statusCode = 500;
                        res.setJsonPayload(untaint errorJson);
                        log:printError(err.message);
                    }
                }
            }
            http:Response response => {
                return response;
            }
        }
    }
    return res;
}

# Validates the paths
# + path - The path to be validated
# + return - Error HTTP Response if the path is not valid
function verifyPath(string path)returns (http:Response?) {
    http:Response res = new;
    json responseJson;
    file:Path file = new(path);
    if (path.contains("..")) {
        json errorJson = { message: "Invalid image requested." };
        res.statusCode = 400;
        res.setJsonPayload(errorJson);
        log:printWarn("Invalid path requested. possible tainted information");
        return res;
    } else if (!file.exists()){
        json errorJson = { message: "Image not found." };
        res.statusCode = 404;
        res.setJsonPayload(errorJson);
        log:printWarn("Image not found. Possible tainted information");
        return res;
    }
    return ();
}
