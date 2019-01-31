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
import ballerina/io;
import ballerina/log;
import ballerina/system;

# Creates a directory in the given path
# + path - The path to the directory
# + return -  RegistryError if the directory could not be created
function createDirectory(string path) returns RegistryError? {
    file:Path directory = new(path);
    if (!directory.exists()) {
        match directory.createDirectory() {
            () => {
                log:printDebug(string `New directory created: {{directory.getPathValue()}}`);
            }
            error err => {
                string msg = string `Unable to create directory at location: {{directory.getPathValue()}}.`;
                RegistryError regErr = new (msg, "R2100", cause = err);
                return regErr;
            }
        }
    }
    return ();
}

# Persist file into file system. Overwrites if file exists.
# + sourceByteChannel - The byte channel of the file source
# + path - The path to persist the file
# + return -  RegistryError if the file could not be peristed
function persistFile(io:ReadableByteChannel sourceByteChannel, string path) returns RegistryError? {
    file:Path file = new(path);
    // Delete if the file already exists.
    if (file.exists()) {
        log:printDebug(string `Deleting existing file: {{file.getPathValue()}}`);
        match file.delete() {
            () => {}
            error ioError => {
                string msg = string `Error while deleting existing file: {{file.getPathValue()}}.`;
                RegistryError regErr = new(msg, "R2100", cause = ioError);
                return regErr;
            }
        }
    }

    // Creating file.
    log:printDebug(string `Creating file: {{file.getPathValue()}}`);
    var artifactFileCreateRes = file.createFile();
    match artifactFileCreateRes {
        () => {}
        error err => {
            string msg = string `Error while creating new file: {{file.getPathValue()}}`;
            RegistryError regErr = new(msg, "R2100", cause = err);
            return regErr;
        }
    }

    log:printDebug(string `Writing content to file: {{file.getPathValue()}}`);
    var writeToFileRes = writeByteChannelToFile(sourceByteChannel, file);
    match writeToFileRes {
        () => {
            log:printDebug(string `File written: {{file.getPathValue()}}`);
        }
        error err => {
            string msg = string `Error while writing to new file: {{file.getPathValue()}}`;
            RegistryError regErr = new(msg, "R2100", cause = err);
            return regErr;
        }
    }

    return ();
}



# Create the file with the given path
# + path - The path of the file to be deleted
# + return -  RegistryError if the file could not be deleted
function createFile(string path) returns RegistryError? {
    file:Path file = new(path);
    log:printDebug(string `Creating file: {{file.getPathValue()}}`);
    var artifactFileCreateRes = file.createFile();
    match artifactFileCreateRes {
        () => {}
        error err => {
            string msg = string `Error while creating new file: {{file.getPathValue()}}`;
            RegistryError regErr = new(msg, "R2100", cause = err);
            return regErr;
        }
    }
    return ();
}

# Delete the file with the given path
# + path - The path of the file to be deleted
# + return -  RegistryError if the file could not be deleted
function deleteFile(string path) returns RegistryError? {
    file:Path file = new(path);
    if (file.exists()) {
        log:printDebug(string `Deleting  file: {{file.getPathValue()}}`);
        match file.delete() {
            () => {}
            error ioError => {
                string msg = string `Error while deleting file: {{file.getPathValue()}}.`;
                RegistryError regErr = new(msg, "R2100", cause = ioError);
                return regErr;
            }
        }
    }
    return ();
}

# Writes a byte channel to the given file
# + src - The source byte channel
# + dstFile - The destination file
# + return -  RegistryError if error occured while writing
function writeByteChannelToFile(io:ReadableByteChannel src, file:Path dstFile) returns error? {
    io:WritableByteChannel dst = io:openWritableFile(untaint dstFile.getPathValue());
    int numberOfBytesWritten = 0;
    int readCount = 0;
    byte[] readContent;
    boolean completed = false;
    try {
        while (!completed) {
            match readBytes(src, 1024) {
                (byte[], int) out => {
                    (readContent, readCount) = out;
                    if (readCount <= 0) {
                        completed = true;
                    }
                    numberOfBytesWritten = writeBytes(dst, readContent);
                }
                () => {
                    error err = { message: "Unexpected error occurred while reading bytes." };
                    return err;
                }
            }
        }
        return ();
    } catch (error err) {
        return err;
    } finally {
        _ = dst.close();
    }
}

# Write the given string to the file with the given path
# + strContent - The string to be written
# + path - The path of the file
# + return -  RegistryError if error occured while writing
function writeStringToFile(string strContent, string path) returns error? {
    file:Path dstFile = new(path);
    io:WritableByteChannel dst = io:openWritableFile(untaint path);
    try {
        byte[] byteContent = strContent.toByteArray("UTF-8");
        _ = writeBytes(dst, byteContent);
        return ();
    } catch (error err) {
        return err;
    } finally {
        _ = dst.close();
    }
}

# Reads the file with the given path and returns the a string content
# + path - The path of the file
# + return -  RegistryError if error occured while reading
function readFile(string path) returns (string|error) {
    io:ReadableByteChannel byteChannel = io:openReadableFile(path);
    io:ReadableCharacterChannel characterChannel = new io:ReadableCharacterChannel(byteChannel, "UTF8");
    match characterChannel.read (100) {
        string characters => {
            _ = byteChannel.close();
            _ = characterChannel.close();
            return characters;
        }
        error err => {
            return err;
        }
    }
}

function writeBytes(io:WritableByteChannel byteChannel, byte[] content) returns (int) {
    match byteChannel.write(content, 0) {
        int numberOfBytesWritten => return numberOfBytesWritten;
        error err => throw err;
    }
}

function readBytes(io:ReadableByteChannel byteChannel, int numberOfBytes) returns (byte[], int)|() {
    match byteChannel.read(numberOfBytes) {
        (byte[], int) content => return content;
        error readError => {
            byte[] emptyContent = [];
            return (emptyContent, 0);
        }
    }
}
