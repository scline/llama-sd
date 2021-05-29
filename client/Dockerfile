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
    APP_GROUP="none" \
    APP_KEEPALIVE=86400 \
    APP_TAGS="hostname: las1.something.com, src_short: las1, city: Las Vegas, Nevada" \
    PYTHONUNBUFFERED=0

# Expose webport
EXPOSE 80

# command to run on container start
CMD [ "python", "./app.py" ]
