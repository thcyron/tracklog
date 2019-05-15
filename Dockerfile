FROM golang:1.12

COPY go.* *.go /usr/src/tracklog/
COPY cmd/ /usr/src/tracklog/cmd/
COPY pkg/ /usr/src/tracklog/pkg/

ENV CGO_ENABLED 0

WORKDIR /usr/src/tracklog/cmd/server
RUN go build

WORKDIR /usr/src/tracklog/cmd/control
RUN go build


FROM node:10

COPY package.json .babelrc /usr/src/tracklog/
COPY css/ /usr/src/tracklog/css/
COPY js/ /usr/src/tracklog/js/

WORKDIR /usr/src/tracklog
RUN npm install
RUN npm run build


FROM scratch

COPY --from=1 /usr/src/tracklog/public /public/
COPY ./templates /templates/
COPY --from=0 /usr/src/tracklog/cmd/server/server /bin/
COPY --from=0 /usr/src/tracklog/cmd/control/control /bin/

CMD ["/bin/server"]
