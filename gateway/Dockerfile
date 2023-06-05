FROM golang:1.16-alpine

ENV GIN_MODE=release
ENV PORT=6433

WORKDIR /go/src/gateway

COPY . /go/src/gateway/

RUN go mod download github.com/bsm/ginkgo/v2
RUN go build gateway/*.go

EXPOSE $PORT

ENTRYPOINT [ "/go/src/gateway/server", "-port=$PORT" ]
