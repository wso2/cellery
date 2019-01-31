#!/bin/sh
# --------------------------------------------------------------------
# Copyright (c) 2018, WSO2 Inc. (http://wso2.com) All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# -----------------------------------------------------------------------
#
# This script will initiate the container and start the registry service
#
# -----------------------------------------------------------------------

REGISTRY_DATA_DIRECTORY=/mnt/cellery-registry-data
CONFIG_DIRECTORY=/mnt/cellery-registry-config

# Copy any configuration changes mounted to cellery-registry-config
if test -d ${CONFIG_DIRECTORY}; then
    cp -r ${CONFIG_DIRECTORY}/* ${WORKING_DIRECTORY}/
fi

ballerina run -c registry.toml registry.api.balx 2>&1 | tee cellery-registry.log