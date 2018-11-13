// Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
//
// WSO2 Inc. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

# Error struct for encapsulating artifact handling errors.
# + message - Error message.
# + code - The error code. This is not an http status code.
# + regError - Inner registry error.
# + cause -  The error(s) that cause the registry error.
type RegistryError object {
    public string message;
    public string code = "";
    public RegistryError? regError;
    public error? cause;
    new (message, code, regError = (), cause = ()) {}

    public function toString() returns @untainted (string) {
        string toStringValue = string `message: {{message}}. code: {{code}} `;
        match regError {
            () => {}
            RegistryError innerRegError => {
                toStringValue = string `{{toStringValue}}
                    - {{innerRegError.toString()}}`;
            }
        }
        match self.cause {
            () => {}
            error bError => {
                toStringValue = string `{{toStringValue}}
                        - {{errorToString(bError)}}`;
            }
        }
        return toStringValue;
    }
};

public function errorToString(error bError) returns @untainted (string) {
    string toStringValue = bError.message;
    match bError.cause {
        error innerBError => {
            toStringValue = string `{{toStringValue}}
                - {{errorToString(innerBError)}}`;
        }
        () => {}
    }
    return toStringValue;
}
