FROM golang:1.23-alpine

ARG USER_ID
ARG GROUP_ID

ENV GOFLAGS="-buildvcs=false"

WORKDIR /app

RUN apk add --no-cache git \
    && go install github.com/air-verse/air@v1.61.5 \
    && go install github.com/go-delve/delve/cmd/dlv@v1.24.0 \
    && go install github.com/pressly/goose/v3/cmd/goose@v3.24.0 \
    && go install golang.org/x/tools/cmd/godoc@latest \
    && go install go.uber.org/mock/mockgen@latest \
    && go install github.com/swaggo/swag/cmd/swag@latest \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]
