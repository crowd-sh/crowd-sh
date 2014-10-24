DOCKER_NAME=workmachine/mart

all: build start

build: *
    sudo docker build -t $(DOCKER_NAME) .

start:
    # Stop the existing instance
    sudo docker ps | grep $(DOCKER_NAME) | awk '{print $1}' | xargs sudo docker kill
    # Start the new instance
    sudo docker run -d -p 127.0.0.1:3000:3000 $(DOCKER_NAME)
