FROM python:3.9

WORKDIR /app

COPY requirements.txt .

# install dependencies
RUN pip install -r requirements.txt

# copy the content of the local src directory to the working directory
COPY src/ .

# environment variables
ENV \
    APP_PORT=80 \
    APP_HOST=0.0.0.0 \
    APP_VERBOSE=True \
    APP_KEEPALIVE=86400 \
    PYTHONUNBUFFERED=0

# Expose webport
EXPOSE 80

# command to run on container start
CMD [ "python", "./app.py" ]
