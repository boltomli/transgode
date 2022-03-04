FROM golang:1.17

RUN sed -i s=http://deb.debian.org=https://mirrors.aliyun.com=g /etc/apt/sources.list
RUN sed -i s=http://security.debian.org=https://mirrors.aliyun.com=g /etc/apt/sources.list
RUN apt-get update && apt-get upgrade -y \
    && apt-get install --no-install-recommends -y libavfilter-dev libavutil-dev libavcodec-dev \
    && rm -rf /var/lib/apt/lists/

WORKDIR /app
COPY main.go .
COPY go.mod .
COPY go.sum .
COPY vendor vendor
RUN go build
ENTRYPOINT [ "./m" ]
