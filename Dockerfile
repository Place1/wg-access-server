### Build stage for the website frontend
FROM node:17-bullseye as website
RUN apt-get update
RUN apt-get install -y protobuf-compiler libprotobuf-dev
WORKDIR /code
COPY ./website/package.json ./
COPY ./website/package-lock.json ./
RUN npm ci --no-audit --prefer-offline
COPY ./proto/ ../proto/
COPY ./website/ ./
RUN npm run codegen
RUN npm run build

### Build stage for the website backend server
FROM golang:1.17.5-alpine as server
RUN apk add --no-cache gcc musl-dev protobuf protobuf-dev
WORKDIR /code
ENV CGO_ENABLED=1
ENV GO111MODULE=on
RUN go install github.com/golang/protobuf/protoc-gen-go@v1.5.2
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
RUN go mod verify
COPY ./proto/ ./proto/
COPY ./codegen.sh ./
RUN ./codegen.sh
COPY ./main.go ./main.go
COPY ./cmd/ ./cmd/
COPY ./pkg/ ./pkg/
COPY ./internal/ ./internal/
RUN go build -o wg-access-server

### Server
FROM alpine:3.15.0
RUN apk add --no-cache iptables ip6tables wireguard-tools curl
ENV WG_CONFIG="/config.yaml"
ENV WG_STORAGE="sqlite3:///data/db.sqlite3"
COPY --from=server /code/wg-access-server /usr/local/bin/wg-access-server
COPY --from=website /code/build /website/build
CMD ["wg-access-server", "serve"]
