import ballerina/internal as file;
import ballerina/io;
import ballerina/log;
import ballerina/system;

# Creates a directory in the given path
# + return - true if successfully created or directory exists
function createDirectory(string path) returns (boolean) {
    file:Path directory = new (path);
    try {
        if (!directory.exists()) {
            match directory.createDirectory() {
                () => {
                    log:printDebug(string `New directory created: {{directory.getPathValue()}}`);
                }
                error err => {
                    error regErrRuntime = {
                        message: string `Unable to create directory at location: {{directory.getPathValue()}}.`,
                        cause: err
                    };
                    throw regErrRuntime;
                }
            }
        }
        return true;
    } catch (error err) {
        error regErrRuntime = {
            message: string `Unable to create directory at location: {{directory.getPathValue()}}.`,
            cause: err
        };
        throw regErrRuntime;
    }
}

# Copy artifact file into file system. Overwrites if file exists.
# + artifactByteChannel - The byte channel of the artifact file.
# + filePath - The path to the file
# + return - Path of the artifact else error cause in creating the artifact file.
function copyArtifactFile (io:ReadableByteChannel artifactByteChannel, string filePath) returns (string | RegistryError) {
    try {
        file:Path artifactFile = new (filePath);
        // Delete if artifact file already exists. This is to overwrite existing artifact content.
        if (artifactFile.exists()) {
            log:printDebug(string `Deleting existing artifact file: {{artifactFile.getPathValue()}}`);
            match artifactFile.delete() {
                () => {}
                error ioError => {
                    RegistryError regErr = new (
                                               string `Error deleting existing artifact file: {{artifactFile.getPathValue()}}.`,
                                               "f2101",
                                               cause = ioError
                    );
                    return regErr;
                }
            }
        }

        // Creating artifact file.
        log:printDebug(string `Creating artifact file: {{artifactFile.getPathValue()}}`);
        var artifactFileCreateRes = artifactFile.createFile();
        match artifactFileCreateRes {
            () => {}
            error err => {
                RegistryError regErr = new (
                                           string `Error creating artifact file: {{artifactFile.getPathValue()}}`,
                                           "f2102",
                                           cause = err
                );
                return regErr;
            }
        }

        log:printDebug(string `Writing content to artifact file: {{artifactFile.getPathValue()}}`);
        var writeToFileRes = writeToFile(artifactByteChannel, artifactFile);
        match writeToFileRes {
            () => {
                log:printDebug(string `Artifact file written: {{artifactFile.getPathValue()}}`);
                return artifactFile.getPathValue();
            }
            error err => {
                RegistryError regErr = new (
                                           string `Error creating artifact file: {{artifactFile.getPathValue()}}`,
                                           "f2103",
                                           cause = err
                );
                return regErr;
            }
        }
    } catch (error err) {
        RegistryError regErr = new (
                                   string `Error creating artifact file`,
                                   "f2000",
                                   cause = err
        );
        return regErr;
    }

}

function writeToFile (io:ReadableByteChannel src, file:Path dstFile) returns (error | ()) {
    io:WritableByteChannel dst = getFileChannel(dstFile);
    // Specifies the number of bytes that should be read from a single read operation.
    int numberOfBytesWritten = 0;
    int readCount = 0;
    byte[] readContent;
    boolean completed = false;
    try {
        // Here is how to specify to read all the content from
        // the source and copy it to the destination.
        while (!completed) {
            match readBytes(src, 1024) {
                (byte[], int) out => {
                    (readContent, readCount) = out;
                    if (readCount <= 0) {
                        //If no content is read we end the loop
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
    }
}

function writeBytes (io:WritableByteChannel byteChannel, byte[] content) returns (int) {
    match byteChannel.write(content, 0) {
        int numberOfBytesWritten => return numberOfBytesWritten;
        error err => throw err;
    }
}

function readBytes (io:ReadableByteChannel byteChannel, int numberOfBytes) returns (byte[], int)|() {
    match byteChannel.read(numberOfBytes) {
        (byte[], int) content => return content;
        error readError => {
            byte[] emptyContent = [];
            return (emptyContent, 0);
        }
    }
}

function getFileChannel (file:Path fileToOpen) returns (io:WritableByteChannel) {
    io:WritableByteChannel byteChannel = io:openWritableFile(untaint fileToOpen.getPathValue());
    return byteChannel;
}
