# Hello World Sample

This is a simple Hello World sample for Cellery with [VICK](https://github.com/wso2/product-vick/) as the runtime.

**Note:**
[VICK](https://github.com/wso2/product-vick/) must be installed to try this out. You can configure Cellery to use a 
VICK installed Kubernetes Cluster by using following command and selecting the particular Kubernetes cluster.

```
$ cellery configure
``` 

## Building the sample
To build the sample, execute the following command.
```
$ cellery build hello-world.bal
``` 

## Running the sample
To run the sample, execute the following command.
```
$ cellery run wso2/hello-world:1.0.0
``` 

To view the status of the cell, execute the following command.
```
$ cellery ps
``` 
Above command will output the following result if all the cells are ready;

| NAME  | STATUS | GATEWAY |   SERVICES | AGE |
| ------ | ------  | ------ | ------ | ------ |
| hello-world | Ready | hello-world--gateway-service | 1 | 3m |

## Invoking the sample

1.  Find your external cluster IP and the port by running the following command (VICK currently uses nginx as the ingress controller),

    ```
	   kubectl get svc -n ingress-nginx ingress-nginx
	```
2.  Update the host entries to point the below hosts to the above IP. 

    ```
	   <External cluster IP>       wso2-apim-gateway
	   <External cluster IP>       wso2-apim
	```
4.  API in the Hello World Cell will be exposed via Global API manager. Therefore log into API Manager Store 
(i.e `https://wso2-apim:<port>/store`, replace the `<port>` with the value you got from #1 above)
as a Subscriber (i.e `admin:admin` in a default setup) and subscribe to the API named `hello_world_global_1_0_0_hello`. 

5. Generate an Access Token following [this document](https://docs.wso2.com/display/AM260/Working+with+Access+Tokens#WorkingwithAccessTokens-Generatingapplicationaccesstokens)

6. To invoke the Hello World API, execute the following curl command. 
Replace the `<port>` with the value you got from #1 above and `<Access_Token>` with the access token value you got from #5 above.
    ```
	   curl -H "Authorization: Bearer <Access_Token>" "https://wso2-apim-gateway:<port>/hello-world/hello/hello/sayHello" -k  
	```