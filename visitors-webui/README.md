# Visitors Web UI

The code base of visitors react app in this repo references https://github.com/jdob/visitors-webui.

The dependent repo which implements visitors backend service is [visitors-service](../visitors-service/).

## Run server in local

To run `visitors-webui` locally, following instructions below:

1. Follow instructions in [Run visitors-service in local](../visitors-service/README.md#run-server-in-local) to start `visitors-service` first. You should see message `Starting development server at http://0.0.0.0:8000/` output in the comand line interface. Otherwise, you need to troubleshoot before returning here to continue.

1. Export environment variables required by `visitors-webui`:

   ```
   export REACT_APP_BACKEND_HOST=localhost
   export REACT_APP_BACKEND_PORT=8000
   export REACT_APP_TITLE="Test Visitors WebUI Title"
   ```

1. Install nvm, node.js, and npm by following instructions in [this guide](https://docs.microsoft.com/windows/dev-environment/javascript/nodejs-on-wsl#install-nvm-nodejs-and-npm). Once it's done, return here to continue.

   > **Note**
   > Recommend to install the stable LTS version of Node.js, install the latest version is optional.

1. Change directory to `visitors-webui`, then run the following commands to start server:

   ```
   cd my-app && npm install && npm run build
   cd ../api && npm install
   node server.js
   ```

   You will see the following message in the output, which indicates the server is successfully started:
   
   ```
   [HPM] Proxy created: /  -> http://localhost:8000
   [HPM] Proxy rewrite rule created: "^/api" ~> ""
   Server listening on the port::3000
   ``` 

1. Open `http://localhost:3000/` in your browser to visit the app. You should see `Test Visitors WebUI Title` displayed in the web page. Refreshing the page and you will see one entry with `127.0.1.1` as **Service IP** and **Client IP** added to the table. 

## Stop server

Follow steps below to stop both frontend and backend server.

1. Press `CONTROL-C` to quit the `visitors-webui` server.

1. Follow instructions in [Stop visitors-service server](../visitors-service/README.md#stop-server) to stop backend server.

# References

* [github.com/jdob/visitors-webui](https://github.com/jdob/visitors-webui)
* [Build a Node.js Proxy Server in Under 10 minutes!](https://www.twilio.com/blog/node-js-proxy-server)
* [How to Set up a Node.js Express Server for React](https://www.section.io/engineering-education/how-to-setup-nodejs-express-for-react/)
* [Install nvm, node.js, and npm](https://docs.microsoft.com/windows/dev-environment/javascript/nodejs-on-wsl#install-nvm-nodejs-and-npm)
