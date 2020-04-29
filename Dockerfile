# builder 源镜像
FROM    golang:1.12.4-alpine as builder

# 安装git
RUN     apk add --no-cache git gcc musl-dev

#LABEL 更改version后，本地build时LABEL以上的Steps使用Cache
LABEL   maintainer="czhang@pharbers.com" PhAuthServer.version="1.0.27"


# 设置工程配置文件的环境变量 && 开启go-module
ENV     PH_AUTH_HOME $GOPATH/src/github.com/PharbersDeveloper/PhAuthServer/resources
ENV     GOPROXY https://goproxy.io
ENV     GO111MODULE on

RUN     git clone https://github.com/PharbersDeveloper/PhAuthServer $GOPATH/src/github.com/PharbersDeveloper/PhAuthServer

# 设置工作目录
WORKDIR $GOPATH/src/github.com/PharbersDeveloper/PhAuthServer

# 构建可执行文件
RUN     go build
#-a && go install


FROM    alpine:latest

WORKDIR /go/src/github.com/PharbersDeveloper/PhAuthServer/resources/resource
ENV     PH_AUTH_HOME=/go/src/github.com/PharbersDeveloper/PhAuthServer/resources

RUN echo http://mirrors.aliyun.com/alpine/edge/main > /etc/apk/repositories \
    && echo http://mirrors.aliyun.com/alpine/edge/community >> /etc/apk/repositories \
    && apk update \
    && apk add --no-cache tzdata bash \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata \
    && rm -rf /var/cache/apk/*

# 设置工作目录
WORKDIR /go/bin

COPY --from=0 /go/src/github.com/PharbersDeveloper/PhAuthServer/ph_auth .

# 暴露端口
EXPOSE  9096
ENTRYPOINT ["./ph_auth"]
