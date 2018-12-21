# Concepts of Cellery

Following are the key concepts of Cellery:

- [Cell](#cell)
- [Component](#component)        
- [Ingress](#ingress)
- [Egress](#egress)
- [API](#api)
     

#### Cell
A cell is a collection of [components](#component), grouped from design and implementation into deployment. A cell is 
independently deployable, manageable, and observable.

Components inside the cell can communicate with each other using supported transports for intra-cell communication. 
External communication must happen through the edge-gateway or proxy, which provides APIs, events, or streams via 
governed network endpoints using standard network protocols.

A cell can have 1:n components grouped. Components inside the cells are reusable and can be instantiated in multiple cells. 
The cell should document its offers. The capabilities of a cell must be provided by network accessible endpoints. In addition, if 
the cell needs access to external dependencies, then these must also be exposed as network endpoints through a cell-gateway. 
These endpoints can expose APIs, events, or streams. Any interfaces that the microservices or serverless components 
offer that are not made available by the control point should be inaccessible from outside the cell. Every component 
within the cell should be versioned. The cell should have a name and a version identifier. The versions should change 
when the cell’s requirements and/or offers change.

#### Component
A component represents business logic running in a container, serverless environment, or an existing runtime. 
This can then be categorized into many subtypes based on the functional capabilities. A component is designed based on 
a specific scope, which can be independently run and reused at the runtime. Runtime requirements and the behavior of 
the component vary based on the component type and the functional capabilities. The user may decide to build and run 
the code as a service, function, or microservice, or choose to reuse an existing legacy service based on the 
architectural need.

#### Ingress
An Ingress represents an entry point into a Cell or Component. These can expose APIs, events, or streams. Any interfaces 
that the components offer that are not made available by the control point should be inaccessible from outside the cell.

#### Egress
An Egress represents an external communication from a Cell or Component. These can be APIs, events or Streams from other
components or from any external source.  

#### API
An API represents a defined set of functions and procedures that the services of the Components inside a Cell exposes as 
resources (i.e ingresses). These could be made accessible internally or Globally.


##

**Next:** [Language Syntax](language.md)
