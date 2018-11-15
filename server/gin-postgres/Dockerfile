FROM golang:latest

RUN mkdir -p /app
WORKDIR /app

RUN go get -u github.com/gin-gonic/gin \
    && go get -u github.com/jinzhu/gorm \
    && go get -u github.com/jinzhu/gorm/dialects/postgres
