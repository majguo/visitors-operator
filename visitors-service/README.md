# Visitors Serivce

The code base of visitors service in this repo references https://github.com/jdob/visitors-service.

It's dependent by the visitors web ui which is located at [visitors-webui](../visitors-webui/).

## Run in local

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

   Quit the server with CONTROL-C.

1. Quit the virtual environment:

   ```
   deactivate
   ```
