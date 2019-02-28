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

// Pet-store sample's orders microservice. Returns a set of order details. 

var fs = require('fs');
const express = require('express');
const orderListFile = 'orderlist.json';
const app = express();
const port = 80;

let obj = null;
fs.readFile(orderListFile, 'utf8', function (err, data) {
    if (!err) {
        obj = JSON.parse(data);
    }
    else {
        obj = JSON.parse("{error detected:error}");
    }

});

app.get('/orders', (req, res) => res.send(obj));
app.listen(port, () => console.log(`Order-service app listening on port ${port}!`));


