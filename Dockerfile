FROM golang:1.11 
RUN mkdir /app 
COPY . /app/ 
WORKDIR /app
RUN go get ./... 
RUN go build -o main .