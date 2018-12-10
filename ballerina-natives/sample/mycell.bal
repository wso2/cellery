import ballerina/io;
import celleryio/cellery;

cellery:Component comp1 = {
    name: "comp1",
    source: {
        image: "docker.io/wso2vick/component1:v1"
    },
    env: ["ENV1", "ENV2"],
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

cellery:Component comp2 = {
    name: "comp2",
    source: {
        image: "docker.io/wso2vick/component2:v1"
    },
    env: ["ENV1", "ENV2"],
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
            name:comp1.ingresses[0],
            env: "ENV1"
        }
    ]
};

cellery:Cell cellA = new;

public function lifeCycleBuild() {
    cellA.addComponent(comp1);
    cellA.addComponent(comp2);
    cellA.egresses = [
        {
            name:comp1.ingresses[0]
        }
    ];
    cellA.apis = [
        {
            context:comp2.ingresses[0],
            global: true
        },
        {
            context: comp1.ingresses[0],
            global: true
        }
    ];

    io:println(cellA);
    var v = cellery:build(cellA);
}

