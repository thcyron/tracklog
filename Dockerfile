FROM golang:1.12

COPY . /usr/src/tracklog

ENV CGO_ENABLED 0

WORKDIR /usr/src/tracklog/cmd/server
RUN go build

WORKDIR /usr/src/tracklog/cmd/control
RUN go build


FROM node:10

COPY . .

RUN npm install
RUN npm run build


FROM scratch

COPY --from=1 public /public
COPY ./templates /templates
COPY --from=0 /usr/src/tracklog/cmd/server/server /bin/
COPY --from=0 /usr/src/tracklog/cmd/control/control /bin/

CMD ["/bin/server"]
