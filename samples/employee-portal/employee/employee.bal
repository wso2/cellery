import ballerina/io;
import ballerina/config;
import celleryio/cellery;

//Employee Component
cellery:Component employeeComponent = {
    name: "employee",
    source: {
        image: "docker.io/wso2vick/sampleapp-employee"
    },
    ingresses: {
        employee: {
            port: 8080,
            context: "employee",
            definitions: [
                {
                    path: "/",
                    method: "GET"
                }
            ]
        }
    },
    parameters: {
        SALARY: {
            paramType: "envVar",
            required: true
        },
        PORT: {
            paramType: "envVar",
            required: true
        },
        CONTEXT: {
            paramType: "envVar",
            required: true
        },
        PASSWORD: {
            paramType: "secret",
            path: "/tmp/dbpsw",
            required: false
        }
    }
};

//Salary Component
cellery:Component salaryComponent = {
    name: "salary",
    source: {
        image: "docker.io/wso2vick/sampleapp-salary"
    },
    ingresses: {
        salaryAPI: {
            context: "payroll",
            port: 8080,
            definitions: [
                {
                    path: "/salary",
                    method: "GET"
                }
            ]
        }
    }
};

public cellery:CellImage employeeCell = new("Employee");

public function celleryBuild() {

    // Build EmployeeCell
    io:println("Building Employee Cell ...");

    //Map component dependecies
    cellery:addParameter(employeeComponent.parameters["SALARY"], cellery:getHost(employeeCell, salaryComponent));
    cellery:addParameter(employeeComponent.parameters["PORT"], 8080);
    cellery:addParameter(employeeComponent.parameters["CONTEXT"], cellery:getContext(salaryComponent.ingresses[
            "salaryAPI"]));

    // Add components to Cell
    employeeCell.addComponent(employeeComponent);
    employeeCell.addComponent(salaryComponent);

    //Expose API from Cell Gateway
    employeeCell.exposeAPIsFrom(employeeComponent);

    //io:println(employeeCell);
    _ = cellery:createImage(employeeCell);
}
