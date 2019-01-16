# Cellery native functions
Cellery native function implementation


## How to build

1. Download and install JDK 8 or later
2. Run the Maven command ``mvn clean install`` from within the ``ballerina-natives`` directory.
3. Copy `./ballerina-natives/target/cellery-<version>.jar` file to `$BALLERINA_HOME/bre/lib/` directory.
4. Copy `./ballerina-natives/target/generated-balo/repo/celleryio` directory to `$BALLERINA_HOME/lib/repo/` directory.

## How to run Sample

1. Build and install native functions as per the instructions above.
2. Navigate to `./ballerina-natives/samples/HRApp` directory.
3. Execute the build function.
```bash
$> ballerina run hr_app.bal:build
Building Employee Cell ...
Building Stock Cell ...
Building HR Cell ...
```  
4. Generated yaml files will be located at `./ballerina-natives/samples/HRApp/target` directory.