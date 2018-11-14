# Cellery Registry

## Overview
Cellery Registry is a repository to maintain the cell images. It provide the capability to store, discover and 
distribute cell definitions (i.e images) within the organizations as well as outside users. 

Cellery registry can be deployed as a private registry. Also there will be a central registry managed by WSO2. 

For more information about Cellery Registry, please see [ARCHITECTURE.md](ARCHITECTURE.md)

## Deploy a registry server

Before you can deploy a Celley Registry, you need to have Docker or Kubernetes installed on the preferred host

#### Run a local registry

You can run the Cellery Registry locally as a docker container with the following command.

```
$ docker run -d -p 9090:9090 --name cellery-registry cellery/cellery-registry:0.0.1
```

To stop the registry, use the same docker container stop command as with any other container.

```
$ docker container stop cellery-registry
```

To remove the container, use docker container rm.

```
$ docker container stop cellery-registry && docker container rm -v cellery-registry
```

## Storage customization

#### Using a Docker Volume

By default, your registry data is persisted as a [docker volume](https://docs.docker.com/storage/volumes/) on the host
 filesystem. 
If you want to use your own docker volume, use the following commands.

```
$ docker volume create cerllery-registry-volume
$ docker run -d -p 9090:9090 --mount source=cellery-registry-volume,target=/mnt/cellery-registry-data --name cellery-registry cellery/cellery-registry:0.0.1
```

#### Using a specific location in Host System

If you want to store your registry data at a specific location on your host filesystem, you can use a bind mount as 
shown below.

```
$ docker run -d -p 9090:9090 -v /tmp/cellery-registry-data:/mnt/cellery-registry-data --name cellery-registry cellery/cellery-registry:0.0.1
```