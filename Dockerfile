#源镜像
FROM    golang:1.12.4-alpine

# 安装git
RUN     apk add --no-cache git gcc musl-dev

#LABEL 更改version后，本地build时LABEL以上的Steps使用Cache
LABEL   maintainer="czhang@pharbers.com" PhAuthServer.version="1.0.12"

# 设置工程配置文件的环境变量 && 开启go-module
ENV     PH_AUTH_HOME $GOPATH/src/github.com/PharbersDeveloper/PhAuthServer/resources
ENV     GO111MODULE on

# 下载依赖
RUN     git clone https://github.com/PharbersDeveloper/PhAuthServer $GOPATH/src/github.com/PharbersDeveloper/PhAuthServer && \
        ln -sf $GOPATH/src/github.com/PharbersDeveloper/PhAuthServer/static  $GOPATH/bin/static

# 设置工作目录
WORKDIR $GOPATH/src/github.com/PharbersDeveloper/PhAuthServer

# 构建可执行文件
RUN     go build -a && go install

# 暴露端口
EXPOSE  9096

# 设置工作目录
WORKDIR $GOPATH/bin

ENTRYPOINT ["ph_auth"]
