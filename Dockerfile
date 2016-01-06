FROM alpine:3.3

RUN apk add --update go && rm -rf /var/cache/apk/*

ENV TRACKLOG_ROOT=/go/src/github.com/thcyron/tracklog
COPY . $TRACKLOG_ROOT
WORKDIR $TRACKLOG_ROOT

ENV GOPATH=/go GO15VENDOREXPERIMENT=1
RUN (cd cmd/server && go build -o /usr/local/bin/tracklog-server)
RUN (cd cmd/import && go build -o /usr/local/bin/tracklog-import)

RUN mkdir -p /usr/local/share/tracklog
RUN mv public/ templates/ /usr/local/share/tracklog
RUN rm -rf /go
RUN apk del go

WORKDIR /usr/local/share/tracklog
ENTRYPOINT ["/usr/local/bin/tracklog-server"]
CMD ["-config", "/etc/tracklog.json"]
