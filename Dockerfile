### Build stage for the website frontend 
FROM node:10 as website

WORKDIR /code

COPY ./website/package.json ./
COPY ./website/package-lock.json ./

# install dependency
RUN npm install

COPY ./website/ ./

RUN npm run build

### Build stage for the website backend server
FROM golang:1.13.8 as server

WORKDIR /code

# environment variable
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

### Server 
FROM alpine:3.10

# environment variable
ENV CONFIG="/config.yaml"
ENV STORAGE_DIRECTORY="/data"

RUN apk add iptables
RUN apk add wireguard-tools

# Copy the final build for the frontend and backend
COPY --from=server /code/server /server
COPY --from=website /code/build /website/build
CMD /server
