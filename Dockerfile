FROM golang:1.13.8-alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git tzdata

ENV USER=appuser
ENV UID=10001 
# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"
WORKDIR $GOPATH/src/github.com/delphis-inc/delphisbe/
COPY . .
RUN go mod download
RUN go mod verify
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/delphis_server
COPY ./config /go/bin

FROM alpine

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

RUN mkdir -p /go/bin

COPY --from=builder /go/bin/delphis_server /go/bin/delphis_server
COPY config /var/delphis/config

USER appuser:appuser

EXPOSE 8080

ENTRYPOINT ["/go/bin/delphis_server"]