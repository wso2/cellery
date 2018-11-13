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
import org.wso2.cellery.utils.Utils;

import java.io.File;
import java.io.IOException;
import java.nio.file.Path;
import java.util.List;

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
    public void process(VariableNode recordNode, List<AnnotationAttachmentNode> annotations) {

    }

    @Override
    public void codeGenerated(PackageID packageID, Path binaryPath) {
        String testYAML = "apiVersion: vick.wso2.com/v1alpha1\n" +
                "kind: Cell\n" +
                "metadata:\n" +
                "  name: my-cell\n" +
                "spec:\n" +
                "  gatewayTemplate:\n" +
                "    spec:\n" +
                "      apis:\n" +
                "      - context: time\n" +
                "        definitions:\n" +
                "        - path: /\n" +
                "          method: GET\n" +
                "        backend: server-time\n" +
                "        global: true\n" +
                "      - context: hello\n" +
                "        definitions:\n" +
                "        - path: /\n" +
                "          method: GET\n" +
                "        backend: node-hello\n" +
                "        global: true\n" +
                "  servicesTemplates:\n" +
                "  - metadata:\n" +
                "      name: time-us\n" +
                "    spec:\n" +
                "      replicas: 1\n" +
                "      container:\n" +
                "        image: docker.io/mirage20/time-us\n" +
                "        ports:\n" +
                "        - containerPort: 8080\n" +
                "      servicePort: 80\n" +
                "  - metadata:\n" +
                "      name: time-uk\n" +
                "    spec:\n" +
                "      replicas: 1\n" +
                "      container:\n" +
                "        image: docker.io/mirage20/time-uk\n" +
                "        ports:\n" +
                "        - containerPort: 8080\n" +
                "      servicePort: 80\n" +
                "  - metadata:\n" +
                "      name: server-time\n" +
                "    spec:\n" +
                "      replicas: 1\n" +
                "      container:\n" +
                "        image: docker.io/mirage20/time\n" +
                "        ports:\n" +
                "        - containerPort: 8080\n" +
                "      servicePort: 80\n" +
                "  - metadata:\n" +
                "      name: debug\n" +
                "    spec:\n" +
                "      replicas: 1\n" +
                "      container:\n" +
                "        image: docker.io/mirage20/k8s-debug-tools\n" +
                "      servicePort: 80\n" +
                "  - metadata:\n" +
                "      name: node-hello\n" +
                "    spec:\n" +
                "      replicas: 1\n" +
                "      container:\n" +
                "        image: docker.io/nipunaprashan/node-hello\n" +
                "        ports:\n" +
                "        - containerPort: 8080\n" +
                "      servicePort: 80";
        try {
            Utils.writeToFile(testYAML,
                    binaryPath.toAbsolutePath().getParent() + File.separator + "target" +
                            File.separator + "cellery" + File.separator + "test" + ".yaml");
        } catch (IOException ex) {
        }
    }

}
