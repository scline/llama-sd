FROM python:3.9

WORKDIR /app

COPY requirements.txt .

# install dependencies
RUN pip install -r requirements.txt

# copy the content of the local src directory to the working directory
COPY src/ .

# environment variables
ENV \
    APP_CONFIG="config.yml" \
    PYTHONUNBUFFERED=0

# application spisific environment variables (overwrites config file if used)
#ENV \
#    APP_PORT=80 \
#    APP_HOST=0.0.0.0 \
#    APP_VERBOSE=False \
#    APP_GROUP="none" \
#    APP_KEEPALIVE=86400 \
#    LLAMA_INTERVAL=10 \
#    APP_VERBOSE=True

# Expose webport
EXPOSE 80

# command to run on container start
CMD [ "python", "./main.py" ]
