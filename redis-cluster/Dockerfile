FROM alpine:redis-cluster
MAINTAINER rain
WORKDIR /MockRedis
ENV REDIS_HOME=/MockRedis
ADD redis-cluster /go/bin/redis-cluster
ENTRYPOINT ["/go/bin/redis-cluster"]
EXPOSE 6381-6386