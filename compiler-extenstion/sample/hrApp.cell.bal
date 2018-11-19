import ballerina/config;
import wso2/cellery;

// HR Application
@cellery:App {}
cellery:Component hrApp = {
    name: "hrApp",
    source: {
        dockerImage: "docker.io/hr-app:v1"
    },
    replicas: 1,
    container: [{ port: 9443, protocol: "TCP" }],
    env: {
        "ADMIN_USERNAME": "admin",
        "ADMIN_PASSWORD": "adminpw"
    },
    apis: {
        context: "/hr",
        global: true,
        definitions: [
            {
                path: "/",
                method: "GET"
            }
        ]
    },
    security: {
        ^"type": "JWT",
        issuer: "account.google.com",
        jwksURI: "https://www.googleapis.com/oauth2/v3/certs"
    }
};

// Employee Application
@cellery:App {}
cellery:Component employeeApp = {
    name: "employeeApp",
    source: {
        dockerImage: "docker.io/employee-app:v1"
    },
    replicas: 2,
    container: [{ port: 9445, protocol: "TCP" }],
    env: {
        "MAX_CONNECTION": 10
    },
    apis: {
        context: "/employee",
        global: false,
        definitions: [
            {
                path: "/info",
                method: "GET"
            }
        ]
    },
    security: {
        ^"type": "JWT",
        issuer: "account.google.com",
        jwksURI: "https://www.googleapis.com/oauth2/v3/certs"
    }
};

public function main(string... args) {

}
