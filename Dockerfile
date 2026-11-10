FROM golang:1.27-alpine AS build 
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /auth-service .
FROM alpine:3.20
WORKDIR /app
COPY --from=build /auth-service .
EXPOSE 8005
CMD ["./auth-service"]