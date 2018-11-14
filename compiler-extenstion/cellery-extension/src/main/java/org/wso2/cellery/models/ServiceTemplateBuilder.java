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

import io.fabric8.kubernetes.api.model.Container;
import io.fabric8.kubernetes.api.model.ObjectMeta;

/**
 * Service Template builder.
 */
public final class ServiceTemplateBuilder {
    private ObjectMeta metadata;
    private int replicas;
    private int servicePort;
    private Container container;

    private ServiceTemplateBuilder() {
    }

    public static ServiceTemplateBuilder aServiceTemplate() {
        return new ServiceTemplateBuilder();
    }

    public ServiceTemplateBuilder withMetadata(ObjectMeta metadata) {
        this.metadata = metadata;
        return this;
    }

    public ServiceTemplateBuilder withReplicas(int replicas) {
        this.replicas = replicas;
        return this;
    }

    public ServiceTemplateBuilder withServicePort(int servicePort) {
        this.servicePort = servicePort;
        return this;
    }

    public ServiceTemplateBuilder withContainer(Container container) {
        this.container = container;
        return this;
    }

    public ServiceTemplate build() {
        ServiceTemplate serviceTemplate = new ServiceTemplate();
        serviceTemplate.setMetadata(metadata);
        serviceTemplate.setReplicas(replicas);
        serviceTemplate.setServicePort(servicePort);
        serviceTemplate.setContainer(container);
        return serviceTemplate;
    }
}
