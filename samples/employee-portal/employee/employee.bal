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
    },
    ingresses: [
        {
            name: "stock",
            port: "8080:80"
        }
    ]
};

public cellery:Cell employeeCell = new("Employee");

public function lifeCycleBuild() {

    // Build EmployeeCell
    io:println("Building Employee Cell ...");
    employeeCell.addComponent(employee);
    employeeCell.addComponent(salary);
    employeeCell.apis = [
    {
            targetComponent: employee.name,
            context: employee.ingresses[0],
            global: false
    }
    ];
    _ = cellery:build(employeeCell);
}
