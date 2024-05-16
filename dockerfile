FROM golang:1.22 as builder
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE="on"
WORKDIR /
COPY . .
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o pay

FROM alpine
WORKDIR /root/
COPY --from=builder /pay .
COPY --from=builder /config/config.yml ./config/
COPY --from=builder /config/com.pem ./config/
# 若数据库为公网则无需暴漏3306
EXPOSE 3636 3306
CMD ["ls"]
ENTRYPOINT [ "/root/pay" ]

# 程序打包 docker build -t pay .
# mysql 配置 docker run --name mysql -e MYSQL_ROOT_PASSWORD=123456 -d mysql:latest
# 网桥 docker network create bridge_1
# mysql上网桥 docker network connect bridge_1 mysql
# pay上网桥 docker network connect bridge_1 pay
# 查看ip地址 docker network inspect bridge_1
# 修改yml的ip  host: xxx.xx.xx.x


