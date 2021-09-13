#FROM docker.u51-inc.com/library/centos69:51nb20180620
#
#ADD jmx-puller /usr/local/bin
#
#ENTRYPOINT jmx-puller

# FROM docker.u51-inc.com/library/golang-builder:0.0.1 as builder
FROM golang:1.16 as builder

ARG GO_PROXY=https://goproxy.cn

WORKDIR /usr/src/app

ENV GOPROXY=$GO_PROXY

COPY . .

RUN apt-get update && apt-get install unzip

RUN cd ../ && unzip ./app/soar.zip && cd ./app && go mod download

RUN go build -o zst_soar ./main.go

FROM ubuntu:20.04 as runner
COPY --from=builder /usr/src/app/zst_soar /opt/app/
CMD ["/opt/app/zst_soar"]

