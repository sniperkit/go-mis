FROM golang:1.9-stretch

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
ADD . /app

RUN chmod +x go-mis

# Make port 80 available to the world outside this container
EXPOSE 80

# Define environment variable
ENV GOOS linux
ENV GOARCH amd64

# Run app when the container launches
CMD ./go-mis