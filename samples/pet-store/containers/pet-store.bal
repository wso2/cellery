import ballerina/io;
import celleryio/cellery;

// Pet-store Components
// Orders Component
cellery:Component ordersComponent = {
    name: "orders",
    source: {
        image: "celleryio/samples-pet-store-orders"
    },
    ingresses: {
        orders: new cellery:HTTPIngress(80,
            "orders-svc",
            [
                {
                    path: "/getOrders",
                    method: "GET"
                }
            ]
        )
    }
};

// Customers Component
cellery:Component customersComponent = {
    name: "customers",
    source: {
        image: "celleryio/samples-pet-store-customers"
    },
    ingresses: {
        customers: new cellery:HTTPIngress(80,
            "customers-svc",
            [
                {
                    path: "/getCustomers",
                    method: "GET"
                }
            ]
        )
    }
};

// Catalog Component
cellery:Component catalogComponent = {
    name: "customers",
    source: {
        image: "celleryio/samples-pet-store-catalog"
    },
    ingresses: {
        customers: new cellery:HTTPIngress(80,
            "catalog-svc",
            [
                {
                    path: "/getCustomers",
                    method: "GET"
                }
            ]
        )
    }
};

// Controller Component
cellery:Component controllerComponent = {
    name: "controller",
    source: {
        image: "celleryio/samples-pet-store-controller"
    },
    ingresses: {
        employee: new cellery:HTTPIngress(
                      80,
                      "controller",
                      "./resources/employee.swagger.json"
        )
    },
    parameters: {
        CATALOG_HOST: new cellery:Env(),
        CATALOG_PORT: new cellery:Env(default=80),
        ORDER_HOST: new cellery:Env(),
        ORDER_PORT: new cellery:Env(default=80),
        CUSTOMER_HOST: new cellery:Env(),
        CUSTOMER_PORT: new cellery:Env(default=80)
    }
};

cellery:CellImage petStoreCell = new();

public function build(string orgName, string imageName, string imageVersion) {
    //Build Pet-store Cell
    io:println("Building Orders Cell ...");

    petStoreCell.addComponent(controllerComponent);
    petStoreCell.addComponent(ordersComponent);
    petStoreCell.addComponent(catalogComponent);
    petStoreCell.addComponent(customersComponent);

    cellery:setParameter(controllerComponent.parameters.CATALOG_HOST, cellery:getHost(imageName, catalogComponent));
    cellery:setParameter(controllerComponent.parameters.ORDER_HOST, cellery:getHost(imageName, ordersComponent));
    cellery:setParameter(controllerComponent.parameters.CUSTOMER_HOST, cellery:getHost(imageName, customersComponent));

    //Expose API from Cell Gateway
    petStoreCell.exposeAPIsFrom(controllerComponent);
    petStoreCell.exposeAPIsFrom(catalogComponent);
    petStoreCell.exposeAPIsFrom(ordersComponent);
    petStoreCell.exposeAPIsFrom(customersComponent);

    //Expose API globally
    petStoreCell.exposeGlobalAPI(controllerComponent);
    _ = cellery:createImage(petStoreCell, orgName, imageName, imageVersion);
}
