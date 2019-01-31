import ballerina/io;
import celleryio/cellery;

//Debug component
cellery:Component debug = {
    name: "debug",
    source: {
        image: "docker.io/mirage20/k8s-debug-tools"
    },
    ingresses: {
        "debug": {
            port: "80:80"
        }
    }
};


//Employee Component
cellery:Component employee = {
    name: "employee",
    source: {
        image: "docker.io/wso2vick/sampleapp-employee"
    },
    ingresses: {
        "employee": {
            port: "8080:80",
            context: "employee",
            definitions: [
                {
                    path: "/",
                    method: "GET"
                }
            ]
        }
    }
};

//Salary Component
cellery:Component salary = {
    name: "salary",
    source: {
        image: "docker.io/wso2vick/sampleapp-salary"
    }
};

//Stock Component
cellery:Component stock = {
    name: "stock",
    source: {
        image: "docker.io/wso2vick/sampleapp-stock"
    },
    ingresses: {
        "stock": {
            port: "8080:80",
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
cellery:CellImage employeeCell = new("EmployeeCell");
cellery:CellImage stockCell = new("StockCell");
cellery:CellImage hrCell = new("HRCell");

public function build() {

    // Build EmployeeCell
    io:println("Building Employee Cell ...");
    employeeCell.addComponent(employee);
    employeeCell.addComponent(salary);
    employeeCell.addComponent(debug);
    //Expose API from Cell Gateway
    employeeCell.exposeAPIsFrom(employee);
    _ = cellery:createImage(employeeCell);

    //Build Stock Cell
    io:println("Building Stock Cell ...");
    stockCell.addComponent(stock);
    stockCell.addComponent(debug);
    //Expose API from Cell Gateway
    stockCell.exposeAPIsFrom(stock);
    _ = cellery:createImage(stockCell);

    // Build HR cell
    io:println("Building HR Cell ...");
    hrCell.addComponent(hr);
    hrCell.addComponent(debug);

    // Expose API from Cell Gateway & Global Gateway
    hrCell.exposeGlobalAPI(hr);
    hrCell.declareEgress(employeeCell.name, "employeegw_url");
    hrCell.declareEgress(stockCell.name, "stockgw_url");
    _ = cellery:createImage(hrCell);

}
