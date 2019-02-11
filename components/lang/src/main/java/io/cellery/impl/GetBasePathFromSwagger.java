/*
 *   Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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
package io.cellery.impl;

import io.swagger.models.Swagger;
import io.swagger.parser.SwaggerParser;
import org.apache.commons.io.FileUtils;
import org.ballerinalang.bre.Context;
import org.ballerinalang.bre.bvm.BlockingNativeCallableUnit;
import org.ballerinalang.model.types.TypeKind;
import org.ballerinalang.model.values.BString;
import org.ballerinalang.natives.annotations.Argument;
import org.ballerinalang.natives.annotations.BallerinaFunction;
import org.ballerinalang.natives.annotations.ReturnType;
import org.ballerinalang.util.exceptions.BallerinaException;

import java.io.File;
import java.io.IOException;
import java.nio.charset.Charset;
import java.nio.file.Files;
import java.nio.file.Paths;

import static io.cellery.CelleryConstants.RESOURCES;
import static io.cellery.CelleryConstants.TARGET;
import static io.cellery.CelleryUtils.readSwaggerFile;

/**
 * Native function cellery/getBasePathFromSwagger.
 */
@BallerinaFunction(
        orgName = "celleryio", packageName = "cellery:0.0.0",
        functionName = "getBasePathFromSwagger",
        args = {@Argument(name = "swaggerFilePath", type = TypeKind.STRING)},
        returnType = {@ReturnType(type = TypeKind.STRING)},
        isPublic = true
)
public class GetBasePathFromSwagger extends BlockingNativeCallableUnit {

    public void execute(Context ctx) {
        String swaggerFilePath = ctx.getNullableStringArgument(0);
        try {
            String basePath = extractBasePath(swaggerFilePath);
            ctx.setReturnValues(new BString(basePath));
        } catch (IOException e) {
            throw new BallerinaException("Unable to read swagger file. " + swaggerFilePath);
        }
    }

    private String extractBasePath(String swaggerFilePath) throws IOException {
        final Swagger swagger = new SwaggerParser().parse(readSwaggerFile(swaggerFilePath, Charset.defaultCharset()));
        copyResourceToTarget(swaggerFilePath);
        return swagger.getBasePath();
    }

    /**
     * Copy file target/resources directory.
     *
     * @param sourcePath source file/directory path
     * @throws IOException if unable to copy file
     */
    private void copyResourceToTarget(String sourcePath) throws IOException {
        File src = new File(sourcePath);
        String targetPath = TARGET + File.separator + RESOURCES + File.separator + src.getName();
        File dst = new File(targetPath);
        // if source is file
        if (Files.isRegularFile(Paths.get(sourcePath))) {
            if (Files.isDirectory(dst.toPath())) {
                // if destination is directory
                FileUtils.copyFileToDirectory(src, dst);
            } else {
                // if destination is file
                FileUtils.copyFile(src, dst);
            }
        } else if (Files.isDirectory(Paths.get(sourcePath))) {
            FileUtils.copyDirectory(src, dst);
        }
    }
}
