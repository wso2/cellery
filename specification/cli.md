#Cellery Command Line Interface

A Command Line Interface (CLI) is implemented to interact with Cellery and its runtime.

Following are the main commands available in the tool;

### Build time commands
|   Command                 | Description                                   |
|           ---             |                   ---                         |
| `cellery configure`       | Configure the cellery installation            |
| `cellery init`            | Initialize a cell project                     |
| `cellery build`           | Build a cell image                            |
| `cellery images`          | List all cell images in the local registry    | 		
| `cellery login`           | Login to the remote cell repository           |
| `cellery tag`             | Tag a cellery image                           | 
| `cellery push`            | Push cell image to the remote repository      |
| `cellery pull`            | Pull cell image from the remote repository    |
| `cellery delete`          | Remove a cell image from the local registry   |
| `cellery describe`        | Describe the cell image                       | 
| `cellery components ls`   | List the components inside the cell           |
| `cellery export`          | Export cell image to .tar archive file        |
| `cellery import`          | Importing cell image from .tar archive file   |
| `cellery descriptor`      | Retrieve cell descriptor from a cell image    | 

### Runtime commands
|   Command                 | Description                                                       |
|           ---             |                           ---                                     |
| `cellery start`		    | Start a previously stopped cell instance                          |      
| `cellery ps`		        | List all running cells                                            |
| `cellery run`			    | Use a cell image to create a  running instance                    |
| `cellery inspect`	        | Describe a running cell instance                                  |
| `cellery rm` 			    | Remove a previously stopped cell instance                         |
| `cellery test`			| Test cell image with dependencies                                 |  
| `cellery stop`			| Stop a cell instance (without removing it)                        |
| `cellery kill`			| Terminate a cell instance                                         |
| `cellery debug`		    | Start the cell in debug mode                                      | 
| `cellery status`      	| Health check of a cell instance                                   |
| `cellery apis`       	    | List the exposed APIs of a cell instance                          |
| `cellery token`		    | Get a token to access cell APIs for a cell instance               |
| `cellery logs`         	| Display logs for a cell or for a particular component of a cell   |
| `cellery exec` 		    | Connect to cell instance component container                      |
| `cellery egress-rules` 	| Retrieve egress rules                                             |
| `cellery ingress-rules`	| Retrieve ingress rules                                            |

##  
