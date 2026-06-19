FROM golang:1.24 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ติดตั้ง CA cert (สำคัญถ้าเรียก HTTPS / DB / API)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/api

FROM alpine:latest

WORKDIR /app
RUN apk --no-cache add ca-certificates

# copy binary จาก build stage
COPY --from=builder /app/main .

# expose port (ปรับตาม app คุณ)
EXPOSE 8080

# run app
CMD ["./main"]
