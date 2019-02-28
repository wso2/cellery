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

// Pet-store sample's customers microservice. Returns a set of customer details. 

var fs = require('fs');
const express = require('express');
const customersFile = 'customerlist.json';
const app = express();
const port = 80;

app.get('/customers', function (req, res) {
    fs.readFile(customersFile, 'utf8', function (err, data) {
        let obj = null;
        if (!err) {
            obj = JSON.parse(data);
        }
        else {
            obj = JSON.parse("{error detected:error}");
        }

        res.send(obj.customers);
    });
});

app.listen(port, () => console.log(`Customer-service app listening on port ${port}!`));


