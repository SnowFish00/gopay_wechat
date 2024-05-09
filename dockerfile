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
EXPOSE 3636
CMD ["ls"]
ENTRYPOINT [ "./pay" ]