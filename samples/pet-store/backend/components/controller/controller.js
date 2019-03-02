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

const express = require('express')
const http = require('http');
const axios = require('axios')

var fs = require('fs');
var obj;
var productObj;
var ordersObj;
var customerObj;

const app = express()
const port = 80

const catalogHost = process.env.CATALOG_HOST
const catalogPort = process.env.CATALOG_PORT
const customerHost = process.env.CUSTOMER_HOST
const customerPort = process.env.CUSTOMER_PORT
const orderHost = process.env.ORDER_HOST
const orderPort = process.env.ORDER_PORT

const catalogUrl = "http://" + catalogHost + ":" + catalogPort  + "/accessories";
const customerUrl = "http://" + customerHost + ":" + customerPort  + "/customers";
const orderUrl = "http://" + orderHost + ":" + orderPort  + "/orders";

console.log(catalogUrl)
console.log(customerUrl)
console.log(orderUrl)

/******************* Fiter Functions *********************/
function findProduct(data, product_id) {
	var catalogArray = data.catalog
    var filtered = catalogArray.filter(item=>item.id==product_id);
    return filtered;
}

function findCustomer(data, customer_id) {
	console.log(customer_id)
	var customerArray = data.customers
    for (var i = 0; i < customerArray.length; i++) {
        if (customerArray[i].customer_id == customer_id) {
            return(customerArray[i]);
        }
    }
}

function findOrder(data, order_id) {
	var orderArray = data.orders
    var filtered = orderArray.filter(item=>item.order_id==order_id);
    return filtered;
}

function findOrderByCustomer(data, customer_id) {
	var orderArray = data.orders
    var filtered = orderArray.filter(a=>a.customer_id==customer_id);
    return filtered;
}

/******************* End Fiter Functions *********************/



/******************* Axios Functions ****************/

const getData = async url => {
	try {
        const response = await axios.get(url);
        return response.data;
    } catch (error) {
        console.log(error);
      }
    };

/******************* End Axios Functions ****************/



/******************* Controller APIs ****************/

//Get product catalog
app.get('/catalog', function (req, res) {
	getData(catalogUrl).then((response) => {
		res.send({
            catalog: response.data
        });
	});
})

// Get product by id
app.get('/catalog/:id',  (req, res) => {
	getData(catalogUrl).then((response) => {
		productObj = findProduct(response, req.params.id);
		res.send(productObj);
	});
});

//Get all orders
app.get('/orders', function (req, res) {
  getData(orderUrl).then((response) => {
     res.send(response);
  });
})

//Get order by id
app.get('/orders/:order_id', function (req, res) {
  getData(orderUrl).then((response) => {
  	 ordersObj = findOrder(response, req.params.order_id);
     res.send(ordersObj);
  });
})

//Get order by id
app.get('/orders/customer/:customer_id', function (req, res) {
  getData(orderUrl).then((response) => {
  	 ordersObj = findOrderByCustomer(response, req.params.customer_id);
     res.send(ordersObj);
  });
})

//Get all customers
app.get('/customers', function (req, res) {
  getData(customerUrl).then((response) => {
     res.send(response.data);
  });
})

// Get customre by id
app.get('/customers/:id',  (req, res) => {
	getData(customerUrl).then((response) => {
		customerObj = findCustomer(response, req.params.id);
		res.send(customerObj);
	});
});

/******************* End Controller APIs ****************/

app.listen(port, () => console.log(`Pet store controller is listening on port ${port}!`))

