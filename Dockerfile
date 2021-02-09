FROM golang:latest

RUN sed -i 's/deb.debian.org/mirrors.163.com/g' /etc/apt/sources.list && \
    sed -i 's/security.debian.org/mirrors.163.com/g' /etc/apt/sources.list && \
    apt-get update

RUN apt-get install -y wkhtmltopdf && \
    mkdir -p /usr/share/fonts/chinese/TrueType

COPY pkg/wkhtmltopdf/fonts/* /usr/share/fonts/chinese/TrueType/

RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,direct

WORKDIR /go/build
COPY . .
RUN go mod download && \
    go build -o save-page ./cmd/save-page.go && \
    mkdir -p /app && \
    cp /go/build/save-page /go/build/config.json /app/
WORKDIR /app
ENV PATH /app:$PATH

COPY docker-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["docker-entrypoint.sh"]
