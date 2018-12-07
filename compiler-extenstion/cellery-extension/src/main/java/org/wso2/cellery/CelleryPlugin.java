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

import org.ballerinalang.compiler.plugins.AbstractCompilerPlugin;
import org.ballerinalang.compiler.plugins.SupportedAnnotationPackages;
import org.ballerinalang.model.elements.PackageID;
import org.ballerinalang.model.tree.AnnotationAttachmentNode;
import org.ballerinalang.model.tree.PackageNode;
import org.ballerinalang.model.tree.VariableNode;
import org.ballerinalang.util.diagnostic.DiagnosticLog;
import org.wso2.ballerinalang.compiler.tree.BLangVariable;
import org.wso2.ballerinalang.compiler.tree.expressions.BLangRecordLiteral;
import org.wso2.ballerinalang.compiler.tree.expressions.BLangSimpleVarRef;
import org.wso2.cellery.models.Cell;
import org.wso2.cellery.models.Component;
import org.wso2.cellery.models.ComponentHolder;
import org.wso2.cellery.utils.Utils;

import java.io.IOException;
import java.io.PrintStream;
import java.nio.file.Path;
import java.util.List;

import static org.wso2.cellery.CelleryParser.generateCell;
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
        PrintStream out = System.out;
        List<BLangRecordLiteral.BLangRecordKeyValue> fields =
                ((BLangRecordLiteral) ((BLangVariable) variableNode).expr).keyValuePairs;
        Component component = new Component();
        for (BLangRecordLiteral.BLangRecordKeyValue keyValue : fields) {
            ComponentKey variableKey = ComponentKey.valueOf(((BLangSimpleVarRef) keyValue.key.expr).variableName.value);
            switch (variableKey) {
                case name:
                    component.setName(keyValue.valueExpr.toString());
                    break;
                case replicas:
                    component.setReplicas(Integer.parseInt(keyValue.valueExpr.toString()));
                    break;
                case apis:
                    component.addApi(CelleryParser.parseAPI(keyValue));
                    break;
                case env:
                    component.setEnvVars(CelleryParser.parseEnv(keyValue));
                    break;
                case source:
                    component.setSource(((BLangRecordLiteral) keyValue.valueExpr).
                            keyValuePairs.get(0).valueExpr.toString());
                    break;
                default:
                    break;
            }
        }
        ComponentHolder.getInstance().addComponent(component);
        out.println(component);
    }

    @Override
    public void codeGenerated(PackageID packageID, Path binaryPath) {
        Cell cell = generateCell();
        try {
            Utils.writeToFile(toYaml(cell), binaryPath);
        } catch (IOException ex) {
        }
    }

    private enum ComponentKey {
        name,
        source,
        replicas,
        container,
        env,
        apis,
        security
    }
}
