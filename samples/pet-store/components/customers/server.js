//Import modules for reading files
var fs = require('fs');
const express = require('express');
const customersFile = 'customerlist.json';
const app = express();
const port = 80;

app.get('/getCustomers', function (req, res) {
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


