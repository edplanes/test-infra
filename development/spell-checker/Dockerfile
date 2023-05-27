FROM golang:1.20-alpine3.18 AS builder

WORKDIR $GOPATH/src/github.com/edplanes/test-infra
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./development/spell-checker -a -installsuffix cgo -ldflags="-w -s" -o /spell-checker

FROM debian:10.13-slim

USER spellchecker

COPY --from=builder /spell-checker /spell-checker
RUN apt-get update -y && apt-get install -y aspell aspell-en aspell-pl

HEALTHCHECK NONE

ENTRYPOINT ["/spell-checker"]
