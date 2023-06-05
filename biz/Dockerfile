FROM golang:1.16-alpine

ENV GIN_MODE=release
ENV PORT=5062

WORKDIR /go/src/biz

COPY . /go/src/biz/

RUN go mod download github.com/bsm/ginkgo/v2
RUN go build biz/*.go

EXPOSE $PORT

ENTRYPOINT [ "/go/src/biz/server", "-port=$PORT" ]
