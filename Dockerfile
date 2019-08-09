FROM node:12 as node-builder
WORKDIR /tracklog
COPY package.json .babelrc ./
RUN npm install
COPY css css/
COPY js js/
RUN mkdir public && npm run production:build

FROM golang:1.12 as go-builder
WORKDIR /tracklog
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/tracklog-server ./cmd/server
RUN CGO_ENABLED=0 go build -o /bin/tracklog-control ./cmd/control

FROM alpine:3.10
WORKDIR /tracklog
COPY templates templates/
COPY --from=node-builder /tracklog/public public/
COPY --from=go-builder /bin/tracklog-server /bin/tracklog-server
COPY --from=go-builder /bin/tracklog-control /bin/tracklog-control
ENTRYPOINT ["/bin/tracklog-server"]
