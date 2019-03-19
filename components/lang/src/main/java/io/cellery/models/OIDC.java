/*
 * Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

import lombok.Data;

/**
 * OIDC config model class.
 */
@Data
public class OIDC {
    String provider;
    String clientId;
    String clientSecret;
    String redirectUrl;
    String baseUrl;
    String subjectClaim;

    public OIDC() {
        provider = "https://accounts.google.com";
        clientId = "";
        clientSecret = "";
        redirectUrl = "http://pet-store.com/_auth/callback";
        baseUrl = "http://pet-store.com/items/";
        subjectClaim = "given_name";
    }
}
