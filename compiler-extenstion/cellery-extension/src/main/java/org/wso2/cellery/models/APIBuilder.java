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

package org.wso2.cellery.models;

import java.util.List;

/**
 * API Builder.
 */
public final class APIBuilder {
    private API aPI;

    private APIBuilder() {
        aPI = new API();
    }

    public static APIBuilder anAPI() {
        return new APIBuilder();
    }

    public APIBuilder withContext(String context) {
        aPI.setContext(context);
        return this;
    }

    public APIBuilder withBackend(String backend) {
        aPI.setBackend(backend);
        return this;
    }

    public APIBuilder withDefinitions(List<APIDefinition> definitions) {
        aPI.setDefinitions(definitions);
        return this;
    }

    public APIBuilder withGlobal(boolean global) {
        aPI.setGlobal(global);
        return this;
    }

    public API build() {
        return aPI;
    }
}
