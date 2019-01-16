import ballerina/io;
import celleryio/cellery;

//Stock Component
cellery:Component stock = {
    name: "stock",
    source: {
        image: "docker.io/wso2vick/sampleapp-stock"
    },
    ingresses: {
        stock: {
            port: 8080,
            context: "stock",
            definitions: [
                {
                    path: "/",
                    method: "GET"
                }
            ]
        }
    }
};

cellery:CellImage stockCell = new("Stock-Options");

public function build() {
    //Build Stock Cell
    io:println("Building Stock Cell ...");
    stockCell.addComponent(stock);
    //Expose API from Cell Gateway
    stockCell.exposeAPIsFrom(stock);
    _ = cellery:createImage(stockCell);
}
