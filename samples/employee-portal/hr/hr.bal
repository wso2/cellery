import ballerina/io;
import celleryio/cellery;


//HR component
cellery:Component hr = {
    name: "hr",
    source: {
        image: "docker.io/wso2vick/sampleapp-hr"
    },
    env: { employeegw_url: "", stockgw_url: "" },
    ingresses: {
        "hr": {
            name: "hr",
            port: "8080:80",
            context: "info",
            definitions: [
                {
                    path: "/",
                    method: "GET"
                }
            ]
        }
    }

};

// Cell Intialization
cellery:CellImage hrCell = new("HR");
cellery:CellImage employeeCell = new("Employee");
cellery:CellImage stocksCell = new("Stock-Options");

public function celleryBuild() {
    _ = cellery:createImage(employeeCell);
    _ = cellery:createImage(stocksCell);

    // Build HR cell
    io:println("Building HR Cell ...");
    hrCell.addComponent(hr);

    // Expose API from Cell Gateway & Global Gateway
    hrCell.exposeGlobalAPI(hr);
    hrCell.declareEgress(employeeCell.name, "employeegw_url");
    hrCell.declareEgress(stocksCell.name, "stockgw_url");
    _ = cellery:createImage(hrCell);

}
