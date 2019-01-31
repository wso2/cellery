# Registry API 

This specifies a RESTful API for Cellery Registry. This will be used by the Cellery API as well as any other client (UI, etc) which will interact with Cellery Registry.

* Version 0.0.1
* URI Scheme
* Base path: /cellery/registry/0.0.1
* Scheme: [HTTPS]

## Resources

#### GET /version
Returns the version of the Cellery registry

#### POST /images/{orgName}/{cellImageName}/{cellImageVersion}
Push a Cellery image to registry

Sample Request : 
```
curl -v -X POST  -F 'file=@hr_cell.zip' 'https://127.0.0.1:9090/registry/0.0.1/images/wso2/hr_app/2.0.0' -k
```


#### GET /images/{orgName}/{cellImageName}/{cellImageVersion}
Pull a Cellery image from registry

Sample Request : 
```
curl -v 'https://127.0.0.1:9090/registry/0.0.1/images/wso2/hr_app/2.0.0' -k --output hr_cell.zip

wget 'https://127.0.0.1:9090/registry/0.0.1/images/wso2/hr_app/2.0.0' --no-check-certificate -O hr_cell.zip
```

#### GET POST PUT DELETE /docker-registry/*
Expose (i.e passthrough) [Docker Registry API](https://docs.docker.com/registry/spec/api/) through the Cellery Registry API


