FROM golang:1.21.3 as builder

WORKDIR /usr/src/app

# COPY go.sum ./
COPY go.mod ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o clean

FROM alpine:latest

COPY --from=builder /usr/src/app/clean /clean

CMD ["/clean"]

# альтернативный вариант

# FROM golang:onbuild AS build
# WORKDIR /build
# COPY . .
# RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./app

# FROM ubuntu:20.04 AS ubuntu
# RUN apt-get update
# RUN DEBIAN_FRONTEND=noninteractive apt-get install -y upx
# RUN apt-get clean && \
#     rm -rf /var/lib/apt/lists/*
# COPY --from=build /build/app /build/app
# RUN upx /build/app

# FROM scratch
# COPY --from=ubuntu /build/app /build/app

# ENTRYPOINT ["/build/app"]