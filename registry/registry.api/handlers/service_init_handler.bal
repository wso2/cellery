import ballerina/internal as file;
import ballerina/system;
import ballerina/config;
import ballerina/log;

@final string FILE_SEPARATOR = "/";
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
    file:Path rootDir = new (config:getAsString("REGISTRY_REPOSITORY_ROOT_LOCATION",
            default = system:getUserHome() + FILE_SEPARATOR + "cellery-registry"));
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
