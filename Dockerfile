# Building - Step 1
FROM golang:1.10 AS build

RUN mkdir -p $GOPATH/src/github.com/anikitenko

WORKDIR $GOPATH/src/github.com/anikitenko

RUN git clone https://github.com/anikitenko/mini-sftp-client.git

RUN go get -u github.com/revel/cmd/revel
RUN go get -u github.com/kardianos/govendor

RUN cd mini-sftp-client && govendor sync

RUN CGO_ENABLED=0 GOOS=linux revel build github.com/anikitenko/mini-sftp-client sftp-client

RUN rm -rf sftp-client/src sftp-client/run.sh sftp-client/run.bat && \
    mv sftp-client/mini-sftp-client sftp-client/mini-sftp-client-linux

RUN CGO_ENABLED=0 GOOS=linux go build -o sftp-client/start-client github.com/anikitenko/mini-sftp-client/run

# Running - Step 2
FROM alpine:3.7

EXPOSE 9000

WORKDIR /app

COPY --from=build /go/src/github.com/anikitenko/sftp-client /app/

ENTRYPOINT ["./start-client", "-p", "9000", "-m", "dev", "--no-ver-check"]