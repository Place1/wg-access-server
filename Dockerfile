### Build stage for the website frontend
FROM node:10 as website
RUN apt-get update
RUN apt-get install -y protobuf-compiler
WORKDIR /code
COPY ./website/package.json ./
COPY ./website/package-lock.json ./
RUN npm ci --no-audit --prefer-offline
COPY ./proto/ ./proto/
COPY ./website/ ./
RUN npm run build

### Build stage for the website backend server
FROM golang:1.13.8 as server
RUN apt-get update
RUN apt-get install -y protobuf-compiler
WORKDIR /code
ENV GOOS=linux
ENV GARCH=amd64
ENV CGO_ENABLED=0
ENV GO111MODULE=on
RUN go get github.com/golang/protobuf/protoc-gen-go@v1.3.5
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY ./proto/ ./proto/
COPY ./codegen.sh ./
RUN ./codegen.sh
COPY ./main.go ./main.go
COPY ./internal/ ./internal/
COPY ./pkg/ ./pkg/
RUN go build -o server

### Server
FROM alpine:3.10
RUN apk add iptables
RUN apk add wireguard-tools
RUN apk add curl
ENV CONFIG="/config.yaml"
ENV STORAGE="file:///data"
COPY --from=server /code/server /server
COPY --from=website /code/build /website/build
CMD /server
