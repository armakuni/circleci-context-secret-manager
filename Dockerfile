FROM circleci/golang:1.13

WORKDIR /go

COPY ccsm /go/bin/.

CMD ["./ccsm"]
