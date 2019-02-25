# Cellery Pet Service Sample with Autoscaling

This is a single cell application which demonstrates the autoscaling of cell components. While Cellery is capable of 
using CPU, memory usage, etc. to autoscale cell components, this sample focuses on the former and hence is written to 
artificially hog cpu for demonstration purposes. 

The component interaction diagram of the application is shown below.

![samples-pet-service](pet-service.png)

## Deploying the sample
1. Run the following command from the sdk root directory:
    ```
    kubectl apply -f ./samples/pet-service/
    ```
2. Wait until Cell gets into ‘Ready’ state. Use the command below to check if the Cell is in Ready state: 
    ```
    kubectl get cells 
    ```

## Invoking the API from the local machine
#### Note: 
By default, the Cellery installation will be using the following domain names,
    ```
    wso2-apim 
    cellery-dashboard 
    wso2sp-observability-api 
    wso2-apim-gateway
    ```
These should be properly mapped to the Kubernetes IngressController IP of the deployment.

1. Login to the API [Store](https://wso2-apim/store/) using admin:admin credentials.

2. Click on ‘pet_service_global_1_0_0_petsvc’ to create a subscription and generate a token. 
(See  [Subscribing to an API](https://docs.wso2.com/display/AM260/Subscribe+to+an+API))

3. Invoke the API using curl

    ```
    curl  https://wso2-apim-gateway/pet-service/petsvc/pets/cats  -H "Authorization: Bearer <access_token>" -k
    ```
    or
    ```    
    curl  https://wso2-apim-gateway/pet-service/petsvc/pets/dogs  -H "Authorization: Bearer <access_token>" -k
    ```

## Simulating high CPU usage

A special debug pod is included in the pet-service cell to induce the effect of high cpu usage. 

1. To simulate high CPU usage, first ssh in to the debug pod using the following command:
   ```
   kubectl exec -it <debug pod id> bash 
   ```
   
2. Issue a a curl command to start high cpu usage in the pet-service component: 
   ```
   curl pet-service--pet-service-service/cpuhogger/start 
   ```
   
   This will start a CPU intensive operation in the pet-service component. As a result of this, the pet-service pod 
   count will get increased. This can be noted by listing the pods in the default namespace, with the following command:
   ```
   kubectl get pods
   ```
   
3. After testing the pet-service as explained in the previous section, to stop the CPU intensive operation, issue the 
following command from the debug pod: 
   ```
   curl pet-service--pet-service-service/cpuhogger/stop
   ```
   
   This will result in the number of pods which belongs to the pet-store service being reduced to one gradually.  
   
   
## Cleaning up
Issue the following command to delete the pet-service sample:
    ```
    kubectl delete -f ./samples/pet-service/
    ```    