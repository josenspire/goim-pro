FROM golang:alpine3.11

MAINTAINER JAMES YANG "josenspire@gmail.com"

WORKDIR /$GOPATH/src/
COPY . .

RUN go build -mod=vendor -v cmd/main.go

EXPOSE 9090

ENTRYPOINT ["./main"]

# docker build -t josenspire/goim-pro .
# docker run --env APP_ENV=PROD -p 9090:9090 --name goim-pro 0620
