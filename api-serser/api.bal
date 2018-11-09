import ballerina/http;
import ballerina/io;
import ballerina/log;
import ballerina/system;
import ballerinax/docker;

@http:ServiceConfig {
    basePath: "/"
}
service<http:Service> cellaryApi bind { port: 8080 } {
    @http:ResourceConfig {
        methods: ["GET"],
        path: "/version"
    }
    celleryVersion (endpoint caller, http:Request req) {
        //create new response
        http:Response res = new;
        res.statusCode = 200;
        res.setContentType("application/json");
        res.setJsonPayload(
            {
                "celleryTool": {
                    "version": "0.50-beta1",
                    "apiVersion": "0.95",
                    "ballerinaVersion": "0.987.0",
                    "gitCommit": "78a6bdb",
                    "built": "Tue Oct 23 22:41:53 2018",
                    "osArch": "darwin/amd64",
                    "experimental": true
                },
                "celleryRepository": {
                    "server": "http://central.cellery.io/api",
                    "apiVersion": "0.95",
                    "authenticated": false
                },
                "kubernetes": {
                    "version": "v1.11.3",
                    "crd": "Not Configured"
                },
                "docker": {
                    "registry": "Not Configured"
                }
            }
        );

        caller->respond(res) but { error e => log:printError("Error sending response", err = e) };
    }

    @http:ResourceConfig {
        methods: ["POST"],
        path: "/cell/init"
    }
    celleryInit (endpoint caller, http:Request req) {
        //create new response
        http:Response res = new;
        res.statusCode = 200;
        res.setContentType("text/plain");
        json jsonValue = check req.getJsonPayload();
        string projectName = check <string> jsonValue.name;
        string responseText = "import wso2/cellery;
// " + projectName + "
cellery:Component " + projectName + " = {
   name: '" + projectName + "',
   source: {
       dockerImage: 'docker.io/" + projectName + ":v1'
   },
   replicas: {
       min: 1,
       max: 1
   },
   container:{},
   env: {},
   apis: {},
   dependencies: {},
   security: {}
};";
       
        res.setTextPayload(untaint responseText);

        caller->respond(res) but { error e => log:printError("Error sending response", err = e) };
    }
}
