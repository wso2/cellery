/*
 * Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import AppBar from "@material-ui/core/AppBar";
import CellDiagramView from "./components/CellDiagramView";
import CelleryLogo from "./icons/CelleryLogo";
import Constants from "./constants";
import CssBaseline from "@material-ui/core/CssBaseline/CssBaseline";
import React from "react";
import Toolbar from "@material-ui/core/Toolbar";
import Typography from "@material-ui/core/Typography";
import {MuiThemeProvider, createMuiTheme, withStyles} from "@material-ui/core/styles";
import * as PropTypes from "prop-types";

const styles = {
    root: {
        flexGrow: 1,
        background: "#fff"
    },
    celleryLogo: {
        width: 100,
        marginRight: 10
    },
    title: {
        color: "#464646"
    },
    appBar: {
        boxShadow: "none"
    }
};

// Create the main theme of the App
const theme = createMuiTheme({
    typography: {
        useNextVariants: true
    },
    palette: {
        primary: {
            main: "#43ab01"
        }
    }
});

const App = ({data, classes}) => {
    const getFQN = (cellData) => `${cellData.org}/${cellData.name}:${cellData.ver}`;
    const diagramData = {
        cells: [], // List of all the cell FQNs
        composites: [], // List of all the composite FQNs
        components: Object.keys(data.components).map((component) => (
            {
                cell: getFQN(data),
                name: component
            }
        )),
        metaInfo: {}, // All the additional information of each Cell/Composite
        dependencyLinks: [] // The dependency links between Cells/Composites
    };

    // Adding the main Cell/Composite
    if (data.kind === Constants.Type.CELL) {
        diagramData.cells.push(getFQN(data));
    } else if (data.kind === Constants.Type.COMPOSITE) {
        diagramData.composites.push(getFQN(data));
    } else {
        throw Error(`Unknown type ${data.kind}`);
    }

    // Recursively extract dependencies (including transitive dependencies if available)
    const extractData = (cell) => {
        if (cell.components) {
            Object.values(cell.components).forEach((component) => {
                Object.entries(component.dependencies.cells).forEach(([alias, dependency]) => {
                    const dependencyName = getFQN(dependency);

                    // Adding the dependency to the Cells/Composites list
                    if (cell.kind === Constants.Type.CELL) {
                        if (!diagramData.cells.includes(dependencyName)) {
                            diagramData.cells.push(dependencyName);
                        }
                    } else if (cell.kind === Constants.Type.COMPOSITE) {
                        if (!diagramData.composites.includes(dependencyName)) {
                            diagramData.composites.push(dependencyName);
                        }
                    } else {
                        throw Error(`Unknown type ${cell.kind}`);
                    }

                    // Adding the link from the Cell to the dependency
                    diagramData.dependencyLinks.push({
                        alias: alias,
                        from: getFQN(cell),
                        to: dependencyName
                    });

                    if (dependency.components) {
                        Object.keys(dependency.components).forEach((component) => {
                            const matches = diagramData.components.find(
                                (datum) => datum.cell === dependencyName && datum.name === component);
                            if (!matches) {
                                diagramData.components.push({
                                    cell: dependencyName,
                                    name: component
                                });
                            }
                        });
                    }
                    extractData(dependency);
                });
            });
        }

        if (!diagramData.metaInfo.hasOwnProperty(getFQN(cell))) {
            const cellMetaInfo = {
                type: cell.type,
                ingresses: [],
                componentDependencyLinks: []
            };

            if (cell.components) {
                Object.entries(cell.components).forEach(([componentName, component]) => {
                    component.dependencies.components.forEach((dependentComponent) => {
                        cellMetaInfo.componentDependencyLinks.push({
                            from: `${getFQN(cell)} ${componentName}`,
                            to: `${getFQN(cell)} ${dependentComponent}`
                        });
                    });
                    if (cell.kind === Constants.Type.CELL && component.exposed) {
                        cellMetaInfo.componentDependencyLinks.push({
                            from: `${getFQN(cell)} gateway`,
                            to: `${getFQN(cell)} ${componentName}`
                        });
                    }
                    if (component.ingressTypes) {
                        component.ingressTypes.forEach((ingressType) => {
                            if (!cellMetaInfo.ingresses.includes(ingressType)) {
                                cellMetaInfo.ingresses.push(ingressType);
                            }
                        });
                    }
                });
            }
            diagramData.metaInfo[getFQN(cell)] = cellMetaInfo;
        }
    };
    extractData(data);

    return (
        <MuiThemeProvider theme={theme}>
            <CssBaseline/>
            <div className={classes.root}>
                <AppBar position="static" color="default" className={classes.appBar}>
                    <Toolbar>
                        <CelleryLogo className={classes.celleryLogo} fontSize="large"/>
                        <Typography variant="h6" className={classes.title}>
                                Image View
                        </Typography>
                    </Toolbar>
                </AppBar>
                <CellDiagramView data={diagramData} focusedNode={getFQN(data)}/>
            </div>
        </MuiThemeProvider>
    );
};

App.propTypes = {
    data: PropTypes.object.isRequired,
    classes: PropTypes.object.isRequired
};

export default withStyles(styles)(App);
