FROM golang:1.21.3-bullseye AS build_base
RUN apt-get update && apt-get install -y git gcc ca-certificates libc6-dev && rm -rf /var/lib/apt/lists/*
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -ldflags "-w -s" -trimpath -o speedtest .

FROM debian:bullseye-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=build_base /build/speedtest ./
COPY settings.toml ./

USER nobody
EXPOSE 8989

CMD ["./speedtest"]
