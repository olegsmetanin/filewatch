FROM golang:1.10.0-alpine3.7 as builder

RUN apk add --no-cache git

RUN go get -u github.com/golang/dep/cmd/dep

ADD . /go/src/olegsmetanin/filewatch

WORKDIR /go/src/olegsmetanin/filewatch

RUN dep init

RUN go build -a -ldflags '-extldflags "-static"' -o filewatch .

FROM alpine:3.7

RUN apk --no-cache add ca-certificates

COPY --from=0 /go/src/olegsmetanin/filewatch/filewatch /usr/local/bin

