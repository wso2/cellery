# Employee Portal Sample

Employee Portal sample demonstrates an application which uses the cell-based architecture. This application contains 
four microservices deployed across three different cells based on their responsibilities.   
  
* HR Cell:  
  * Contains one microservice written in ballerina.  
  * This microservice invokes Employee cell and Stock cell to build the employee information and returns the result.  
  * Exposes an API at the global Gateway which can be consumed by any API consumer from outside VICK.  
  
* Employee cell  
   
  * Contains two microservices (employee and salary) written in ballerina.  
  * Employee microservice returns employee details such as name, designation.  
  * Salary microservice returns employee salary details.  
  * Exposes a cell level API, when invoked returns the aggregated response from employee and salary response.  
  
* Stock cell  
   
  * Contains one microservice written in ballerina.  
  * This microservice returns the stock related data for the employee.  
  * Exposes a cell level API, when invoked returns the response from stock microservice. 
  
Following diagram illustrates the sample application flow.  
<p align="center">
  <img src="https://raw.githubusercontent.com/wso2/product-vick/master/samples/employee-portal/src/images/employee-portal-architecture-diagram.png">
</p>  
  
**Note:**
[VICK](https://github.com/wso2/product-vick/) must be installed to try this out. You can configure Cellery to use a 
VICK installed Kubernetes Cluster by using following command and selecting the particular Kubernetes cluster.

```
$ cellery configure
``` 

## Building and Running the sample
Each cell need to be build and run individually. 

* Employee Cell

    To build employee cell, execute the following commands
    ```
    $ cd employee
    $ cellery build employee.bal
    ``` 

    To run employee cell, execute the following command.
    ```
    $ cellery run wso2/employee:1.0.0
    ```
* Stocks Cell

    To build employee cell, execute the following commands
    ```
    $ cd stocks
    $ cellery build stocks.bal
    ``` 

    To run stocks cell, execute the following command.
    ```
    $ cellery run wso2/stocks:1.0.0
    ```  
* HR Cell

    To build HR cell, execute the following commands
    ```
    $ cd hr
    $ cellery build hr.bal
    ``` 

    To run HR cell, execute the following command.
    ```
    $ cellery run wso2/hr:1.0.0
    ```       
    

Above commands will deploy the above three cells in the default namespace. To check the cell status run the following command,
```
$ cellery ps
``` 

It will output the following result if all the cells are ready;

| NAME  | STATUS | GATEWAY |   SERVICES | AGE |
| ------ | ------  | ------ | ------ | ------ |
| employee | Ready | employee--gateway-service | 2 | 3m |
| hr | Ready | hr--gateway-service | 1 | 3m |
| stock-options | Ready | stock-options--gateway-service | 1 | 3m |

## Invoking the application

1.  Find your external cluster IP and the port by running the following command (VICK currently uses nginx as the ingress controller),

    ```
	   kubectl get svc -n ingress-nginx ingress-nginx
	```
2.  Update the host entries to point the below hosts to the above IP

    ```
	   <External cluster IP>       wso2-apim-gateway
	   <External cluster IP>       wso2-apim
	```
3. Go to API Manager Store (`https://wso2-apim:<port>/store`) and click `Sign Up` and register a API consumer named `alice` 
(i.e Username). Populate the other values accordingly. 

4.  API in the HR Cell will be exposed via Global API manager. Therefore log into API Manager Store (`https://wso2-apim:<port>/store`)
as a Subscriber (i.e `admin:admin` in a default setup) and subscribe to the API named `hr_global_1_0_0_info`. 
Get the Consumer Key and Consumer secret of the application used for subscribing.

5. To generate an Access Token, execute the following curl command; Replace the `<port>` with the value you got from #1
above and `<Base64 encoded consumer_key:consumer_secret>` with the values you got from #4 above.

    ```
	   curl -k -d "grant_type=password&username=alice&password=password" -H "Authorization: Basic <Base64 encoded consumer_key:consumer_secret>" https://wso2-apim-gateway:<port>/token  
	```
6. To invoke the HR API, execute the following curl command; Replace the `<port>` with the value you got from #1 above 
and `<Access_Token>` with the access token value you got from #5 above.

    ```
	   curl -k -H "Authorization: Bearer <Access_Token>" https://wso2-apim-gateway:<port>/hr/info/info 
	```
	
##### Sample Response

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