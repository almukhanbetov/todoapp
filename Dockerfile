# Build stage
FROM golang:1.24-alpine AS build

WORKDIR /app
COPY . .

RUN apk add --no-cache git
RUN go mod init todoapp && go mod tidy
RUN go build -o todoapp ./main.go

# Runtime stage
FROM alpine:3.19
WORKDIR /app

COPY --from=build /app/todoapp /app/

EXPOSE 8081
CMD ["./todoapp"]
