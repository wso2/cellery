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

import {createMuiTheme} from "@material-ui/core/styles";

const generateTheme = () => createMuiTheme({
    typography: {
        useNextVariants: true
    },
    palette: {
        primary: {
            light: "#E5EAEA",
            main: "#29ABE0",
            contrastText: "#FFF",
            dark: "#1F88b3"
        },
        secondary: {
            main: "#424245"
        }
    }
});

const renderFullPage = (css, content) => (
    "<!DOCTYPE html>" +
    "<html lang='en'>" +
    "<head>" +
    "<meta charset='utf-8'>" +
    "<title>Pet Store</title>" +
    `<style id='jss-server-side'>${css}</style>` +
    "</head>" +
    "<body>" +
    `<div id='app'>${content}</div>` +
    "<script src='./app/index.bundle.js'></script>" +
    "<script src='./app/vendors~index.bundle.js'></script>" +
    "</body>" +
    "</html>"
);

export {
    generateTheme,
    renderFullPage
}
