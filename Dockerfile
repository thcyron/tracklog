FROM scratch

COPY ["public", "/public"]
COPY ["templates", "/templates"]
COPY ["cmd/server/server", "cmd/control/control", "/bin/"]

CMD ["/bin/server"]
