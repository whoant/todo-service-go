FROM todo-service-cached as builder

WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
ENTRYPOINT ["/app/main"]
