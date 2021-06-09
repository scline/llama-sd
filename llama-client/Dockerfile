FROM python:3.9

WORKDIR /app

COPY requirements.txt .

# install dependencies
RUN pip install -r requirements.txt

# copy the content of the local src directory to the working directory
COPY $PDW/src/ .

# environment variables
ENV \
    APP_CONFIG="config.yml" \
    PYTHONUNBUFFERED=0

# application spisific environment variables (overwrites config file if used)
#ENV \
#    LLAMA_SERVER="http://127.0.0.1" \
#    LLAMA_KEEPALIVE=300 \
#    LLAMA_GROUP="default" \
#    LLAMA_SOURCE_IP="192.168.1.2"
#    LLAMA_TAGS="[ tag1: value1, tag2: value2 ]" \
#    APP_VERBOSE=True

# command to run on container start
CMD [ "python", "./app.py" ]
