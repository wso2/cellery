import ballerina/io;
import celleryio/cellery;

//Stock Component
cellery:Component stock = {
    name: "stock",
    source: {
        image: "docker.io/wso2vick/sampleapp-stock"
    },
    ingresses: [
        {
            name: "stock",
            port: "8080:80",
            context: "stock",
            definitions: [
                {
                    path: "/",
                    method: "GET"
                }
            ]
        }
    ]
};

cellery:Cell stockCell = new("Stock-Options");

public function lifeCycleBuild() {


    //Build Stock Cell
    io:println("Building Stock Cell ...");
    stockCell.addComponent(stock);
    stockCell.apis = [
    {
            parent: stock.name,
            context: stock.ingresses[0],
            global: false
    }
    ];
    _ = cellery:build(stockCell);
}
