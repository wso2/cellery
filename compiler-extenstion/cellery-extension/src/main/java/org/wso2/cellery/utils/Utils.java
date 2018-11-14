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

package org.wso2.cellery.utils;

import java.io.File;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.nio.file.StandardOpenOption;

/**
 * Utility methods for Cellery extension.
 */
public class Utils {
    /**
     * Write content to a File. Create the required directories if they don't not exists.
     *
     * @param context        context of the file
     * @param outputFileName target file path
     * @throws IOException If an error occurs when writing to a file
     */
    public static void writeToFile(String context, String outputFileName) throws IOException {
        File newFile = new File(outputFileName);
        // append if file exists
        if (newFile.exists()) {
            Files.write(Paths.get(outputFileName), context.getBytes(StandardCharsets.UTF_8), StandardOpenOption.APPEND);
            return;
        }
        //create required directories
        if (newFile.getParentFile().mkdirs()) {
            Files.write(Paths.get(outputFileName), context.getBytes(StandardCharsets.UTF_8));
            return;
        }
        Files.write(Paths.get(outputFileName), context.getBytes(StandardCharsets.UTF_8));
    }

//    public static <T> String toYaml(T object) {
//        try (StringWriter stringWriter = new StringWriter()) {
//            YamlWriter writer = new YamlWriter(stringWriter);
//            writer.write(object);
//            writer.getConfig().writeConfig.setWriteRootTags(false); //replaces only root tag
//            writer.close(); //don't add this to finally, because the text will not be flushed
//            return removeTags(stringWriter.toString());
//        } catch (IOException e) {
//            throw new RuntimeException(e);
//        }
//    }
//
//    private static String removeTags(String string) {
//        //a tag is a sequence of characters starting with ! and ending with whitespace
//        return removePattern(string, " ![^\\s]*");
//    }
//
//    public static void main(String[] args) {
//        Cell cell = Cell.getInstance();
//        cell.setApiVersion("v1");
//        cell.setKind("Cell");
//        cell.setMetadata(new ObjectMetaBuilder().addToLabels("DD", "bb").withName("myCell").build());
//        List<APIDefinition> apiDefinitions = new ArrayList<>();
//        apiDefinitions.add(APIDefinitionBuilder.anAPIDefinition().withMethod("GET").withPath("/").build());
//        List<API> apis = new ArrayList<>();
//        apis.add(APIBuilder.anAPI().withBackend("blah")
//                .withContext("/")
//                .withGlobal(false)
//                .withDefinitions(apiDefinitions)
//                .build());
//
//
//        CellSpec cellSpec = new CellSpec();
//        cellSpec.setGatewayTemplate(GatewayTemplateBuilder.aGatewayTemplate().withSpec(
//                GatewaySpecBuilder.aGatewaySpec().withApis(apis).build()
//        ).build());
//        cell.setSpec(cellSpec);
//    }
}
