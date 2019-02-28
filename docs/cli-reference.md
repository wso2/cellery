## Cellery Commands
### Init
Initializes a cellery project. In the command prompt, provide the project name.
Usage : 
```
cellery init
```

### Build
Build an immutable cell image with required dependencies.  
Usage : 
```
cellery build <BAL_FILE_NAME> -t <ORGANIZATION_NAME>/<IMAGE_NAME>:<VERSION>
``` 
Example : 
```
cellery build my-project.bal -t wso2/my-cell:1.0.0
```

### Run
Use a cellery image to create a  running instance.You can also run the cell image with a name by passing -n parameter, 
and this is an optional field. 
Usage : 
```
cellery run <ORGANIZATION_NAME>/<IMAGE_NAME>:<VERSION>  
cellery run <ORGANIZATION_NAME>/<IMAGE_NAME>:<VERSION> -n <NAME>
``` 
Example : 
``` 
cellery run wso2/my-cell:1.0.0 
cellery run wso2/my-cell:1.0.0 -n my-cell-name1
```

### Run linking to dependencies
Use a cellery image to create a running instance specifying an instance name
Input the same instance name to establish the link between the running cell instance and another cell which depends on it.  
Usage : 
```
cellery run <ORGANIZATION_NAME>/<IMAGE_NAME>:<VERSION> -l <NAME-1>
```
Example : 
```
cellery run wso2/cell2:1.0.1 -n mycell-2 -l mycell-1
```

### Ps
List all running cells.  
Usage : 
```
cellery ps
```

### Push
Push cell image to the remote repository.  
Usage : 
```
cellery push <ORGANIZATION_NAME>/<IMAGE_NAME>:<VERSION>
```
Example : 
```
cellery push wso2/my-cell:1.0.0
```

### Pull 
Pulls cell image from the remote repository.  
Usage : 
```
cellery pull <ORGANIZATION_NAME>/<IMAGE_NAME>:<VERSION>
```
Example : 
```
cellery pull wso2/my-cell:1.0.0
```

### Images
List cell images that was pulled and built in the current machine.  
Usage : 
```
cellery images
```

### Stop
Stop a running cell instance.  
Usage : 
```
cellery stop <CELL_NAME>
```
Example : 
```
cellery stop hello-cell
```

### Status 
Performs a health check of a cell.  
Usage : 
```
cellery status <CELL_NAME>
```
Example : 
```
cellery status hello-cell
```

### Apis 
List the exposed APIs of a cell instance.  
Usage : 
```
cellery apis <CELL_NAME>
```
Example : 
```
cellery apis hello-cell
```

### Logs 
Displays logs for either the cell instance, or a component of a running cell instance.  
Usage : 
```
cellery logs <CELL_NAME> / cellery logs <CELL_NAME> -c <COMPONENT_NAME>
```
Example: 
```
cellery logs hello-cell / cellery logs hello-cell -c my-component
```

### Components
Lists the components which the cell encapsulates.  
Usage : 
```
cellery components <CELL_NAME>
```  
Example : 
```
cellery components hello-cell
```

### List-files
List the files (directory structure) of a cell images.  
Usage : 
```
cellery list-files <ORGANIZATION_NAME>/<IMAGE_NAME>:<VERSION>
```
Example : 
```
cellery list-files wso2/my-cell:1.0.0
```

### Extract-resources
Extract the resource files of a pulled image to the provided location. This is useful when cells are packaged with 
swagger files, therefore any components or micro services that uses the cell can generate the code from the swagger.  
Usage : 
```
cellery extract-resources <ORGANIZATION_NAME>/<IMAGE_NAME>:<VERSION> ./resources
```
Example : 
```
cellery extract-resources cellery-samples/employee:1.0.0 ./resources
```
