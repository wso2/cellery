/*
 *  Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 * http://www.apache.org/licenses/LICENSE-2.0

 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

const express = require("express");
const axios = require("axios");

const service = express();
const port = process.env.SERVICE_PORT || 3004;

const CATALOG_HOST = process.env.CATALOG_HOST;
const CATALOG_PORT = process.env.CATALOG_PORT;
const CUSTOMERS_HOST = process.env.CUSTOMER_HOST;
const CUSTOMERS_PORT = process.env.CUSTOMER_PORT;
const ORDERS_HOST = process.env.ORDER_HOST;
const ORDERS_PORT = process.env.ORDER_PORT;

const CATALOG_SERVICE_URL = "http://" + CATALOG_HOST + ":" + CATALOG_PORT;
const CUSTOMERS_SERVICE_URL = "http://" + CUSTOMERS_HOST + ":" + CUSTOMERS_PORT;
const ORDERS_SERVICE_URL = "http://" + ORDERS_HOST + ":" + ORDERS_PORT;

service.use(express.json());

/**
 * Handle a success response from the API invocation.
 *
 * @param res The express response object
 * @param data The returned data from the API invocation
 */
const handleSuccess = (res, data) => {
    const response = {
        status: "SUCCESS"
    };
    if (data) {
        response.data = data;
    }
    res.send(response);
};

/**
 * Handle an error which occurred during the API invocation.
 *
 * @param res The express response object
 * @param message The error message
 */
const handleError = (res, message) => {
    console.log("[ERROR] " + message);
    res.status(500).send({
        status: "ERROR",
        message: message
    });
};

/**
 * Call an API.
 *
 * @param config Axios configuration to be used
 * @return {Promise<any>} The promise for the data fetch request
 */
const callAPI = (config) => new Promise((resolve, reject) => {
    axios(config)
        .then((response) => {
            const responseBody = response.data;
            if (responseBody.status === "SUCCESS") {
                resolve(responseBody.data);
            } else {
                console.log("[ERROR] Failed to call API " + config.url + " using method " + config.method + " due to "
                    + responseBody.message);
                reject("Failed to fetch data");
            }
        })
        .catch((error) => {
            console.log("[ERROR] Failed to call API " + config.url + " using method " + config.method + " due to "
                + error);
            reject(error);
        });
});

/**
 * Call an API in the Catalog Service.
 *
 * @param endpoint The endpoint to call in the service
 * @param method The HTTP method to use
 * @return {Promise<any>} The promise for the data fetch request
 */
const callCatalogService = (endpoint, method) => {
    return callAPI({
        url: CATALOG_SERVICE_URL + endpoint,
        method: method
    });
};

/**
 * Call an API in the Customers Service.
 *
 * @param endpoint The endpoint to call in the service
 * @param method The HTTP method to use
 * @return {Promise<any>} The promise for the data fetch request
 */
const callCustomersService = (endpoint, method) => {
    return callAPI({
        url: CUSTOMERS_SERVICE_URL + endpoint,
        method: method
    });
};

/**
 * Call an API in the Orders Service.
 *
 * @param endpoint The endpoint to call in the service
 * @param method The HTTP method to use
 * @return {Promise<any>} The promise for the data fetch request
 */
const callOrdersService = (endpoint, method) => {
    return callAPI({
        url: ORDERS_SERVICE_URL + endpoint,
        method: method
    });
};

/*
 * API endpoint for getting a list of accessories available in the catalog.
 */
service.get("/catalog", (req, res) => {
    callCatalogService("/accessories", "GET")
        .then((data) => {
            handleSuccess(res, {
                accessories: data
            });
        })
        .catch((error) => {
            handleError(res, "Failed to fetch catalog data due to " + error);
        });
});

/*
 * API endpoint for getting a list of accessories available in the catalog.
 */
service.get("/orders", (req, res) => {
    callOrdersService("/orders", "GET")
        .then((data) => {
            handleSuccess(res, {
                orders: data
            });
        })
        .catch((error) => {
            handleError(res, "Failed to fetch orders data due to " + error);
        });
});

/*
 * Starting the server
 */
const server = service.listen(port, () => {
    const host = server.address().address;
    const port = server.address().port;

    console.log("[INFO] Pet Store Controller Service listening at http://%s:%s", host, port);
});
