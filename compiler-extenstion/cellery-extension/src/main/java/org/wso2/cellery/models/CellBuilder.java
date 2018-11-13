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

import io.fabric8.kubernetes.api.model.ObjectMeta;

/**
 * API model.
 */
public final class CellBuilder {
    private String apiVersion;
    private String kind;
    private ObjectMeta metadata;
    private CellSpec spec;

    public CellBuilder() {
    }

    public static CellBuilder aCell() {
        return new CellBuilder();
    }

    public CellBuilder withApiVersion(String apiVersion) {
        this.apiVersion = apiVersion;
        return this;
    }

    public CellBuilder withKind(String kind) {
        this.kind = kind;
        return this;
    }

    public CellBuilder withMetadata(ObjectMeta metadata) {
        this.metadata = metadata;
        return this;
    }

    public CellBuilder withSpec(CellSpec spec) {
        this.spec = spec;
        return this;
    }

    public Cell build() {
        Cell cell = Cell.getInstance();
        cell.setApiVersion(apiVersion);
        cell.setKind(kind);
        cell.setMetadata(metadata);
        cell.setSpec(spec);
        return cell;
    }
}
