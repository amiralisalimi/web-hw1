FROM golang:latest

ENV PORT=6433

WORKDIR /go/src/gateway/
COPY . /go/src/gateway/

RUN go build /go/src/gateway/server.go

EXPOSE $PORT

ENTRYPOINT [ "/go/src/gateway/server" ]
