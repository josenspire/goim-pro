FROM golang:alpine3.11

MAINTAINER JAMES YANG "josenspire@gmail.com"

WORKDIR /$GOPATH/src/
COPY . .

RUN chmod 777 ./wait-for-it.sh

RUN apk update && apk add bash

RUN go build -mod=vendor -v cmd/main.go

EXPOSE 9090

CMD ["./main"]

# docker build -t josenspire/goim-pro .
# docker run -d -i -t --env APP_ENV=PROD -p 9090:9090 --name goim-pro 0620
# docker run -d -i -t --env APP_ENV=PROD -p 9090:9090 --name goim-pro --network app-net f061
