FROM golang:1.8 AS go
ENV GOBIN /go/bin
RUN mkdir -p /go/src/github.com/oceanoverflow/sidecar
COPY . /go/src/github.com/oceanoverflow/sidecar
WORKDIR /go/src/github.com/oceanoverflow/sidecar
RUN go get -u github.com/golang/dep/...
RUN dep ensure
RUN go build -o sidecar .

FROM registry.cn-hangzhou.aliyuncs.com/aliware2018/services AS builder

FROM registry.cn-hangzhou.aliyuncs.com/aliware2018/debian-jdk8
COPY --from=builder /root/workspace/services/mesh-provider/target/mesh-provider-1.0-SNAPSHOT.jar /root/dists/mesh-provider.jar
COPY --from=builder /root/workspace/services/mesh-consumer/target/mesh-consumer-1.0-SNAPSHOT.jar /root/dists/mesh-consumer.jar
COPY --from=builder /usr/local/bin/docker-entrypoint.sh /usr/local/bin
COPY --from=go /go/src/github.com/oceanoverflow/sidecar/start-agent.sh /usr/local/bin
COPY --from=go /go/src/github.com/oceanoverflow/sidecar/sidecar /usr/local/bin
RUN set -ex && chmod a+x /usr/local/bin/start-agent.sh && mkdir -p /root/logs
ENTRYPOINT ["docker-entrypoint.sh"]