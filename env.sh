export PATH=$PATH:/usr/local/go/bin # add Go path
export PATH=$PATH:/home/fabric/phd-2/cybernet/creditdb/bin # add peer path
go env -w GO111MODULE=on
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
