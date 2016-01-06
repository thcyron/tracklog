FROM alpine:3.3

ENV GOPATH=/go GO15VENDOREXPERIMENT=1
ENV TRACKLOG_ROOT=/go/src/github.com/thcyron/tracklog

RUN apk add --update go && rm -rf /var/cache/apk/*
RUN mkdir -p $TRACKLOG_ROOT
COPY . $TRACKLOG_ROOT

WORKDIR $TRACKLOG_ROOT

RUN (cd cmd/server && go build)
RUN (cd cmd/import && go build)

ENTRYPOINT ["cmd/server/server"]
CMD ["-config", "config.json"]
