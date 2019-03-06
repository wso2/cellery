/*
 * Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

const express = require("express");
const fs = require("fs");

const service = express();
const port = process.env.SERVICE_PORT || 3002;
const customersDataFile = "data/customers.json";

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
 * Handle when the requested resource was not found.
 *
 * @param res The express response object
 * @param message The error message
 */
const handleNotFound = (res, message) => {
    res.status(404).send({
        status: "NOT_FOUND",
        message: message
    });
};

/*
 * API endpoint for getting a list of customers available in the catalog.
 */
service.get("/customers", (req, res) => {
    fs.readFile(customersDataFile, "utf8", function (err, data) {
        if (err) {
            handleError(res, "Failed to read data file " + customersDataFile + " due to " + err);
        } else {
            handleSuccess(res, JSON.parse(data));
        }
    });
});

/*
 * API endpoint for creating a new customer.
 */
service.post("/customers/:username", (req, res) => {
    fs.readFile(customersDataFile, "utf8", function (err, data) {
        if (err) {
            handleError(res, "Failed to read data file " + customersDataFile + " due to " + err);
        } else {
            // Creating the new customer data.
            const match = data.filter((customer) => customer.name === req.params.username);
            if (match.length === 0) {
                data.push({
                    ...req.body,
                    name: req.params.username
                });

                // Creating the new customer
                fs.writeFile(customersDataFile, data, "utf8", function (err) {
                    if (err) {
                        handleError(res, "Failed to create new customer due to " + err)
                    } else {
                        handleSuccess(res);
                    }
                });
            } else {
                handleError(res, "Customer " + req.params.username + " already exists");
            }
        }
    });
});

/*
 * API endpoint for getting a single customer from the catalog.
 */
service.get("/customers/:username", (req, res) => {
    fs.readFile(customersDataFile, "utf8", function (err, data) {
        if (err) {
            handleError(res, "Failed to read data file " + customersDataFile + " due to " + err);
        } else {
            let match = JSON.parse(data).filter((customer) => customer.username === req.params.username);
            if (match.length === 1) {
                handleSuccess(res, match[0]);
            } else {
                handleNotFound("Customer not available");
            }
        }
    });
});

/*
 * API endpoint for updating a customer in the catalog.
 */
service.put("/customers/:username", (req, res) => {
    fs.readFile(customersDataFile, "utf8", function (err, data) {
        if (err) {
            handleError(res, "Failed to read data file " + customersDataFile + " due to " + err);
        } else {
            const match = data.filter((customer) => customer.username === req.params.username);

            if (match.length === 1) {
                Object.assign(match[0], req.body);

                // Updating the customer
                fs.writeFile(customersDataFile, data, "utf8", function (err) {
                    if (err) {
                        handleError(res, "Failed to update customer " + req.params.username + " due to " + err)
                    } else {
                        handleSuccess(res);
                    }
                });
            } else {
                handleNotFound("Customer not available");
            }
        }
    });
});

/*
 * API endpoint for deleting a customer in the catalog.
 */
service.delete("/customers/:username", (req, res) => {
    fs.readFile(customersDataFile, "utf8", function (err, data) {
        if (err) {
            handleError(res, "Failed to read data file " + customersDataFile + " due to " + err);
        } else {
            const newData = data.filter((customer) => customer.username !== req.params.username);

            if (newData.length === data.length) {
                handleNotFound("Customer not available");
            } else {
                // Deleting the customer
                fs.writeFile(customersDataFile, newData, "utf8", function (err) {
                    if (err) {
                        handleError(res, "Failed to delete customer " + req.params.username + " due to " + err);
                    } else {
                        handleSuccess(res);
                    }
                });
            }
        }
    });
});

/*
 * Starting the server
 */
const server = service.listen(port, () => {
    const host = server.address().address;
    const port = server.address().port;

    console.log("[INFO] Pet Store Customers Service listening at http://%s:%s", host, port);
});
