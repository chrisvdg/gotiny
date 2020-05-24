ARG go_version=latest
ARG alpine_version=latest

FROM golang:$go_version as builder

WORKDIR /go/src/gotiny

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go install -mod vendor


FROM alpine:$alpine_version

# install certificates
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=builder /go/bin/gotiny /bin
RUN mkdir /data

EXPOSE 80 443
ENTRYPOINT [ "gotiny" ]
CMD ["-l", ":80", "-j", "-v", "-f", "/data/backend.js"]