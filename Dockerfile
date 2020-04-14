FROM circleci/golang:1.13

WORKDIR /go

COPY ccsm /go/bin/.

USER root
RUN apt-get update && apt-get install -y git-crypt

USER circleci
CMD ["./bin/ccsm"]
