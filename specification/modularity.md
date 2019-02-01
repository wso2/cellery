# Modularity & Versioning

Cellery programs can be written in one or more files organized into packages.

## Files

A Cellery file is structured as follows:
```
[import PackageName [version ImportVersionNumber] [as Identifier];]*
import celleryio/cellery;

(ComponentDefinition | CellInitialization)+

build()
```

A sample file would be as follows;
```ballerina
import ballerina/io;
import celleryio/cellery;

//Components

cellery:Component componentA = {
    name: "ComponentA",
    source: {
        image: "docker.io/celleryio/component-a:v1"
    },
    env: { ENV1: "", ENV2: "" },
    ingresses: [
        {
            name: "foo",
            port: "8080:80",
            context: "foo",
            definitions: [
                {
                    path: "*",
                    method: "GET,POST,PUT,DELETE"
                }
            ]
        }
    ]
};

cellery:Component componentB = {
    name: "ComponentB",
    source: {
        image: "docker.io/celleryio/component-b:v1"
    },
    env: { ENV1: "", ENV2: "" },
    ingresses: [
        {
            name: "bar",
            port: "8080:80",
            context: "bar",
            definitions: [
                {
                    path: "*",
                    method: "GET,POST,PUT,DELETE"
                }
            ]
        }
    ],
    egresses: [
        {
            targetComponent: componentA.name,
            ingress: componentA.ingresses[0],
            envVar: "ENV1",
            resiliency: {
                retryConfig: {
                    interval: 100,
                    count: 10,
                    backOffFactor: 0.5,
                    maxWaitInterval: 20000
                }
            }
        }
    ]
};

//Cell
cellery:CellImage cellA = new ("CellA");
cellery:CellImage cellB = new ("CellB");

//Build Function
public function build() {
    componentA.env["ENV1"] = "2VALUE1";
    componentA.env["ENV2"] = "2VALUE2";
    cellA.addComponent(componentA);
    
    componentB.env["ENV2"] = "1VALUE2";
    cellB.addComponent(componentB);

    cellB.egresses = [
        {
            targetComponent: componentA.name,
            ingress: componentA.ingresses[0],
            envVar: "ENV1"
        }
    ];

    cellB.apis = [
        {
            targetComponent: componentB.name,
            context:componentB.ingresses[0],
            global: true
        }
    ];

    _ = cellery:createImage(cellA);
}
```

## Versioning

##
**Next:** [Command Line Interface](cli.md)

