FROM circleci/golang:1.13

WORKDIR /root/

COPY ccsm .

CMD ["./ccsm"]
