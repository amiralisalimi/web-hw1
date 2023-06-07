FROM golang:latest

ENV PORT=5052

WORKDIR /go/src/auth/
COPY . /go/src/auth/

RUN go build /go/src/auth/server.go

EXPOSE $PORT

ENTRYPOINT [ "/go/src/auth/server" ]
