# ベースとなるDockerイメージ指定
FROM golang:1.16.6-alpine3.14 as build

RUN apk add --update --no-cache ca-certificates git make

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src/api

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN make build

FROM alpine

COPY --from=build /go/src/api/api /api
ENTRYPOINT ["./api"]

EXPOSE 5000