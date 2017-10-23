FROM golang:1.8.1

ENV TIME_ZONE=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TIME_ZONE /etc/localtime && echo $TIME_ZONE > /etc/timezone

COPY . /go/src/github.com/yiyiyaya/luntan

WORKDIR /go/src/github.com/yiyiyaya/luntan 

RUN go build

CMD ["sh", "-c", "./luntan"]
