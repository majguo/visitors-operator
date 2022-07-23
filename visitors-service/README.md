# Visitors Serivce

The code base of visitors service in this repo references https://github.com/jdob/visitors-service.

It's dependent by the visitors web ui which is located at [visitors-webui](../visitors-webui/).

## Run server in local

To run `visitors-service` locally, following instructions below:

1. Check if python3.8 is installed in your local envrionment:
   
   ```
   python3 --version
   ```

1. If not, run the following commands to install python3.8:

   ```
   sudo add-apt-repository ppa:deadsnakes/ppa
   sudo apt-get update
   sudo apt-get install python3.8
   ```

1. If you have other python version installed (e.g., python3.10), select python3.8:

   ```
   update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.8 1
   update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.10 2
	
   update-alternatives --config python3
   ```

   Enter number `1` to select python3.8.

1. Install pip, python3.8-dev and python3.8-venv packages

   ```
   sudo apt install python3-pip
   sudo apt-get install python3.8-dev python3.8-venv
   ```

1. Change directory to `visitors-service`, then run the following commands to start server:

   ```
   python3 -m venv .venv
   source .venv/bin/activate
   
   pip install -r requirements.txt
   ./startup.sh
   ```

   You will see the following message in the output, which indicates the server is successfully started:
   
   ```
   Starting development server at http://0.0.0.0:8000/
   Quit the server with CONTROL-C.
   ``` 

1. Open `http://localhost:8000/visitors/` in the browser and you should see an empty array `[]` displayed in the page.

## Stop server

Follow steps below to stop server and quit the virtual environment.

1. Press `CONTROL-C` to quit the server.

1. Quit the virtual environment:

   ```
   deactivate
   ```

## Run server in Docker

Containerize the application so that it can run as a contanier in Docker or Kubernetes cluster.

### Build application image

Change directory to `visitors-service`, then run the following commands to build a Docker image and push to DockerHub:

```
docker build -t visitors-service:1.0.0 .
docker tag visitors-service:1.0.0 <DockerHub-account>/visitors-service:1.0.0
docker push <DockerHub-account>/visitors-service:1.0.0
```

### Run as container

Execute the command below to run the containerzed application in Docker.

```
docker run -it --rm -p 8000:8000 --name visitors-service visitors-service:1.0.0
```

Open `http://localhost:8000/visitors/` in the browser and you should see an empty array `[]` displayed in the page.

## Stop the container

To stop running the containerzed application, you can press `CONTROL-C` or execute the command below in a separate CLI.

```
docker stop visitors-service
```

## References

* [github.com/jdob/visitors-service](https://github.com/jdob/visitors-service)
* [Managing Multiple Versions of Python on Ubuntu 20.04](https://hackersandslackers.com/multiple-python-versions-ubuntu-20-04/)
* [Install Python, pip, and venv](https://docs.microsoft.com/windows/python/web-frameworks#install-python-pip-and-venv)
