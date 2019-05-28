FROM golang:alpine as builder 

RUN apk --no-cache add ca-certificates

COPY src /go/src/lock-manager
WORKDIR /go/src/lock-manager

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/lockmgr main.go


FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/lockmgr /lockmgr

ENTRYPOINT [ "/lockmgr" ]
