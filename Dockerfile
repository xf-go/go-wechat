FROM golang:1.15.3-alpine3.12

WORKDIR /data
COPY . .
ENV GOPROXY=https://goproxy.io
RUN go build -o go-wechat

# RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone

EXPOSE 8121
CMD [ "./go-wechat" ]