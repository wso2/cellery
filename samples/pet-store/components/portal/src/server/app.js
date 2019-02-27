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

import App from "../components/App";
import {CssBaseline} from "@material-ui/core";
import {JssProvider} from "react-jss";
import React from "react";
import ReactDOMServer from "react-dom/server";
import {SheetsRegistry} from "jss";
import {StaticRouter} from "react-router-dom";
import * as path from "path";
import {MuiThemeProvider, createGenerateClassName} from "@material-ui/core/styles";
import {generateTheme, renderFullPage} from "../utils";
import * as express from "express";

const createServer = (port) => {
    const app = express();

    app.use("/app", express.static(path.join(__dirname, "/app")));

    /*
     * Serving the App
     */
    app.get("*", (req, res) => {
        const sheetsRegistry = new SheetsRegistry();
        const sheetsManager = new Map();
        const context = {};
        const app = (
            <JssProvider registry={sheetsRegistry} generateClassName={createGenerateClassName()}>
                <MuiThemeProvider theme={generateTheme()} sheetsManager={sheetsManager}>
                    <CssBaseline/>
                    <StaticRouter context={context} location={req.url}>
                        <App/>
                    </StaticRouter>
                </MuiThemeProvider>
            </JssProvider>
        );
        const css = sheetsRegistry.toString();
        res.send(renderFullPage(css, ReactDOMServer.renderToString(app)));
    });

    /*
     * Binding to a port and listening for requests
     */
    const server = app.listen(port, () => {
        const host = server.address().address;
        const port = server.address().port;

        console.log("Pet Store Portal listening at http://%s:%s", host, port);
    });
};

export default createServer;
