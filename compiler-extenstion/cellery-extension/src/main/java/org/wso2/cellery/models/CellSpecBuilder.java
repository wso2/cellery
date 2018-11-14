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
 * Cell Spec Builder.
 */
public final class CellSpecBuilder {
    private GatewayTemplate gatewayTemplate;
    private List<ServiceTemplate> serviceTemplates;

    private CellSpecBuilder() {
    }

    public static CellSpecBuilder aCellSpec() {
        return new CellSpecBuilder();
    }

    public CellSpecBuilder withGatewayTemplate(GatewayTemplate gatewayTemplate) {
        this.gatewayTemplate = gatewayTemplate;
        return this;
    }

    public CellSpecBuilder withServiceTemplates(List<ServiceTemplate> serviceTemplates) {
        this.serviceTemplates = serviceTemplates;
        return this;
    }

    public CellSpec build() {
        CellSpec cellSpec = new CellSpec();
        cellSpec.setGatewayTemplate(gatewayTemplate);
        cellSpec.setServiceTemplates(serviceTemplates);
        return cellSpec;
    }
}
