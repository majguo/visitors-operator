const express = require('express');
const path = require('path');
const morgan = require("morgan");
const { createProxyMiddleware } = require('http-proxy-middleware');
const app = express(),
      port = 3000;

// configuration
const BACKEND_SVC_URL = `http://${process.env.REACT_APP_BACKEND_HOST}:${process.env.REACT_APP_BACKEND_PORT}`;

app.use(morgan('dev'));
app.use(express.static(path.join(__dirname, '../my-app/build')));

// Proxy endpoints
app.use('/api/visitors', createProxyMiddleware({
    target: BACKEND_SVC_URL,
    changeOrigin: true,
    pathRewrite: {
        [`^/api`]: '',
    },
 }));

 // Static home page
app.get('/', (req, res) => {
    res.sendFile(path.join(__dirname, '../my-app/build/index.html'));
});

// Home page title endpoint
app.get('/title', (req, res) => {
    console.log(`React App Title: ${process.env.REACT_APP_TITLE}`);
    res.json({
        title: process.env.REACT_APP_TITLE
    });
});

app.listen(port, () => {
    console.log(`Server listening on the port::${port}`);
});