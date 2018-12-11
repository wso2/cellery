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

import com.esotericsoftware.yamlbeans.YamlWriter;
import io.cellery.models.API;
import io.cellery.models.APIDefinition;
import io.cellery.models.Cell;
import io.cellery.models.CellSpec;
import io.cellery.models.Component;
import io.cellery.models.ComponentHolder;
import io.cellery.models.GatewaySpec;
import io.cellery.models.GatewayTemplate;
import io.cellery.models.ServiceTemplate;
import io.fabric8.kubernetes.api.model.ContainerBuilder;
import io.fabric8.kubernetes.api.model.ContainerPortBuilder;
import io.fabric8.kubernetes.api.model.ObjectMetaBuilder;
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

import java.io.File;
import java.io.IOException;
import java.io.StringWriter;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Locale;
import java.util.Map;

import static org.apache.commons.lang3.StringUtils.removePattern;

/**
 * Native function cellery/build.
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
        processComponents(
                ((BRefValueArray) ((BMap) ctx.getNullableRefArgument(0)).getMap().get("components")).getValues());
        processAPIs(((BRefValueArray) ((BMap) ctx.getNullableRefArgument(0)).getMap().get("apis")).getValues());
        generateCell(((BMap) ctx.getNullableRefArgument(0)).getMap().get("name").toString());
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

    private void generateCell(String name) {
        List<Component> components =
                new ArrayList<>(ComponentHolder.getInstance().getComponentNameToComponentMap().values());
        GatewaySpec spec = new GatewaySpec();
        List<ServiceTemplate> serviceTemplateList = new ArrayList<>();
        for (Component component : components) {
            spec.setApis(component.getApis());
            ServiceTemplate serviceTemplate = new ServiceTemplate();
            serviceTemplate.setReplicas(component.getReplicas());
            serviceTemplate.setServicePort(component.getServicePort());
            serviceTemplate.setMetadata(new ObjectMetaBuilder().withName(component.getService()).build());
            serviceTemplate.setContainer(new ContainerBuilder()
                    .withImage(component.getSource())
                    .withPorts(new ContainerPortBuilder().
                            withContainerPort(component.getContainerPort())
                            .build())
                    .build());
            serviceTemplateList.add(serviceTemplate);
        }
        GatewayTemplate gatewayTemplate = new GatewayTemplate();
        gatewayTemplate.setSpec(spec);
        CellSpec cellSpec = new CellSpec();
        cellSpec.setGatewayTemplate(gatewayTemplate);
        cellSpec.setServiceTemplates(serviceTemplateList);
        Cell cell = new Cell();
        cell.setMetadata(new ObjectMetaBuilder().withName(name).build());
        cell.setSpec(cellSpec);
        String targetPath = System.getProperty("user.dir") + File.separator + name + ".yaml";
        try {
            writeToFile(toYaml(cell), targetPath);
        } catch (IOException e) {
            throw new BallerinaException(e.getMessage());
        }
    }

    /**
     * Write content to a File. Create the required directories if they don't not exists.
     *
     * @param context    context of the file
     * @param targetPath target file path
     * @throws IOException If an error occurs when writing to a file
     */
    private void writeToFile(String context, String targetPath) throws IOException {
        File newFile = new File(targetPath);
        // delete if file exists
        if (newFile.exists()) {
            Files.delete(Paths.get(newFile.getPath()));
        }
        //create required directories
        if (newFile.getParentFile().mkdirs()) {
            Files.write(Paths.get(targetPath), context.getBytes(StandardCharsets.UTF_8));
            return;
        }
        Files.write(Paths.get(targetPath), context.getBytes(StandardCharsets.UTF_8));
    }

    /**
     * Generates Yaml from a object.
     *
     * @param object Object
     * @param <T>    Any Object type
     * @return Yaml as a string.
     */
    private <T> String toYaml(T object) {
        try (StringWriter stringWriter = new StringWriter()) {
            YamlWriter writer = new YamlWriter(stringWriter);
            writer.write(object);
            writer.getConfig().writeConfig.setWriteRootTags(false); //replaces only root tag
            writer.close(); //don't add this to finally, because the text will not be flushed
            return removeTags(stringWriter.toString());
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }

    private String removeTags(String string) {
        //a tag is a sequence of characters starting with ! and ending with whitespace
        return removePattern(string, " ![^\\s]*");
    }


}


