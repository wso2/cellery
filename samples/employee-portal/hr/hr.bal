import ballerina/io;
import celleryio/cellery;


//HR component
cellery:Component hr = {
    name: "hr",
    source: {
        image: "docker.io/wso2vick/sampleapp-hr"
    },
    env: { employeegw_url: "", stockgw_url: "" },
    ingresses: [
        {
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
    ]
};

// Cell Intialization
cellery:Cell hrCell = new("HR");
cellery:CellStub employeeCell = new("Employee");
cellery:Cell stocksCell = new("Stock-Options");

public function lifeCycleBuild() {
    _ = cellery:build(employeeCell);
    _ = cellery:build(stocksCell);
    // Build HR cell
    io:println("Building HR Cell ...");
    hrCell.addComponent(hr);
    hrCell.apis = [
    {
            parent: hr.name,
            context: hr.ingresses[0],
            global: true
    }
    ];
    hrCell.egresses = [
    {
            targetCell: "Employee",
            ingress: {
                name: "employee",
                port: "8080:80",
                context: "employee",
                definitions: [
                    {
                        path: "/",
                        method: "GET"
                    }
                ]
            },
            envVar: "employeegw_url"
    },
    {
            targetCell: "Stock-Options",
            ingress: {
                name: "stock",
                port: "8080:80",
                context: "stock",
                definitions: [
                    {
                        path: "/",
                        method: "GET"
                    }
                ]
            },
            envVar: "stockgw_url"
    }
    ];
    _ = cellery:build(hrCell);

}
