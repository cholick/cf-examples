const express = require('express');
const app = express();

const port = process.env.PORT || 3000;

app.get('/crash', function (req, res) {
    res.end('Crashing backing instance');
    process.exit(1);
});

app.get('*', function (req, res) {
    res.end('Success: Node + Express');
});

var server = app.listen(port, function () {
    console.log(`Listening on port [${port}]`);
});
