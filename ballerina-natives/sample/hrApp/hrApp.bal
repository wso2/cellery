import ballerina/io;
import celleryio/cellery;

//Employee Component
cellery:Component employee = {
    name: "employee",
    source: {
        image: "docker.io/wso2vick/sampleapp-employee"
    },
    ingresses: [
        {
            name: "employee",
            port: "8080:80",
            context: "employee",
            definitions: [
                {
                    path: "/",
                    method: "GET"
                }
            ]
        }
    ]
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

//Employee Cell
cellery:Cell employeeCell = new("EmployeeCell");
cellery:Cell stockCell = new("StockCell");
cellery:Cell hrCell = new("HRCell");

public function lifeCycleBuild() {

    //Set EmployeeCell build time parameters
    io:println("Building Employee Cell ...");
    employeeCell.addComponent(employee);
    employeeCell.addComponent(salary);
    employeeCell.apis = [
        {
            parent:employee.name,
            context: employee.ingresses[0],
            global: false
        }
    ];
    _ = cellery:build(employeeCell);

    //Build Stock Cell
    io:println("Building Stock Cell ...");
    stockCell.addComponent(stock);
    stockCell.apis = [
        {
            parent:stock.name,
            context: stock.ingresses[0],
            global: false
        }
    ];
    _ = cellery:build(stockCell);

    io:println("Building HR Cell ...");
    hrCell.addComponent(hr);
    hrCell.apis = [
        {
            parent:hr.name,
            context: hr.ingresses[0],
            global: true
        }
    ];
    hrCell.egresses = [
        {
            targetCell:employeeCell.name,
            ingress: employee.ingresses[0],
            envVar: "employeegw_url"
        },
        {
            targetCell: stockCell.name,
            ingress: stock.ingresses[0],
            envVar: "stockgw_url"
        }
    ];
    _ = cellery:build(hrCell);

}