FROM golang:1.22 AS builder

RUN go version
COPY ../ .

RUN cd services/core && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -C services/core -o main main.go

FROM alpine
COPY --from=builder /go/services/core/main .
COPY --from=builder /go/services/core/config.json ./config.json
COPY --from=builder /go/services/core/migrations ./migrations

ENV ADDR=0.0.0.0
ENV LOCAL_CONFIG_FILE=./config.json
EXPOSE 8080 8081
CMD ["./main"]