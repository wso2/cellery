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

import CellDiagram from "./components/CellDiagram";
import React from "react";
import * as PropTypes from "prop-types";

const App = ({data}) => {
    const diagramData = {
        cells: [data.name],
        components: data.components.map((component) => (
            {
                cell: data.name,
                name: component
            }
        )),
        dependencyLinks: []
    };

    // Recursively extract dependencies (including transitive dependencies if available)
    const extractData = (cell) => {
        if (cell.dependencies) {
            cell.dependencies.forEach((dependency) => {
                diagramData.cells.push(dependency.name);
                diagramData.dependencyLinks.push({
                    from: cell.name,
                    to: dependency.name
                });

                if (dependency.components) {
                    dependency.components.forEach((component) => {
                        diagramData.components.push({
                            cell: dependency.name,
                            name: component
                        });
                    });
                }

                extractData(dependency);
            });
        }
    };
    extractData(data);

    return (
        <CellDiagram data={diagramData} focusedCell={data.name}/>
    );
};

App.propTypes = {
    data: PropTypes.object.isRequired
};

export default App;
