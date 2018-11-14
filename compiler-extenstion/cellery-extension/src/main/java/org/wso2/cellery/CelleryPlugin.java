/*
 * Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.cellery;

import io.fabric8.kubernetes.api.model.ContainerBuilder;
import io.fabric8.kubernetes.api.model.ContainerPortBuilder;
import io.fabric8.kubernetes.api.model.ObjectMetaBuilder;
import org.ballerinalang.compiler.plugins.AbstractCompilerPlugin;
import org.ballerinalang.compiler.plugins.SupportedAnnotationPackages;
import org.ballerinalang.model.elements.PackageID;
import org.ballerinalang.model.tree.AnnotationAttachmentNode;
import org.ballerinalang.model.tree.PackageNode;
import org.ballerinalang.model.tree.VariableNode;
import org.ballerinalang.util.diagnostic.DiagnosticLog;
import org.wso2.cellery.models.APIBuilder;
import org.wso2.cellery.models.APIDefinitionBuilder;
import org.wso2.cellery.models.Cell;
import org.wso2.cellery.models.CellSpec;
import org.wso2.cellery.models.GatewaySpecBuilder;
import org.wso2.cellery.models.GatewayTemplateBuilder;
import org.wso2.cellery.models.ServiceTemplateBuilder;
import org.wso2.cellery.utils.Utils;

import java.io.File;
import java.io.IOException;
import java.nio.file.Path;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.wso2.cellery.utils.Utils.toYaml;

/**
 * Compiler plugin to generate Cellery artifacts.
 */
@SupportedAnnotationPackages(
        value = "wso2/cellery:0.0.0"
)
public class CelleryPlugin extends AbstractCompilerPlugin {

    @Override
    public void init(DiagnosticLog diagnosticLog) {
    }

    @Override
    public void process(PackageNode packageNode) {
    }

    @Override
    public void process(VariableNode variableNode, List<AnnotationAttachmentNode> annotations) {
        //TODO:Implement to process variable node.
    }

    @Override
    public void codeGenerated(PackageID packageID, Path binaryPath) {
        Cell cell = Cell.getInstance();
        cell.setApiVersion("v1");
        cell.setKind("Cell");
        cell.setMetadata(new ObjectMetaBuilder()
                .addToLabels("app", "test")
                .withName("myCell")
                .build());
        CellSpec cellSpec = new CellSpec();
        cellSpec.setGatewayTemplate(GatewayTemplateBuilder.aGatewayTemplate().withSpec(
                GatewaySpecBuilder.aGatewaySpec()
                        .withApis(
                        Collections.singletonList(APIBuilder.anAPI()
                                .withBackend("mock-backend")
                                .withContext("/")
                                .withGlobal(false)
                                .withDefinitions(Collections.singletonList(
                                        APIDefinitionBuilder.anAPIDefinition()
                                                .withMethod("GET")
                                                .withPath("/")
                                                .build()))
                                .build())
                ).build()
        ).build());
        cellSpec.setServiceTemplates(
                Arrays.asList(
                        ServiceTemplateBuilder.
                                aServiceTemplate()
                                .withReplicas(1)
                                .withServicePort(9090)
                                .withMetadata(new ObjectMetaBuilder()
                                        .addToLabels("app", "myApp1")
                                        .build())
                                .withContainer(new ContainerBuilder()
                                        .withImage("ballerina-service1:v1.0.0")
                                        .withPorts(Collections.singletonList(
                                                new ContainerPortBuilder()
                                                        .withContainerPort(9090)
                                                        .build()))
                                        .build())
                                .build(),
                        ServiceTemplateBuilder.
                                aServiceTemplate()
                                .withReplicas(1)
                                .withServicePort(9091)
                                .withMetadata(new ObjectMetaBuilder()
                                        .addToLabels("app", "myApp2")
                                        .build())
                                .withContainer(new ContainerBuilder()
                                        .withImage("ballerina-service2:v1.0.0")
                                        .withPorts(Collections.singletonList(
                                                new ContainerPortBuilder()
                                                        .withContainerPort(9091)
                                                        .build()))
                                        .build())
                                .build())
        );
        cell.setSpec(cellSpec);
        try {
            Utils.writeToFile(toYaml(cell),
                    binaryPath.toAbsolutePath().getParent() + File.separator + "target" +
                            File.separator + "cellery" + File.separator + "test" + ".yaml");
        } catch (IOException ex) {
        }
    }

}
