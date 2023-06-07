FROM golang:latest

ENV PORT=5062

WORKDIR /go/src/biz/
COPY . /go/src/biz/

RUN go build /go/src/biz/server.go

EXPOSE $PORT

ENTRYPOINT [ "/go/src/biz/server" ]
