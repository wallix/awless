FROM alpine:latest

RUN apk --update add ca-certificates

COPY awless-scheduler /usr/bin/awless-scheduler

EXPOSE 8082
EXPOSE 8083

VOLUME ["/root/.awless-scheduler"]

ENTRYPOINT ["/usr/bin/awless-scheduler", "-discovery-hostport", ":8082", "scheduler-hostport", ":8083", "-http-mode"]