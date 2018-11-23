#!/usr/bin/env bash
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


function showUsageAndExit() {
    echo "USAGE$"
    echo "  $ ./build.sh [OPTIONS]"
    echo
    echo "OPTIONS"
    echo "  -v\t(OPTIONAL) Build version. If not specified a default value (0.0.1) will be used."
    echo
    echo "EXAMPLES"
    echo "  $ ./build.sh -v 1.0.0 \t-  Builds Cellery Registry for version 1.0.0"
    echo
    exit 1
}

function showErrorAndExit() {
    if [[ -z "$1" ]]; then
        echo "[FATAL] Error message is required"
        exit 1
    fi
    if [[ -z "$2" ]]; then
        echo "[FATAL] Error status is required"
        exit 1
    fi
    echo "[ERROR]$1"
    exit $2
}

while getopts :v:h FLAG; do
  case ${FLAG} in
    v)
      BUILD_VERSION=$OPTARG
      ;;
    h)
      showUsageAndExit
      ;;
    \?)
      showUsageAndExit
      ;;
  esac
done

if [[ -z "$BUILD_VERSION" ]]; then
    echo "Build version is not specified. Default version (0.0.1) will be used."
    BUILD_VERSION="0.0.1"
fi


echo
echo "Building Docker images ..."
docker build --no-cache -f Dockerfile -t cellery/cellery-registry:${BUILD_VERSION} .

STATUS=$?

if [[ ${STATUS} != 0 ]]; then
    showErrorAndExit "Docker image build failed." ${STATUS}
fi

echo
echo "Building Cellery Registry Completed"
