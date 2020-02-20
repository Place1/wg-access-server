FROM node:10 as website
WORKDIR /code
COPY ./website/package.json ./
COPY ./website/package-lock.json ./
RUN npm install
COPY ./website/ ./
RUN npm run build

FROM golang:1.13.8 as server
WORKDIR /code
ENV GOOS=linux
ENV GARCH=amd64
ENV CGO_ENABLED=0
ENV GO111MODULE=on
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY ./proto/ ./proto/
COPY ./main.go ./main.go
COPY ./internal/ ./internal
RUN go build -o server

FROM alpine:3.10
RUN apk add iptables
RUN apk add wireguard-tools
ENV CONFIG="/config.yaml"
ENV STORAGE_DIRECTORY="/data"
COPY --from=server /code/server /server
COPY --from=website /code/build /website/build
CMD /server
