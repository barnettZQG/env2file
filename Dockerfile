FROM golang:1.11 as builder
ADD . /go/src/env2file
WORKDIR /go/src/env2file
RUN go build -ldflags " -w" -o build/env2file

FROM alpine

COPY --from=builder /go/src/env2file/build/env2file /usr/local/bin/env2file
