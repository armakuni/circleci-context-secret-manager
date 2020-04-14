FROM circleci/golang:1.13

WORKDIR /go

COPY ccsm /go/bin/.

RUN apt-get update && apt-get install -y git-crypt

CMD ["./bin/ccsm"]
