FROM golang:alpine as builder
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
ENV USER=webapp
ENV UID=1001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    webapp
WORKDIR $GOPATH/src/hub4_exporter/
COPY . .
RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0  go build -ldflags="-w -s" -o /go/bin/hub4-exporter
RUN chmod +x /go/bin/hub4-exporter

#############################
FROM scratch
# Import from builder.
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /go/bin/hub4-exporter /hub4-exporter
COPY ./config.yaml /config.yaml
USER webapp:webapp
ENTRYPOINT ["/hub4-exporter"]
