FROM alpine:3.3

RUN mkdir -p /usr/local/share/tracklog
COPY ["dist/public", "/usr/local/share/tracklog/public"]
COPY ["dist/templates", "/usr/local/share/tracklog/templates"]
COPY ["dist/tracklog-server", "dist/tracklog-import", "/usr/local/bin/"]

WORKDIR /usr/local/share/tracklog
CMD ["/usr/local/bin/tracklog-server", "-config", "/etc/tracklog.json"]
