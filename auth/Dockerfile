FROM golang:1.16-alpine

ENV GIN_MODE=release
ENV PORT=5052

WORKDIR /go/src/auth

COPY . /go/src/auth/

RUN go get github.com/bsm/ginkgo/v2
RUN go build auth/*.go

EXPOSE $PORT

ENTRYPOINT [ "/go/src/auth/server.go", "-port=$PORT" ]