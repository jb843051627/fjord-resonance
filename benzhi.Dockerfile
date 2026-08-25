FROM golang:1.22-bookworm

ENV GOPROXY=https://goproxy.cn,direct
ENV GOTOOLCHAIN=local
WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build ./...
RUN go build -o /usr/local/bin/fjord-resonance .

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/fjord-resonance"]
