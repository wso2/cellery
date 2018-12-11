/*
 *   Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 *  WSO2 Inc. licenses this file to you under the Apache License,
 *  Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package io.cellery;

import io.cellery.models.API;
import io.cellery.models.APIDefinition;
import io.cellery.models.Cell;
import io.cellery.models.Component;
import io.cellery.models.ComponentHolder;
import org.ballerinalang.bre.Context;
import org.ballerinalang.bre.bvm.BlockingNativeCallableUnit;
import org.ballerinalang.model.types.TypeKind;
import org.ballerinalang.model.values.BMap;
import org.ballerinalang.model.values.BRefType;
import org.ballerinalang.model.values.BRefValueArray;
import org.ballerinalang.model.values.BValue;
import org.ballerinalang.natives.annotations.Argument;
import org.ballerinalang.natives.annotations.BallerinaFunction;
import org.ballerinalang.natives.annotations.ReturnType;
import org.ballerinalang.util.exceptions.BallerinaException;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Locale;
import java.util.Map;

/**
 * Native function cellery/registry:hash.
 */
@BallerinaFunction(
        orgName = "celleryio", packageName = "cellery:0.0.0",
        functionName = "build",
        args = {@Argument(name = "cell", type = TypeKind.OBJECT)},
        returnType = {@ReturnType(type = TypeKind.STRING)},
        isPublic = true
)
public class CelleryBuild extends BlockingNativeCallableUnit {

    public void execute(Context ctx) {
        System.out.println("----------------------------------");
        System.out.println("Build Called....");
        System.out.println(ctx.getNullableRefArgument(0));
        processComponents(
                ((BRefValueArray) ((BMap) ctx.getNullableRefArgument(0)).getMap().get("components")).getValues());
        processAPIs(((BRefValueArray) ((BMap) ctx.getNullableRefArgument(0)).getMap().get("apis")).getValues());
        System.out.println("done");
    }

    private void processComponents(BRefType<?>[] components) {
        for (BRefType componentDefinition : components) {
            if (componentDefinition == null) {
                continue;
            }
            Component component = new Component();
            ((BMap<?, ?>) componentDefinition).getMap().forEach((key, value) -> {
                switch (key.toString()) {
                    case "name":
                        component.setName(value.toString());
                        component.setService(getValidName(value.toString()));
                        break;
                    case "replicas":
                        component.setReplicas(Integer.parseInt(value.toString()));
                        break;
                    case "source":
                        component.setSource(Collections.
                                singletonList(((BMap) value).getMap()).get(0).values().toArray()[0].toString());
                        break;
                    case "ingresses":
                        String portString = processIngressPort(((BRefValueArray) value).getValues());
                        component.setServicePort(Integer.parseInt(portString.substring(portString.indexOf(":") + 1)));
                        component.setContainerPort(Integer.parseInt(portString.substring(0, portString.indexOf(":"))));
                        break;
                    default:
                        break;
                }
            });
            ComponentHolder.getInstance().addComponent(component);
        }
        System.out.println("Holder" + ComponentHolder.getInstance().getComponentNameToComponentMap());
    }

    /**
     * Extract the ingress port
     *
     * @param ingressMap list of ingressee defined
     * @return A string with container port and Service port
     */
    private String processIngressPort(BRefType<?>[] ingressMap) {
        String portString = null;
        for (BRefType ingressDefinition : ingressMap) {
            if (ingressDefinition == null) {
                continue;
            }

            for (Map.Entry<?, ?> entry : ((BMap<?, ?>) ingressDefinition).getMap().entrySet()) {
                Object key = entry.getKey();
                BValue value = (BValue) entry.getValue();
                switch (key.toString()) {
                    case "port":
                        portString = value.toString();
                        break;
                    default:
                        break;
                }
            }
        }
        if (portString == null) {
            throw new BallerinaException("Ingress port is not defined");
        }
        return portString;
    }

    private void processAPIs(BRefType<?>[] apiMap) {
        for (BRefType apiDefinition : apiMap) {
            if (apiDefinition == null) {
                continue;
            }
            API api = new API();
            String componentName = null;
            for (Map.Entry<?, ?> entry : ((BMap<?, ?>) apiDefinition).getMap().entrySet()) {
                Object key = entry.getKey();
                BValue value = (BValue) entry.getValue();
                switch (key.toString()) {
                    case "global":
                        api.setGlobal(Boolean.parseBoolean(value.toString()));
                        break;
                    case "parent":
                        componentName = value.toString();
                        break;
                    case "context":
                        ((BMap<?, ?>) value).getMap().forEach((contextKey, contextValue) -> {
                            switch (contextKey.toString()) {
                                case "context":
                                    api.setContext(contextValue.toString());
                                    break;
                                case "definitions":
                                    api.setDefinitions(processDefinitions(((BRefValueArray) contextValue).getValues()));
                                    break;
                                default:
                                    break;
                            }
                        });
                        break;
                    default:
                        break;

                }
            }
            if (componentName != null) {
                ComponentHolder.getInstance().addAPI(componentName, api);
            } else {
                throw new BallerinaException("Undefined parent component");
            }
        }
        System.out.println("Holder" + ComponentHolder.getInstance().getComponentNameToComponentMap());
    }

    private List<APIDefinition> processDefinitions(BRefType<?>[] definitions) {
        List<APIDefinition> apiDefinitions = new ArrayList<>();
        for (BRefType definition : definitions) {
            if (definition == null) {
                continue;
            }
            APIDefinition apiDefinition = new APIDefinition();
            ((BMap<?, ?>) definition).getMap().forEach((key, value) -> {
                switch (key.toString()) {
                    case "path":
                        apiDefinition.setPath(value.toString());
                        break;
                    case "method":
                        apiDefinition.setMethod(value.toString());
                        break;
                    default:
                        break;

                }
            });
            apiDefinitions.add(apiDefinition);
        }
        return apiDefinitions;
    }

    /**
     * Returns valid kubernetes name.
     *
     * @param name actual value
     * @return valid name
     */
    private String getValidName(String name) {
        return name.toLowerCase(Locale.getDefault()).replace("_", "-").replace(".", "-");
    }

    private void generateCell(){
        Cell cell = new Cell();
    }
}


