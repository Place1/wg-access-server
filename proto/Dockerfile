FROM node:17-bullseye as proto-js
RUN apt-get update && apt-get install -y protobuf-compiler libprotobuf-dev
WORKDIR /code
COPY ./website/package.json ./
COPY ./website/package-lock.json ./
RUN npm ci --no-audit --prefer-offline
# Mount ./proto to /proto
# Mount ./website/src/sdk to /code/src/sdk
CMD ["npm", "run", "codegen"]

FROM golang:1.18-alpine as proto-go
RUN apk add --no-cache protobuf protobuf-dev
WORKDIR /code
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
# Mount ./proto to /code/proto
COPY ./codegen.sh ./
CMD ["./codegen.sh"]
