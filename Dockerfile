FROM golang:1.10

EXPOSE 9000

RUN go get -u github.com/revel/cmd/revel
RUN go get -u github.com/kardianos/govendor
RUN mkdir -p $GOPATH/src/github.com/anikitenko
RUN cd $GOPATH/src/github.com/anikitenko && git clone https://github.com/anikitenko/mini-sftp-client.git
RUN cd $GOPATH/src/github.com/anikitenko/mini-sftp-client && govendor sync

CMD ["revel", "run", "github.com/anikitenko/mini-sftp-client"]