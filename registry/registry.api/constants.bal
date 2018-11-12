import ballerina/config;

// Configuration for Web API.
@final int API_SERVICE_PORT = config:getAsInt("API_SERVICE_PORT", default = 9090);
@final string API_VERSION = config:getAsString("API_VERSION", default = "0.0.1");

// Multiple CORS can be given with separated commas
@final string CORS_ORIGINS = config:getAsString("CORS_ORIGINS", default = "*");

// Keystore to secure Web API Service
@final string KEYSTORE_FILE = config:getAsString("KEYSTORE_FILE", default = "${ballerina.home}/bre/security/ballerinaKeystore.p12");
@final string KEYSTORE_PASSWORD = config:getAsString("KEYSTORE_PASSWORD", default = "ballerina");


@final string ORG_NAME_REGEX = "^[a-z0-9_]*$";
@final string PACKAGE_NAME_REGEX = "^[a-zA-Z0-9_.]*$";
@final string VERSION_REGEX = "^(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(-(0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(\\.(0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(\\+[0-9a-zA-Z-]+(\\.[0-9a-zA-Z-]+)*)?$";

@final string ARTIFACT_EXTENSION = config:getAsString("ARTIFACT_EXTENSION", default = ".zip");

// Webserver
@final string FILE_SERVER_URL = config:getAsString("FILE_SERVER_URL", default = "https://localhost:8080");
@final string FILE_SERVER_URL_ARTIFACT_LOCATION = config:getAsString("FILE_SERVER_URL_ARTIFACT_LOCATION", default = "/");