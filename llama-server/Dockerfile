FROM python:3.11

WORKDIR /app

COPY requirements.txt .

# install dependencies
RUN pip install -r requirements.txt

# copy the content of the local src directory to the working directory
COPY src/ .

# environment variables
ENV PYTHONUNBUFFERED=0

# Expose webport
EXPOSE 80

# command to run on container start
CMD [ "python", "./main.py" ]

# for debugging container
#ENTRYPOINT ["tail", "-f", "/dev/null"]