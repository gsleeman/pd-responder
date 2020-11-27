FROM golang:1.15

RUN mkdir -p /go/src/github.com/gsleeman/pd-responder
WORKDIR /go/src/github.com/gsleeman/pd-responder
EXPOSE 8888
COPY responder.go /go/src/github.com/gsleeman/pd-responder
RUN go build

CMD ["./pd-responder"]
