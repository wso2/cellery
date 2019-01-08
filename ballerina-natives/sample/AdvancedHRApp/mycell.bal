import ballerina/io;
import ballerina/config;
import celleryio/cellery;

// Salary Component
cellery:Component salary = {
    name: "Salary",
    source: {
        image: "docker.io/wso2vick/sampleapp-salary"
    },
    ingresses: [
        {
            name: "salary",
            port: "8080:80",
            context: "salary",
            definitions: [
                {
                    path: "*",
                    method: "GET,POST,PUT,DELETE"
                }
            ]
        }
    ]
};

// Employee Component
cellery:Component employee = {
    name: "Employee",
    source: {
        image: "docker.io/wso2vick/sampleapp-employee"
    },
    env: { SALARYSVC_URL: "" },
    ingresses: [
        {
            name: "employee",
            port: "8080:80",
            context: "employee",
            definitions: [
                {
                    path: "*",
                    method: "GET,POST,PUT,DELETE"
                }
            ]
        }
    ],
    egresses: [
        {
            parent:salary.name,
            ingress: salary.ingresses[0],
            envVar: "SALARYSVC_URL"
        }
    ]
};

// MySQL Component
cellery:Component mysql = {
    name: "MySQL",
    source: {
        image: "mysql:5.7.24"
    },
    env: { MYSQL_ROOT_PASSWORD: "" },
    ingresses: [
        {
            name: "mysql",
            port: "3306:3306"
        }
    ]
};

// Stocks Component
cellery:Component stocks = {
    name: "Stocks",
    source: {
        image: "docker.io/wso2vick/sampleapp-stocks"
    },
    env: { MYSQL_HOST: "", MYSQL_USER: "", MYSQL_PW: "" },
    ingresses: [
        {
            name: "stocks",
            port: "8080:80",
            context: "stocks",
            definitions: [
                {
                    path: "*",
                    method: "GET,POST,PUT,DELETE"
                }
            ]
        }
    ]
};

// HR Component
cellery:Component hr = {
    name: "HR",
    source: {
        image: "docker.io/wso2vick/sampleapp-hr"
    },
    env: { "EMPLOYEEGW_URL": "", "STOCKGW_URL": "" },
    ingresses: [
        {
            name: "info",
            port: "8080:80",
            context: "info",
            definitions: [
                {
                    path: "*",
                    method: "GET,POST,PUT,DELETE"
                }
            ]
        }
    ]
};

// Cells
cellery:CellImage hrCell = new("HRCell");
cellery:CellImage employeeCell = new("EmployeeCell");
cellery:CellImage stocksCell = new("StocksCell");
cellery:CellImage mysqlCell = new("MySQLCell");

// Build Function
public function celleryBuild() {

    //Employee Cell
    io:println("Building Employee Cell ...");
    employee.replicas = config:getAsInt("employee.EMPLOYEE_REPLICAS");
    salary.replicas = config:getAsInt("salary.SALARY_REPLICAS");
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

    //MySQL cell
    io:println("Building MySQL Cell ...");
    mysql.env["MYSQL_ROOT_PASSWORD"] = config:getAsString("mysql.MYSQL_ROOT_PASSWORD");
    mysqlCell.addComponent(mysql);
    _ = cellery:build(mysqlCell);

    //Stocks Cell
    io:println("Building Stocks Cell ...");
    stocks.replicas = config:getAsInt("stocks.STOCKS_REPLICAS");
    stocks.env["MYSQL_USER"] = config:getAsString("stocks.MYSQL_USER");
    stocks.env["MYSQL_PW"] = config:getAsString("stocks.MYSQL_PW");
    stocksCell.addComponent(stocks);
    stocksCell.apis = [
        {
            parent:stocks.name,
            context: stocks.ingresses[0],
            global: false
        }
    ];
    stocksCell.egresses = [
        {
            targetCell:mysqlCell.name,
            ingress: mysql.ingresses[0],
            envVar: "MYSQL_HOST"
        }
    ];
    _ = cellery:build(stocksCell);

    io:println("Building HR Cell ...");
    hr.replicas = config:getAsInt("hr.HR_REPLICAS");
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
            envVar: "EMPLOYEEGW_URL",
            resiliency: {
                retryConfig: {
                    interval: 100,
                    count: 10,
                    backOffFactor: 0.5,
                    maxWaitInterval: 20000
                }
            }
        },
        {
            targetCell: stocksCell.name,
            ingress: stocks.ingresses[0],
            envVar: "STOCKGW_URL",
            resiliency: {
                retryConfig: {
                    interval: 100,
                    count: 10,
                    backOffFactor: 0.5,
                    maxWaitInterval: 20000
                }
            }
        }
    ];
    _ = cellery:build(hrCell);

}