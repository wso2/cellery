import ballerina/io;
import celleryio/cellery;

//Employee Component
cellery:Component employee = {
    name: "employee",
    source: {
        image: "docker.io/wso2vick/sampleapp-employee"
    },
    ingresses: {
        employee: {
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
    }
};

//Salary Component
cellery:Component salary = {
    name: "salary",
    source: {
        image: "docker.io/wso2vick/sampleapp-salary"
    },
    ingresses: {
        salaryAPI: {
            name: "salary",
            port: "8080:80"
        }
    }
};

public cellery:CellImage employeeCell = new("Employee");

public function celleryBuild() {

    // Build EmployeeCell
    io:println("Building Employee Cell ...");
    employeeCell.addComponent(employee);
    employeeCell.addComponent(salary);
    //Expose API from Cell Gateway
    employeeCell.exposeAPIsFrom(employee);
    _ = cellery:createImage(employeeCell);
}
