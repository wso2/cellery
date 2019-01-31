# Introduction

## Target Audience: 
- Developers building cloud-native applications, as they go beyond the first 10-20 microservices
- Architects who want to structure their MSA/CNA into an effective application architecture. 

## Pain
Organisations are at a long term cross roads in Enterprise IT as systems shift away from legacy monolithic applications 
and towards cloud-native, microservices and serverless architectures. Fundamentally, though, organisations really only 
want an outcome: adaptivity. Adaptivity means being able to be more agile, to respond to market and competitive demands
and to scale to meet varying real-time demand. 

To do this, it is beginning to be accepted that we need a “composable enterprise”. Unfortunately, there isn’t a clear 
view of what that is. In addition, there are both too many and too few cloud native tools. Too many small projects that 
solve detailed problems. Too few that really address the challenge of become adaptive and composable. 

## What is Cellery?
Cellery is an architectural description language, tooling and extensions for popular cloud-native runtimes that enables 
agile development teams to easily create a composable enterprise in a technology neutral way. 

##How does Cellery work?
The Cellery Language is a subset of [Ballerina](http://ballerina.io) and Ballerina extensions. The language allows developers to define cells
in a simple, effective and compilation-validated way. The standard file extension is a `.cel` file. The language has an 
opinionated view of a composable enterprise, with well-defined concepts of control plane, data plane, ingress/egress 
management and observability. The Cellery Language is completely technology neutral. It supports any container runtimes,
cloud orchestration model, control plane, data plane and gateway. It also enables existing legacy systems as well as 
external systems (including SaaS and partner APIs) to be described as cells. 

The Cellery tooling compiles and validates the cell definitions written in the Cellery Language. The tool generates an 
immutable, versioned, compiled representation of the input language (`.celx`). The `.celx` file refers to external 
dependencies; especially container images and it bakes in the specific image ids so that a given version of a `celx` will 
deploy in a consistent immutable way forever. 

The deployment definitions works for cloud-native runtimes such as Kubernetes and CloudFoundry. It also creates 
well-defined visual representations of the cell. Further tooling will potentially be able to extrapolate the cell 
interface and dependencies completely. The tool is a key part of CI/C.D pipelines as it can also help build external 
dependencies such as Dockerfile or Buildpack containers.

**(TODO : Verify this line)** What about test? Can we define a test for a cell as part of the Cell definition? 

The Cellery language understands that cloud native applications are deployed in multiple environments from integration 
test through to staging and production, and also aids in patterns such as Blue/Green and Canary deployment. 

The Cellery runtime extensions for different cloud native orchestration systems, enable those systems to directly 
understand the concept of a cell. In addition, observability extensions allow a cell to be monitored, logged, and traced. 

Cellery can be used standalone, but over time the project will also create a repository manager for cell definitions 
that can be deployed inside an organisation, run in the cloud and in hybrid scenarios (a.k.a Cellery Central). 

## What would the project contain (initially)?
- A specification of the Cellery Language: Cellery Language is a subset (and/or extension) of Ballerina that allows 
you to define Cells. 
- The Cellery compiler:  This validates and compiles Cellery into deployable artifacts. It also generates graphical 
views of each cell definition. It works with…
- The Cellery runtime extensions for Kubernetes: A set of extensions for K8s including CRDs and Admission Controllers, 
observability plugins etc. 
- Extensions for observability: ….
- A repository for cell definitions (later): e.g. Cellery Central to allow push/pull of cell definitions. 


## How does Cellery compare to existing approaches?

Cellery overlaps with Docker Compose, Helm and Kubernetes YAML (plus….)
Comparison points:
- It only generates deployment artefacts so it works with existing runtimes. Runtime extensions can improve this default 
behaviour.
- Cellery has an opinionated view of what a composable architecture should look like
- Compilation ensures development time validation of links between containers in a way that those other approaches do 
not have.
- It is not YAML

## Organization of the Specification

## Notation

The syntax of Cellery Language is defined informally in this document using a pseudo grammar. Future revisions of this 
document will include more formal grammar notations as well.

The following syntactic conventions are followed in the pseudo grammar used in this document:

- All terminals of the language (keywords) are lowercase words.
- All non-terminals are uppercase words.
- All language keywords are in `this font`.

##
**Next:** [Concepts of Cellery](concepts.md)
