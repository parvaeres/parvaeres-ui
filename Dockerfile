FROM golang:1.15

WORKDIR /go/src/github.com/parvaeres/go8s

COPY . .

RUN go get github.com/revel/revel
RUN go get github.com/revel/cmd/revel

ENTRYPOINT revel run /go/src/github.com/parvaeres/go8s dev 9000

EXPOSE 9000
