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
 * API model.
 */
public final class APIBuilder {
    private String context;
    private String backend;
    private List<APIDefinition> definitions;
    private boolean global;

    private APIBuilder() {
    }

    public static APIBuilder anAPI() {
        return new APIBuilder();
    }

    public APIBuilder withContext(String context) {
        this.context = context;
        return this;
    }

    public APIBuilder withBackend(String backend) {
        this.backend = backend;
        return this;
    }

    public APIBuilder withDefinitions(List<APIDefinition> definitions) {
        this.definitions = definitions;
        return this;
    }

    public APIBuilder withGlobal(boolean global) {
        this.global = global;
        return this;
    }

    public API build() {
        API aPI = new API();
        aPI.setContext(context);
        aPI.setBackend(backend);
        aPI.setDefinitions(definitions);
        aPI.setGlobal(global);
        return aPI;
    }
}
