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
package io.cellery.registry;

import org.apache.commons.codec.digest.DigestUtils;
import org.ballerinalang.bre.Context;
import org.ballerinalang.bre.bvm.BLangVMErrors;
import org.ballerinalang.bre.bvm.BLangVMStructs;
import org.ballerinalang.bre.bvm.BlockingNativeCallableUnit;
import org.ballerinalang.model.types.TypeKind;
import org.ballerinalang.model.values.BMap;
import org.ballerinalang.model.values.BString;
import org.ballerinalang.model.values.BValue;
import org.ballerinalang.natives.annotations.Argument;
import org.ballerinalang.natives.annotations.BallerinaFunction;
import org.ballerinalang.natives.annotations.ReturnType;

import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.nio.file.Files;
import java.nio.file.NoSuchFileException;
import java.nio.file.Paths;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;

import org.ballerinalang.util.codegen.PackageInfo;
import org.ballerinalang.util.codegen.StructureTypeInfo;
import org.ballerinalang.util.exceptions.BallerinaException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import static org.ballerinalang.util.BLangConstants.BALLERINA_BUILTIN_PKG;

/**
 * Native function cellery/registry:hash.
 */
@BallerinaFunction(
        orgName = "cellery", packageName = "registry:0.0.0",
        functionName = "hash",
        args = {@Argument(name = "path", type = TypeKind.STRING),
                @Argument(name = "algorithm", type = TypeKind.STRING)},
        returnType = {@ReturnType(type = TypeKind.STRING)},
        isPublic = true
)
public class RegistryHash extends BlockingNativeCallableUnit {

    public void execute(Context ctx) {
        String path = ctx.getStringArgument(0);
        BString algorithm = (BString) ctx.getNullableRefArgument(0);
        String hashAlgorithm;

        switch (algorithm.stringValue()) {
            case "SHA1":
                hashAlgorithm = "SHA-1";
                break;
            case "SHA256":
                hashAlgorithm = "SHA-256";
                break;
            case "MD5":
                hashAlgorithm = "MD5";
                break;
            default:
                throw new BallerinaException("Unsupported algorithm " + algorithm + " for HMAC calculation");
        }

        if (null != path && !path.isEmpty()) {
            try {
                byte[] fileContent = Files.readAllBytes(Paths.get(path));
                MessageDigest messageDigest;
                messageDigest = MessageDigest.getInstance(hashAlgorithm);
                messageDigest.update(fileContent);
                byte[] bytes = messageDigest.digest();

                StringBuilder builder = new StringBuilder();
                for (byte aByte : bytes) {
                    builder.append(Integer.toString((aByte & 0xff) + 0x100, 16).substring(1));
                }
                ctx.setReturnValues(new BString(builder.toString()));
            } catch (NoSuchAlgorithmException e) {
                String msg = "Error occurred while generating hash of the file: " + path
                        + ", No Such Algorithm : " + hashAlgorithm;
                ctx.setReturnValues(BLangVMErrors.createError(ctx, msg));
            } catch (NoSuchFileException e) {
                String msg = "Error occurred while generating hash of the file: " + path
                        + ", File Not found.";
                ctx.setReturnValues(BLangVMErrors.createError(ctx, msg));
            } catch (IOException e) {
                String msg = "Error occurred while generating hash of the file: " + path;
                ctx.setReturnValues(BLangVMErrors.createError(ctx, msg));
            }
        } else {
            String msg = "Path cannot be empty";
            ctx.setReturnValues(BLangVMErrors.createError(ctx, msg));
        }
    }
}


