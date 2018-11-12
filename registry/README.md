# Cellery Registry

Cellery Registry is a repository to maintain the cell images. It provide the capability to store, discover and distribute cell definitions (i.e images) within the organizations as well as outside users. 

Cellery registry can be deployed as a private registry. Also there will be a central registry managed by WSO2. Following diagram illustrates the basic components of the registry.

<img src="cellery_registry.png" width="40%">

The Cell image will be persisted into a Storage Volume (i.e File system) and a Docker registry will be used to keep the docker images used in the Cells. This Docker registry will be fronted by the Cellery Registry API, which will block any direct docker pushes to the docker registry used by Cellery. Furthermore, Native Basic Auth will be used to make sure that Docker registry is accessed only via the Cellery registry API. This will guarantee the integrity of the images pushed into the docker registry 

**(1) (2) (3)** - When pushing an image to the registry, the daemon will send the compressed image file and Registry API will verify and validate it and persist it in the Cell Image storage. It will then generate a token which has to be passed along when pushing the docker images.
**(4)** - If the above call is successful, Daemon will push the docker images with the token
**(5)** - Registry API will verify the token and allow the images to be pushed

**(6) (7)** - When pulling the images, again it will go through the Registry API, however, it will not require any tokens. 

Registry API can be secured with Basic Auth or JWT if required.

Default Cell Image Storage will be File System. There should be support for different storage drivers from popular cloud service providers such as AWS S3, Google Cloud Storage etc. Storage structure would be as follows;
```
.
├── blobs
│   └── sha256
│       ├── 0addb6fece630456e0ab187b0aa4304d0851ba60576e7f6f9042a97ee908a796
│       │   └── hr-app.cel
│       ├── 18d680d616571900d78ee1c8fff0310f2a2afe39c6ed0ba2651ff667af406c3e
│       │   └── hr-app.cel
│       ├── 473ede7ed136b710ab2dd51579af038b7d00fbbf6a1790c6294c93666203c0a6
│       │   └── hr-app.cel
│       ├── 4a689991aa24aeb0339a27a1b5f42d040f28b0411d37d4812816e79f7049516c
│       │   └── travel-service.cel
│       ├── 5de5f69f42d765af6ffb6753242b18dd4a33602ad7d76df52064833e5c527cb4
│       │   └── travel-service.cel
│       └── 6b9eb699512656fc6ef936ddeb45ab25edcd17ab94901790989f89dbf782344a
│           └── hello-world.cel
└── repository
    ├── foo
    │   └── hello-world
    │       ├── revisions
    │       │   └── sha256
    │       │       └── 6b9eb699512656fc6ef936ddeb45ab25edcd17ab94901790989f89dbf782344a
    │       │           └── link
    │       └── tags
    │           ├── 2.0.0
    │           │   └── link
    │           └── latest
    │               └── link
    └── wso2
        ├── hr-app
        │   ├── revisions
        │   │   └── sha256
        │   │       ├── 0addb6fece630456e0ab187b0aa4304d0851ba60576e7f6f9042a97ee908a796
        │   │       │   └── link
        │   │       ├── 18d680d616571900d78ee1c8fff0310f2a2afe39c6ed0ba2651ff667af406c3e
        │   │       │   └── link
        │   │       └── 473ede7ed136b710ab2dd51579af038b7d00fbbf6a1790c6294c93666203c0a6
        │   │           └── link
        │   └── tags
        │       ├── 1.0.0
        │       │   └── link
        │       ├── 1.0.1
        │       │   └── link
        │       └── latest
        │           └── link
        └── travel-service
            ├── revisions
            │   ├── 4a689991aa24aeb0339a27a1b5f42d040f28b0411d37d4812816e79f7049516c
            │   │   └── link
            │   └── 5de5f69f42d765af6ffb6753242b18dd4a33602ad7d76df52064833e5c527cb4
            │       └── link
            └── tags
                ├── 1.0.0
                │   └── link
                └── 1.0.1
                    └── link
```
For each revision of a cell file, SHA-256 hash would be generated and a directory would be created inside blobs/sha256/ directory with the hash as the name. Particular revision of the cell file will be saved inside that directory.

Repository directory would be organized according to the tag of an image. If the tag name is foo/hello-world:2.0.0, following will be the directory structure;
``` 
└── repository
    ├── foo --- ** Organization/Team/Username **
    │   └── hello-world --- ** Cell Name **
    │       ├── revisions
    │       │   └── sha256
    │       │       └── 6b9eb699512656fc6ef936ddeb45ab25edcd17ab94901790989f89dbf782344a
    │       │           └── link
    │       └── tags
    │           ├── 2.0.0
    │           │   └── link
    │           └── latest
    │               └── link
```

revisions/sha256 directory will have a directory for hash values of all the revisions and inside the directory a file will be create with the name link, which will contain the link to the actual cell file location within blobs directory. 
