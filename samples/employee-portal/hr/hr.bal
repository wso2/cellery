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
cellery:CellStub employeeCell = new("Employee");
cellery:CellStub stocksCell = new("Stock-Options");

public function celleryBuild() {
    // Build HR cell
    io:println("Building HR Cell ...");
    hrCell.addComponent(hr);
    hrCell.addStub(employeeCell);
    hrCell.addStub(stocksCell);

    // Expose API from Cell Gateway & Global Gateway
    hrCell.exposeGlobalAPI(hr);
    hrCell.declareEgress(employeeCell.name, "employeegw_url");
    hrCell.declareEgress(stocksCell.name, "stockgw_url");
    _ = cellery:createImage(hrCell);

}
