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
import org.wso2.ballerinalang.compiler.tree.expressions.BLangRecordLiteral;
import org.wso2.ballerinalang.compiler.tree.expressions.BLangSimpleVarRef;
import org.wso2.cellery.models.API;
import org.wso2.cellery.models.Cell;
import org.wso2.cellery.models.CellSpec;
import org.wso2.cellery.models.ComponentHolder;
import org.wso2.cellery.models.GatewaySpec;
import org.wso2.cellery.models.GatewayTemplate;
import org.wso2.cellery.models.ServiceTemplate;

import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Parse Cell Components.
 */
class CelleryParser {

    /**
     * Create {@link API} object from key value pairs.
     *
     * @param apiKeyValues Record values
     * @return {@link API} object
     */
    static API parseAPI(BLangRecordLiteral.BLangRecordKeyValue apiKeyValues) {
        API api = new API();
        ((BLangRecordLiteral) apiKeyValues.valueExpr).keyValuePairs.forEach(bLangRecordKeyValue -> {
            switch (APIEnum.valueOf(((BLangSimpleVarRef) bLangRecordKeyValue.key.expr).variableName.value)) {
                case context:
                    api.setContext(bLangRecordKeyValue.valueExpr.toString());
                    break;
                case global:
                    api.setGlobal(Boolean.parseBoolean(bLangRecordKeyValue.valueExpr.toString()));
                    break;
                default:
                    break;
            }
        });
        return api;
    }

    /**
     * Create Map object from key value pairs.
     *
     * @param envKeyValues Record values
     * @return Map object
     */
    static Map<String, String> parseEnv(BLangRecordLiteral.BLangRecordKeyValue envKeyValues) {
        Map<String, String> envVars = new HashMap<>();
        ((BLangRecordLiteral) envKeyValues.valueExpr).keyValuePairs.forEach(bLangRecordKeyValue ->
                envVars.put(bLangRecordKeyValue.key.expr.toString(), bLangRecordKeyValue.valueExpr.toString()));
        return envVars;
    }

    static Cell generateCell() {

        List<API> apis = new ArrayList<>();
        List<ServiceTemplate> serviceTemplates = new ArrayList<>();
        ComponentHolder.getInstance().getComponents().forEach(component -> {
            apis.addAll(component.getApis());
            serviceTemplates.add(ServiceTemplate.builder().replicas(component.getReplicas()).build());
        });

        return Cell.builder()
                .metadata(new ObjectMetaBuilder()
                        .addToLabels("app", "test")
                        .withName("my-cell")
                        .build())
                .spec(CellSpec.builder()
                        .gatewayTemplate(
                                GatewayTemplate.builder().spec(GatewaySpec.builder().apis(apis).build()).build())
                        .serviceTemplate(
                                ServiceTemplate.builder()
                                        .replicas(1)
                                        .servicePort(9090)
                                        .metadata(new ObjectMetaBuilder()
                                                .addToLabels("app", "myApp1")
                                                .build())
                                        .container(new ContainerBuilder()
                                                .withImage("ballerina-service1:v1.0.0")
                                                .withPorts(Collections.singletonList(
                                                        new ContainerPortBuilder()
                                                                .withContainerPort(9090)
                                                                .build()))
                                                .build())
                                        .build())
                        .serviceTemplate(
                                ServiceTemplate.builder()
                                        .replicas(1)
                                        .servicePort(9091)
                                        .metadata(new ObjectMetaBuilder()
                                                .addToLabels("app", "myApp2")
                                                .build())
                                        .container(new ContainerBuilder()
                                                .withImage("ballerina-service2:v1.0.0")
                                                .withPorts(Collections.singletonList(
                                                        new ContainerPortBuilder()
                                                                .withContainerPort(9091)
                                                                .build()
                                                        )
                                                )
                                                .build()
                                        )
                                        .build()
                        )
                        .build())
                .build();
    }

    private enum APIEnum {
        context,
        global,
        definitions,
    }
}
