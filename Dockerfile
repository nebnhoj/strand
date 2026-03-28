# STEP 1 build executable binary
FROM golang:alpine AS builder

RUN adduser -D -g '' appuser
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/app .

RUN chmod 0555 /go/bin/app

# STEP 2 build a small image
FROM scratch

WORKDIR /app

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/app ./app_binary
COPY --from=builder /app/docs ./docs

USER appuser

EXPOSE 3000

CMD ["./app_binary"]
