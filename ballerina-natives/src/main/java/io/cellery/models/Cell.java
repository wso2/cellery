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

package io.cellery.models;

import io.fabric8.kubernetes.api.model.ObjectMeta;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * Cell POJO model.
 */
@Data
@NoArgsConstructor
public class Cell {
    private String apiVersion;
    private String kind;
    private ObjectMeta metadata;
    private CellSpec spec;

    public Cell(ObjectMeta metadata, CellSpec spec) {
        kind = "Cell";
        apiVersion = "vick.wso2.com/v1alpha1";
        this.metadata = metadata;
        this.spec = spec;
    }
}
