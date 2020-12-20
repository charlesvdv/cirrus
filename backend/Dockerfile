FROM golang:alpine AS builder

ENV GO111MODULE=on

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o cirrus ./cmd

FROM alpine
COPY --from=builder /build/cirrus /cirrus

EXPOSE 80
ENTRYPOINT ["/cirrus", "-p", "80"]