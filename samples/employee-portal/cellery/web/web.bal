import ballerina/io;
import celleryio/cellery;


// Web Component
cellery:Component webComponent = {
    name: "web-ui",
    source: {
        image: "docker.io/celleryio/sampleapp-employee"
    },
    ingresses: {
        employee: new cellery:WebIngress(
                      8080,
                      uri = "/testApp"
        )
    }
};


cellery:CellImage webCell = new();

public function build(string orgName, string imageName, string imageVersion) {
    io:println("Building Web Cell ...");


    // Add components to Cell
    webCell.addComponent(webComponent);
    webCell.exposeGlobal(webComponent, properties = <cellery:WebProperties>{ vhost: "abc.com", context: "/demo", TLS:
        cellery:TLS_REQUIRED });

    _ = cellery:createImage(webCell, orgName, imageName, imageVersion);
}
