FROM golang:1.20-alpine3.18 AS builder

WORKDIR $GOPATH/src/github.com/edplanes/test-infra
COPY . .

RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

ENV USER=appuser
ENV UID=10001

RUN adduser \    
  --disabled-password \    
  --gecos "" \    
  --home "/nonexistent" \    
  --shell "/sbin/nologin" \    
  --no-create-home \    
  --uid "${UID}" \    
  "${USER}"

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./development/image-cleaner -ldflags="-w -s" -o /image-cleaner

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /image-cleaner /image-cleaner

USER appuser:appuser

HEALTHCHECK NONE

ENTRYPOINT [ "/image-cleaner" ]