export PATH=$PATH:/usr/local/go/bin:$PWD/bin
go env -w GO111MODULE=on
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
