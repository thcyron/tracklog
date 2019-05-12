FROM golang:1.12

COPY . /go/src/github.com/thcyron/tracklog/

ENV GO15VENDOREXPERIMENT 1
ENV CGO_ENABLED 0

WORKDIR /go/src/github.com/thcyron/tracklog/cmd/server
RUN go build

WORKDIR /go/src/github.com/thcyron/tracklog/cmd/control
RUN go build


FROM node:10

COPY . .

RUN npm install
RUN npm run build


FROM scratch

COPY --from=1 public /public
COPY ./templates /templates
COPY --from=0 /go/src/github.com/thcyron/tracklog/cmd/server/server /bin/
COPY --from=0 /go/src/github.com/thcyron/tracklog/cmd/control/control /bin/

CMD ["/bin/server"]
