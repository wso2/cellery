
## Employee Portal  
  
Employee Portal sample demonstrates an application which uses the cell-based architecture. This application contains  
four microservices deployed across three different cells based on their responsibilities.   
  
* HR Cell:  
  * Contains one microservice written in ballerina.  
  * This microservice invokes Employee cell and Stock cell to build the employee information and returns the result.  
  * Exposes an API at the global Gateway which can be consumed by any API consumer from outside.  
  
* Employee cell  
   
  * Contains two microservices (employee and salary) written in ballerina.  
  * Employee microservice returns employee details such as name, designation.  
  * Salary microservice returns employee salary details.  
  * Exposes a cell level API, when invoked returns the aggregated response from employee and salary response.  
  
* Stock cell  
   
  * Contains one microservice written in ballerina.  
  * This microservice returns the stock related data for the employee.  
  * Exposes a cell level API, when invoked returns the response from stock microservice.  
  
<p align="center">
  <img src="https://raw.githubusercontent.com/celleryio/sdk/master/samples/employee-portal/src/images/employee-portal-architecture-diagram.png">
</p>
  
  
## Deploying the application   
Once you set up the Cellery runtime on Kubernetes, you can use kubectl to deploy the sample cells into Cellery runtime.  
  
Run following commands to deploy three cells,  

```
kubectl apply -f https://raw.githubusercontent.com/celleryio/sdk/master/samples/employee-portal/artifacts/sample-hr-cell.yaml
kubectl apply -f https://raw.githubusercontent.com/celleryio/sdk/master/samples/employee-portal/artifacts/sample-employee-cell.yaml
kubectl apply -f https://raw.githubusercontent.com/celleryio/sdk/master/samples/employee-portal/artifacts/sample-stock-cell.yaml  
```  
  
The above will deploy the three cells in the default namespace. To check the cell status run the following command,  

```  
kubectl get cells  
```  
  
Which will output the following result if all the cells are ready  
 
| NAME  | STATUS | GATEWAY |   SERVICES | AGE |
| ------ | ------  | ------ | ------ | ------ |
| employee | Ready | employee--gateway-service | 3 | 21s |
| hr | Ready | hr--gateway-service | 2 | 21s |
| stock-options | Ready | stock-options--gateway-service | 2 | 21s |

  
> NOTE: You need to have kubectl version 1.11 or above which support server-side column rendering to see the all 
columns in the above table.


## Invoking the application

1.  Find your external cluster IP by running the following command (Cellery currently uses nginx as the ingress controller),

    ```
	   kubectl get svc -n ingress-nginx ingress-nginx
	```
  
2.  First of all, employee needs to have access to get his employee details. Here, we have registered employee-portal application in Global API manager and generated a token for the employee by providing employee username and password.  
     
3.  Invoke Global API Gateway with the generated token to get employee details.  
      
	As shown in above architecture, cell to cell communication and the intra cell communication is secured by default using JWT.

  

4.  JWT secured invoking call will hit the HR cell gateway (/hr/details) from Global API gateway.  
      
    
5.  The microservice deployed in HR cell will invoke the Employee cell gateway to get employee details and invoke the Stock cell to get stock related details of the employee.  
      
    
6.  Finally, the microservice in HR cell will return the result back to global API gateway.


## Sample Response

```
{
  "employee": {
    "details": {
      "id": "0410",
      "designation": "Senior Software Engineer",
      "salary": "$1500"
    },
    "stocks": {
      "options": {
        "total": 120,
        "vestedAmount": 105
      }
    }
  }
}
```