# Update cells

This readme we focus on how the cells can be updated in the runtime without having a considerable downtime. Cellery supports rolling updates, 
and advanced deployments such as blue-green and canary updates.

This README includes,

- [Rolling update](#rolling-update)
- [Blue/Green update](#blue-green-and-canary-update)
- [Canary update](#blue-green-and-canary-update)

## Rolling update
A components within the cell can updated via rolling update. This will terminate the components one-by-one and applying 
the changes, eventually the whole cell will be updated with new version. Let us assume, there is a cell `pet-be` running in the current runtime, 
with cell image `wso2cellery/pet-be-cell:1.0.0`. And now, some changes are made to the application source and hence the users will have to update the running instance. 
As this operation simply updates the currently running instance, the client cells invoking this cell may not be aware of this. Therefore, this should be only permitted when there is no API changes. 
 
Let us assume, the above mentioned `pet-be` cell is going to be updated with cell image `wso2cellery/pet-be-cell:1.0.1`, 
and hence below steps should be followed to perform the rolling update. 
 
1) Update the currently running `pet-be` instance with below command. This will update its components respectively.
```
cellery update pet-be wso2cellery/pet-be-cell:1.0.1
```
3) Now execute `kubectl get pods` and you can see the pods of the `pet-be` are getting initialized. And finally, older pods are getting terminated.
```
$ kubectl get pods
NAME                                             READY   STATUS            RESTARTS   AGE
pet-be--catalog-deployment-54b8cd64-knhnc        0/2     PodInitializing   0          4s
pet-be--catalog-deployment-67b8565469-fq86w      2/2     Running           0          26m
pet-be--controller-deployment-6f89fdb47c-rn4mn   2/2     Running           0          24m
pet-be--controller-deployment-75f5db95f4-2dt96   0/2     PodInitializing   0          4s
pet-be--customers-deployment-7997974649-22hft    2/2     Running           0          26m
pet-be--customers-deployment-7d8df7fb84-h48xs    0/2     PodInitializing   0          4s
pet-be--gateway-deployment-7f787575c6-vmg4p      2/2     Running           0          26m
pet-be--orders-deployment-7d874dfd98-vnhdw       0/2     PodInitializing   0          4s
pet-be--orders-deployment-7d9fd8f5ff-4czdx       2/2     Running           0          26m
pet-be--sts-deployment-7f4f56b5d5-bjhww          3/3     Running           0          26m
pet-fe--gateway-deployment-67ccf688fb-dnhhw      2/2     Running           0          4h6m
pet-fe--portal-deployment-69bb57c466-25nqd       2/2     Running           0          4h6m
pet-fe--sts-deployment-59dbb995c7-g7tc7          3/3     Running           0          4h6m
```

## Blue-Green and Canary update
To update a cell using Blue-Green/ Canary patterns, you should first create a cell instance using the cellery run command. 
Then, the traffic for the specific instance can be switch 100% (Blue-Green) or partially to a new version (canary). 

Let us assume that the pet-be cell is an instance of `wso2cellery/pet-be-cell:1.0.0`, and should be updated with new cell image version ` wso2cellery/pet-be-cell:2.0.0`.
And therefore, as a first step a new pet-be cell instance `pet-be-v2` should be started. Canary deployment can be achieve by having 50% traffic routed to the `pet-be` and `pet-be-v2` 
cell instances. Then, we can  completely switch 100% traffic to `pet-be-v2` and still have the both cell instances running as per the blue-green deployment pattern. 
Finally, terminate `pet-be`.

- Route the 50% of the traffic to the new `pet-be-v2` cell instance. 
```
$ cellery route-traffic pet-be -p pet-be-v2=50
```

- The traffic can be completely switched to 100% to the `pet-be-v2` as shown below. 
```
$ cellery route-traffic pet-be -p pet-be-v2=100
```
- The old instance `pet-be` cell instance, and only have the `pet-be-v2` cell running. 
```
cellery terminate pet-be
```

## Try with sample
[Pet-store application sample](https://github.com/wso2-cellery/samples/tree/master/cells/pet-store) walks through the cell upate scenario. 
Find more information on the steps [here](https://github.com/wso2-cellery/samples/blob/master/docs/pet-store/update-cell.md).

# What's Next?
- [Developing and runing a Cell](writing-a-cell.md) - step by step explanation on how you could define your own cells.
- [CLI Commands](cli-reference.md) - reference for CLI commands.
- [How to code cells?](cellery-syntax.md) - explains how Cellery cells are written.
- [Scale up/down](cell-scaling.md) - scalability of running cell instances with zero scaling and horizontal autoscaler.
- [Observe cells](cellery-observability.md) - provides the runtime insight of cells and components.

