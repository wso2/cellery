# Cellery Registry Architecture

Following diagram illustrates the basic components of the registry.

<img src="cellery_registry.png" width="40%">

The Cell image will be persisted into a Storage Volume (i.e File system) and a Docker registry will be used to keep 
the docker images used in the Cells. This Docker registry will be fronted by the Cellery Registry API, which will 
block any direct docker pushes to the docker registry used by Cellery. Furthermore, Native Basic Auth will be used to 
make sure that Docker registry is accessed only via the Cellery registry API. This will guarantee the integrity of the 
images pushed into the docker registry 

**(1) (2) (3)** - When pushing an image to the registry, the daemon/cli will send the compressed image file and 
Registry API will verify and validate it and persist it in the Cell Image storage. It will then generate a token which
has to be passed along when pushing the docker images.

**(4)** - If the above call is successful, Daemon/CLI will push the docker images with the token

**(5)** - Registry API will verify the token and allow the images to be pushed

**(6) (7)** - When pulling the images, again it will go through the Registry API, however, it will not require any tokens. 

Registry API can be secured with Basic Auth or JWT if required.

Default Cell Image Storage will be File System. There should be support for different storage drivers from popular cloud service providers such as AWS S3, Google Cloud Storage etc. Storage structure would be as follows;
```
.
├── archives
│   ├── 042cc1aacd630d66f4035164a1bda63e474bab3184b2112da2859fde59750d92
│   │   └── hr_app.zip
│   ├── 0fe03374e85c87008f445ae42276a91869c3c41f80a0fcc8e6b70588c78a22d6
│   │   └── travel_app.zip
│   ├── 15798e37cb39e494644ae5376510514f7103c70bf4552c9b459f9e50ba80db98
│   │   └── hello-cell.zip
│   ├── 5a0d367d43e688301c2fdfc8cc7d43d5657c294096a0233476ebd5f6dbcdb120
│   │   └── hr_app.zip
│   ├── c4223810867a925bcdcc058cf405f4ddb963785297012d07f0038c5ea83f93a8
│   │   └── hr_app.zip
│   └── e91451e14c466b19b52b96b438192c3336db1dbdae58f7522ce89d61333cb5e8
│       └── hello-cell.zip
├── repository
│   ├── foo
│   │   └── hello-cell
│   │       ├── revisions
│   │       │   ├── 15798e37cb39e494644ae5376510514f7103c70bf4552c9b459f9e50ba80db98
│   │       │   └── e91451e14c466b19b52b96b438192c3336db1dbdae58f7522ce89d61333cb5e8
│   │       └── tags
│   │           ├── 1.0.0
│   │           │   ├── current
│   │           │   │   └── link
│   │           │   └── revisions
│   │           │       ├── 15798e37cb39e494644ae5376510514f7103c70bf4552c9b459f9e50ba80db98
│   │           │       └── e91451e14c466b19b52b96b438192c3336db1dbdae58f7522ce89d61333cb5e8
│   │           └── latest
│   │               ├── current
│   │               │   └── link
│   │               └── revisions
│   │                   └── 15798e37cb39e494644ae5376510514f7103c70bf4552c9b459f9e50ba80db98
│   └── wso2
│       ├── hr_app
│       │   ├── revisions
│       │   │   ├── 042cc1aacd630d66f4035164a1bda63e474bab3184b2112da2859fde59750d92
│       │   │   ├── 5a0d367d43e688301c2fdfc8cc7d43d5657c294096a0233476ebd5f6dbcdb120
│       │   │   └── c4223810867a925bcdcc058cf405f4ddb963785297012d07f0038c5ea83f93a8
│       │   └── tags
│       │       ├── 1.0.0
│       │       │   ├── current
│       │       │   │   └── link
│       │       │   └── revisions
│       │       │       └── 042cc1aacd630d66f4035164a1bda63e474bab3184b2112da2859fde59750d92
│       │       ├── 1.0.1
│       │       │   ├── current
│       │       │   │   └── link
│       │       │   └── revisions
│       │       │       ├── 5a0d367d43e688301c2fdfc8cc7d43d5657c294096a0233476ebd5f6dbcdb120
│       │       │       └── c4223810867a925bcdcc058cf405f4ddb963785297012d07f0038c5ea83f93a8
│       │       └── latest
│       │           ├── current
│       │           │   └── link
│       │           └── revisions
│       │               └── 5a0d367d43e688301c2fdfc8cc7d43d5657c294096a0233476ebd5f6dbcdb120
│       └── travel_app
│           ├── revisions
│           │   └── 0fe03374e85c87008f445ae42276a91869c3c41f80a0fcc8e6b70588c78a22d6
│           └── tags
│               └── 1.0.0
│                   ├── current
│                   │   └── link
│                   └── revisions
│                       └── 0fe03374e85c87008f445ae42276a91869c3c41f80a0fcc8e6b70588c78a22d6
└── tmp
```
For each revision of a cell file, SHA-256 hash would be generated and a directory would be created inside `archives` directory with the hash as the name. Particular revision of the cell file will be saved inside that directory.

Repository directory would be organized according to the tag of an image. If the tag name is `foo/hello-cell:1.0.0` and has two revisions, following will be the directory structure;
``` 
├── repository
│   ├── foo
│   │   └── hello-cell
│   │       ├── revisions
│   │       │   ├── 15798e37cb39e494644ae5376510514f7103c70bf4552c9b459f9e50ba80db98
│   │       │   └── e91451e14c466b19b52b96b438192c3336db1dbdae58f7522ce89d61333cb5e8
│   │       └── tags
│   │           ├── 1.0.0
│   │           │   ├── current
│   │           │   │   └── link
│   │           │   └── revisions
│   │           │       ├── 15798e37cb39e494644ae5376510514f7103c70bf4552c9b459f9e50ba80db98
│   │           │       └── e91451e14c466b19b52b96b438192c3336db1dbdae58f7522ce89d61333cb5e8
│   │           └── latest
│   │               ├── current
│   │               │   └── link
│   │               └── revisions
│   │                   └── 15798e37cb39e494644ae5376510514f7103c70bf4552c9b459f9e50ba80db98
```

- `revisions` directory will have a directory for hash values of all the revisions of the cell image **irrespective of the version**.
- `tags` directory will contain a directory for each version of the cell
    - `{version}/current` will contain a file (i.e `link`) with the hash of the current revision of the particular version
    - `{version}/revisions` will have a directory for hash values of all the revisions of the particular cell image version
    
    **Note:** `latest` would be a special version to represent the latest cell image 