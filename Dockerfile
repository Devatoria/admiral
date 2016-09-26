FROM golang:1.7.1

ADD . /go/src/github.com/Devatoria/admiral
WORKDIR /go/src/github.com/Devatoria/admiral

RUN go get -u github.com/kardianos/govendor \
    && govendor sync \
    && go build

ENTRYPOINT ./admiral
