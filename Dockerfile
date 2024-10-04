ARG GO_VERSION=1.22

FROM golang:${GO_VERSION}-alpine AS builder

RUN go env -w GOPROXY=direct
RUN apk add --no-cache git
RUN apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /src 

COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./

RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /chi-api cmd/main.go

FROM scratch AS runner

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY .env ./

COPY --from=builder /chi-api /chi-api

EXPOSE 9000

ENTRYPOINT ["/chi-api"]
