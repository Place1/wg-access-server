FROM rust:1.38.0-alpine as boringtun
RUN apk add musl-dev
RUN cargo install \
  --root / \
  --bin boringtun \
  --git https://github.com/cloudflare/boringtun.git \
  --rev 0b980a2f5a5f8622bf0e3a024ace63ad01b5d0f6

FROM node:10 as website
WORKDIR /code
COPY ./website/package.json ./
COPY ./website/package-lock.json ./
RUN npm install
COPY ./website/ ./
RUN npm run build

FROM golang:1.13 as server
WORKDIR /code
ENV GOOS=linux
ENV GARCH=amd64
ENV CGO_ENABLED=0
ENV GO111MODULE=on
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY ./main.go ./main.go
COPY ./internal/ ./internal
RUN go build -o server

FROM alpine:3.10
RUN apk add iptables
RUN apk add wireguard-tools
ENV WIREGUARD_USERSPACE_IMPLEMENTATION=boringtun
ENV STORAGE_DIRECTORY="/data"
COPY --from=boringtun /bin/boringtun /usr/local/bin/boringtun
COPY --from=server /code/server ./server
COPY --from=website /code/build ./website/build
CMD ./server
