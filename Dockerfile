FROM golang:latest

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main

EXPOSE 8080

CMD ["./main"]
 
#BUILD COMMAND : docker build -t app .

#RUN COMMAND : docker run -p 8080:8080 -e PORT=8080 -e MONGO_URI=mongodb://host.docker.internal:27017 -e APP_NAME=sample app
