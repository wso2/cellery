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
import {MuiThemeProvider, createGenerateClassName} from "@material-ui/core/styles";
import {StaticRouter, matchPath} from "react-router-dom";
import {generateTheme, renderFullPage} from "../utils";
import * as express from "express";
import * as path from "path";
import * as petStoreApi from "../gen/petStoreApi";

const routes = [
    "/",
    "/orders"
];

const renderApp = (req, res, initialState, basePath) => {
    const sheetsRegistry = new SheetsRegistry();
    const sheetsManager = new Map();
    const context = {};
    const app = (
        <JssProvider registry={sheetsRegistry} generateClassName={createGenerateClassName()}>
            <MuiThemeProvider theme={generateTheme()} sheetsManager={sheetsManager}>
                <CssBaseline/>
                <StaticRouter context={context} location={req.url}>
                    <App initialState={initialState}/>
                </StaticRouter>
            </MuiThemeProvider>
        </JssProvider>
    );
    const html = ReactDOMServer.renderToString(app);
    const css = sheetsRegistry.toString();
    res.send(renderFullPage(css, html, initialState, basePath));
};

const createServer = (port) => {
    const app = express();

    app.use("/app", express.static(path.join(__dirname, "/app")));

    /*
     * Serving the App
     */
    app.get("*", (req, res) => {
        const match = routes.reduce((acc, route) => matchPath(req.url, {path: route, exact: true}) || acc, null);

        const initialState = {
            petStoreCell: process.env.PET_STORE_CELL_URL
        };
        const basePath = process.env.BASE_PATH;

        // Setting the Pet Store Cell URL for the Swagger Generated Client
        petStoreApi.setDomain(initialState.petStoreCell);

        if (match) {
            if (match.path === routes[0]) {
                petStoreApi.getCatalog()
                    .then((response) => {
                        let responseBody = response.data;
                        initialState.catalog = {
                            accessories: responseBody.data.accessories
                        };
                        renderApp(req, res, initialState, basePath);
                    })
                    .catch((e) => {
                        console.log("[ERROR] Failed to fetch the catalog due to " + e);
                    });
            } else if (match.path === routes[1]) {
                petStoreApi.getOrders()
                    .then((response) => {
                        const responseBody = response.data;
                        initialState.orders = responseBody.data.orders;
                        renderApp(req, res, initialState, basePath);
                    })
                    .catch((e) => {
                        console.log("[ERROR] Failed to fetch the orders due to " + e);
                    });
            } else {
                renderApp(req, res, initialState, basePath);
            }
        } else {
            renderApp(req, res, initialState, basePath);
        }
    });

    /*
     * Binding to a port and listening for requests
     */
    const server = app.listen(port, () => {
        const host = server.address().address;
        const port = server.address().port;

        console.log("[INFO] Pet Store Portal listening at http://%s:%s", host, port);
    });
};

export default createServer;
