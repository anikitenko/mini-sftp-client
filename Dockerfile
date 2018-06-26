# Building - Step 1
FROM golang:1.9 AS build

RUN mkdir -p $GOPATH/src/github.com/anikitenko

WORKDIR $GOPATH/src/github.com/anikitenko

RUN git clone https://github.com/anikitenko/mini-sftp-client.git

RUN go get -u github.com/revel/cmd/revel
RUN go get -u github.com/kardianos/govendor
RUN go get -v github.com/swaggo/swag/cmd/swag

RUN cd mini-sftp-client && govendor sync && swag init -g api_v1.go -d app/controllers

RUN CGO_ENABLED=0 GOOS=linux revel build github.com/anikitenko/mini-sftp-client sftp-client

RUN rm -f sftp-client/run.sh sftp-client/run.bat && \
    mv sftp-client/mini-sftp-client sftp-client/mini-sftp-client-linux

RUN find sftp-client/src/github.com/anikitenko/mini-sftp-client \
 -maxdepth 1 ! -path sftp-client/src/github.com/anikitenko/mini-sftp-client \
 -not -name app -not -name conf -not -name public -exec rm -rf {} +

RUN cd mini-sftp-client/run && govendor sync && CGO_ENABLED=0 GOOS=linux go build -o ../../sftp-client/start-client

# Running - Step 2
FROM alpine:3.7

EXPOSE 9000

WORKDIR /app

COPY --from=build /go/src/github.com/anikitenko/sftp-client /app/

CMD ["./start-client", "-p", "9000", "-m", "dev", "--no-ver-check"]
