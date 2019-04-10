### Existing setup

Cellery runtime can be installed onto an existing k8s clusters. Cellery runtime requires MySQL server to store the control plane state, and MySQL server can be started within kubernetes with persistent volume mounted or not. 

Hence we have persistence mode and non-persisted mode. Also it requires network or local file system access to store the runtime artifacts of the cellery runtime.

Cellery installer is tested on K8s environments with following component versions.

Kubernetes:  1.11.x 

Docker : docker-ce 18.06.0

MySQL Server: 5.7

NFS Server : NFSv3 compatible


#### Tested Kubernetes Providers

1. [GCP GKE](#gcp-gke)

2. [Docker for Desktop on MacOS](#docker-for-desktop-on-macos) 

3. [Kube Admin](#kube-admin)

4. [Minikube](#minikube)

##### GCP GKE

User needs to provide an NFS share and set of MySQL databases to proceed with the Cellery installation on GCP GKE.

User can find the [SQL script](https://raw.githubusercontent.com/wso2-cellery/distribution/master/installer/k8s-artefacts/mysql/dbscripts/init.sql) file in the Cellery distribution repository

User needs to replace the database use and the password accordingly.

GKE k8s cluster, NFS server and the MySQL server should be hosted in the same GCP compute region and the zone.

##### Docker for Desktop on macOS

User needs to increase the Docker Desktop resources to support Cellery runtime. 

Minimum resources required by Cellery runtime:

 * CPU : 4 core or more
 * Memory : 8 GB or more

Add /var/tmp to the Docker Desktop file sharing to support persistence deployments.  

User may need to restart the Docker Desktop to update the ingress-nginx EXTERNAL-IP after deploying the Cellery runtime.  ( This is a known issue in Docker Desktop.)


##### Kube Admin

Tested on Ubuntu 18.04

Cellery supports both volatile and persistent runtime deployment on kubeadm based k8s. 

##### Minikube

Cellery only supports volatile runtime deployment on Minikube.


#### Cellery runtime deployment with Persistent volume

In this Cellery installation option Cellery CLI uses the default kubernetes cluster configured in the $HOME/.kube/config file.

Cellery needs persistent volume to keep MySQL server files and WSO2 APIM deployable artifacts.


#### Cellery runtime deployment with Persistent volume and access to NFS

If the user has access to an NFS server he can use it as the persistent volume. 


Cellery runtime can be installed on top of the existing K8s cluster with the CLI.

#### Interactive Method

   i. Execute `cellery setup` command to configure Cellery runtime. This 
    will prompt a list of selections. By selecting `create ` section users can setup the Cellery runtime: 
   ```
     $ cellery setup
     [Use arrow keys]
     ? Setup Cellery runtime
         Manage
       ➤ Create
         Modify
         Switch
         EXIT
   ```
   
   ii. From the selections available for environment type, select `Existing cluster` to proceed with local installation:
   ```
     $ ✔ Create
     [Use arrow keys]
     ? Select an environment to be installed
         Local
         GCP
       ➤ Existing cluster
         BACK
   ```
   
##### With Persistent Volume

   i. Once the option `Existing cluster` is selected, the CLI will prompt to select whether to use Persistent Volume or not:

   ```
    $ cellery setup
    ✔ Create
    ✔ Existing cluster
    [Use arrow keys]
    ? Select the type of runtime
     ➤ Persistent volume
       Non persistent volume
       BACK
   ```
   
   ii. Next, select whether you want to use an existing NFS server or not:
   
   ```
    $ cellery setup 
     ✔ Create
     ✔ Existing cluster
     ✔ Persistent volume
     Use the arrow keys to navigate: ↓ ↑ → ←
     ? Use NFS server:
      ▸ Yes
        No
        BACK
   ```
   
   iii. Provide the relevant details and confirm:
   ```
     $ cellery setup
      ✔ Create
      ✔ Existing cluster
      ✔ Persistent volume
      ✔ Yes
      ? NFS server ip:  192.168.2.1
      ? File share name:  data
      ? Database host:  192.168.2.100
      
      ? Mysql credentials required
      Username: mysqlroot
      Password:
      Confirm Password:
   ```

##### Without Persistent Volume

   i. Select the option `Non persistent volume` and continue.
   ```
     $ cellery setup
     ✔ Create
     ✔ Existing cluster
     [Use arrow keys]
     ? Select the type of runtime
         Persistent volume
       ➤ Non persistent volume
   ```
   
   ii. Once the Persistent volume selection is done, select the setup type and continue.

   ```
    $ cellery setup
    ✔ Create
    ✔ Existing cluster
    ✔ Non persistent volume
     [Use arrow keys]
     ? Select the type of runtime
       ➤ Basic
         Complete
   ```

Once the setup is complete, cellery system hostnames should be mapped with the ip of the ingress. 
For this purpose, its currently assumed that the K8s cluster has a nginx ingress controller functioning. 

Run the following kubectl command to get the IP address.

   ```
    kubectl get ingress -n cellery-system
   ```
    
   Then update the /etc/hosts file with that Ip as follows. 
     
   ```
    <IP Address> wso2-apim cellery-dashboard wso2sp-observability-api wso2-apim-gateway cellery-k8s-metrics idp.cellery-system pet-store.com hello-world.com my-hello-world.com
   ```
 

#### Trying Out

Once the installation process is completed, you can try out [quick start with cellery](../../README.md#quick-start-with-cellery).

#### Cleaning Up

Please refer readme for [managing cellery runtimes](./manage-setup.md) for details on how to clean up the setup.
